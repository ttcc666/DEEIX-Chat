package postgres

import (
	"fmt"
	"net/url"
	"os"
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
}
