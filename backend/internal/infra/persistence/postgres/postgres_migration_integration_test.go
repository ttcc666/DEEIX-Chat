package postgres

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/DEEIX-AI/DEEIX-Chat/backend/internal/infra/config"
	persistenceSchema "github.com/DEEIX-AI/DEEIX-Chat/backend/internal/infra/persistence/schema"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestNewMigratesAllRegisteredModelsInIsolatedPostgresSchema(t *testing.T) {
	dsn := strings.TrimSpace(os.Getenv("DEEIX_TEST_DATABASE_DSN"))
	if dsn == "" {
		t.Skip("set DEEIX_TEST_DATABASE_DSN to run isolated PostgreSQL migration catalog test")
	}

	adminDB, err := gorm.Open(gormpostgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open PostgreSQL admin connection: %v", err)
	}
	adminSQLDB, err := adminDB.DB()
	if err != nil {
		t.Fatalf("resolve PostgreSQL admin connection: %v", err)
	}

	schemaName := fmt.Sprintf("deeix_migration_catalog_%d", time.Now().UnixNano())
	if err = adminDB.Exec("CREATE SCHEMA " + schemaName).Error; err != nil {
		_ = adminSQLDB.Close()
		t.Fatalf("create isolated schema: %v", err)
	}
	t.Cleanup(func() {
		_ = adminDB.Exec("DROP SCHEMA IF EXISTS " + schemaName + " CASCADE").Error
		_ = adminSQLDB.Close()
	})

	isolatedURL, err := url.Parse(dsn)
	if err != nil {
		t.Fatalf("parse PostgreSQL DSN: %v", err)
	}
	query := isolatedURL.Query()
	query.Set("options", "-c search_path="+schemaName+",public")
	isolatedURL.RawQuery = query.Encode()

	db, err := New(config.Config{PostgresDSN: isolatedURL.String()})
	if err != nil {
		t.Fatalf("run full PostgreSQL migration: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("resolve migrated PostgreSQL connection: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })

	models := persistenceSchema.Models()
	if len(models) != 65 {
		t.Fatalf("registered model count = %d, want 65", len(models))
	}

	var columnCount, indexCount, constraintCount int64
	if err = db.Raw("SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = current_schema()").Scan(&columnCount).Error; err != nil {
		t.Fatalf("count migrated columns: %v", err)
	}
	if err = db.Raw("SELECT COUNT(*) FROM pg_indexes WHERE schemaname = current_schema()").Scan(&indexCount).Error; err != nil {
		t.Fatalf("count migrated indexes: %v", err)
	}
	if err = db.Raw("SELECT COUNT(*) FROM pg_constraint c JOIN pg_namespace n ON n.oid = c.connamespace WHERE n.nspname = current_schema()").Scan(&constraintCount).Error; err != nil {
		t.Fatalf("count migrated constraints: %v", err)
	}

	for _, model := range models {
		if !db.Migrator().HasTable(model) {
			t.Errorf("registered model table is missing: %T", model)
			continue
		}
		columns, columnErr := db.Migrator().ColumnTypes(model)
		if columnErr != nil {
			t.Errorf("inspect columns for %T: %v", model, columnErr)
			continue
		}
		if len(columns) == 0 {
			t.Errorf("registered model has no columns: %T", model)
		}
	}

	if columnCount == 0 || indexCount == 0 {
		t.Fatalf("migration catalog is unexpectedly empty: columns=%d indexes=%d", columnCount, indexCount)
	}
	t.Logf(
		"POSTGRES_MIGRATION_CATALOG schema=%s tables=%d columns=%d indexes=%d constraints=%d",
		schemaName,
		len(models),
		columnCount,
		indexCount,
		constraintCount,
	)
	if err = writePostgresCatalogSnapshot(db, schemaName); err != nil {
		t.Fatalf("export PostgreSQL migration catalog: %v", err)
	}
}

type postgresCatalogSnapshot struct {
	Schema      string                      `json:"schema"`
	Tables      []string                    `json:"tables"`
	Columns     []postgresCatalogColumn     `json:"columns"`
	Indexes     []postgresCatalogIndex      `json:"indexes"`
	Constraints []postgresCatalogConstraint `json:"constraints"`
}

type postgresCatalogColumn struct {
	TableName       string `json:"table_name" gorm:"column:table_name"`
	ColumnName      string `json:"column_name" gorm:"column:column_name"`
	DataType        string `json:"data_type" gorm:"column:data_type"`
	UDTName         string `json:"udt_name" gorm:"column:udt_name"`
	IsNullable      string `json:"is_nullable" gorm:"column:is_nullable"`
	ColumnDefault   string `json:"column_default" gorm:"column:column_default"`
	OrdinalPosition int    `json:"ordinal_position" gorm:"column:ordinal_position"`
}

type postgresCatalogIndex struct {
	TableName string `json:"table_name" gorm:"column:table_name"`
	IndexName string `json:"index_name" gorm:"column:index_name"`
	IndexDef  string `json:"index_def" gorm:"column:index_def"`
}

type postgresCatalogConstraint struct {
	TableName      string `json:"table_name" gorm:"column:table_name"`
	ConstraintName string `json:"constraint_name" gorm:"column:constraint_name"`
	ConstraintType string `json:"constraint_type" gorm:"column:constraint_type"`
	Definition     string `json:"definition" gorm:"column:definition"`
}

func writePostgresCatalogSnapshot(db *gorm.DB, schemaName string) error {
	outputDir := strings.TrimSpace(os.Getenv("DEEIX_TEST_CATALOG_DIR"))
	if outputDir == "" {
		return nil
	}
	snapshot := postgresCatalogSnapshot{Schema: schemaName}
	var tableRows []struct {
		TableName string `gorm:"column:table_name"`
	}
	if err := db.Raw(`SELECT table_name FROM information_schema.tables WHERE table_schema = current_schema() AND table_type = 'BASE TABLE' ORDER BY table_name`).Scan(&tableRows).Error; err != nil {
		return fmt.Errorf("query tables: %w", err)
	}
	for _, row := range tableRows {
		snapshot.Tables = append(snapshot.Tables, row.TableName)
	}
	if err := db.Raw(`SELECT table_name, column_name, data_type, udt_name, is_nullable, COALESCE(column_default, '') AS column_default, ordinal_position FROM information_schema.columns WHERE table_schema = current_schema() ORDER BY table_name, ordinal_position`).Scan(&snapshot.Columns).Error; err != nil {
		return fmt.Errorf("query columns: %w", err)
	}
	if err := db.Raw(`SELECT tablename AS table_name, indexname AS index_name, indexdef AS index_def FROM pg_indexes WHERE schemaname = current_schema() ORDER BY tablename, indexname`).Scan(&snapshot.Indexes).Error; err != nil {
		return fmt.Errorf("query indexes: %w", err)
	}
	if err := db.Raw(`SELECT c.conrelid::regclass::text AS table_name, c.conname AS constraint_name, c.contype::text AS constraint_type, pg_get_constraintdef(c.oid) AS definition FROM pg_constraint c JOIN pg_namespace n ON n.oid = c.connamespace WHERE n.nspname = current_schema() ORDER BY c.conrelid::regclass::text, c.conname`).Scan(&snapshot.Constraints).Error; err != nil {
		return fmt.Errorf("query constraints: %w", err)
	}
	payload, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return fmt.Errorf("encode catalog: %w", err)
	}
	if err = os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("create catalog output directory: %w", err)
	}
	if err = os.WriteFile(filepath.Join(outputDir, "postgres-catalog.json"), payload, 0o644); err != nil {
		return fmt.Errorf("write catalog snapshot: %w", err)
	}
	return nil
}
