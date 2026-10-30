package sqlite

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/DEEIX-AI/DEEIX-Chat/backend/internal/infra/config"
	persistenceSchema "github.com/DEEIX-AI/DEEIX-Chat/backend/internal/infra/persistence/schema"
)

func TestNewMigratesAllRegisteredModelsAndVectorTables(t *testing.T) {
	db, err := New(config.Config{
		SQLitePath:          filepath.Join(t.TempDir(), "migration-catalog.sqlite"),
		SQLiteBusyTimeoutMS: 5000,
		SQLiteCacheSizeKB:   2048,
		SQLiteSynchronous:   "NORMAL",
		SQLiteTempStore:     "MEMORY",
		SQLiteMaxOpenConns:  1,
	})
	if err != nil {
		t.Fatalf("run full SQLite migration: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("resolve SQLite connection: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })

	models := persistenceSchema.Models()
	if len(models) != 65 {
		t.Fatalf("registered model count = %d, want 65", len(models))
	}
	for _, model := range models {
		if !db.Migrator().HasTable(model) {
			t.Errorf("registered model table is missing: %T", model)
		}
		columns, columnErr := db.Migrator().ColumnTypes(model)
		if columnErr != nil {
			t.Errorf("inspect columns for %T: %v", model, columnErr)
		} else if len(columns) == 0 {
			t.Errorf("registered model has no columns: %T", model)
		}
	}

	vectorTables := []string{
		"file_chunk_vectors",
		"chat_message_chunk_vectors",
		"user_memory_vectors",
	}
	for _, table := range vectorTables {
		var ddl string
		if err = db.Raw("SELECT sql FROM sqlite_master WHERE type = 'table' AND name = ?", table).Scan(&ddl).Error; err != nil {
			t.Errorf("read vector table %s: %v", table, err)
			continue
		}
		if ddl == "" {
			t.Errorf("vector virtual table is missing: %s", table)
			continue
		}
		if !containsSQLiteVec(ddl) {
			t.Errorf("vector table %s is not backed by sqlite-vec: %s", table, ddl)
		}
	}

	if err = writeSQLiteCatalogSnapshot(db); err != nil {
		t.Fatalf("export SQLite migration catalog: %v", err)
	}
}

func containsSQLiteVec(ddl string) bool {
	const marker = "vec0"
	for i := 0; i+len(marker) <= len(ddl); i++ {
		if ddl[i:i+len(marker)] == marker {
			return true
		}
	}
	return false
}

type sqliteCatalogSnapshot struct {
	Tables  []string                  `json:"tables"`
	Objects []sqliteCatalogObject     `json:"objects"`
	Columns []sqliteCatalogColumn     `json:"columns"`
	Indexes []sqliteCatalogIndex      `json:"indexes"`
}

type sqliteCatalogObject struct {
	Type      string `json:"type" gorm:"column:type"`
	Name      string `json:"name" gorm:"column:name"`
	TableName string `json:"table_name" gorm:"column:table_name"`
	SQL       string `json:"sql" gorm:"column:sql"`
}

type sqliteCatalogColumn struct {
	TableName    string `json:"table_name"`
	CID          int    `json:"cid" gorm:"column:cid"`
	Name         string `json:"name" gorm:"column:name"`
	Type         string `json:"type" gorm:"column:type"`
	NotNull      int    `json:"not_null" gorm:"column:notnull"`
	DefaultValue string `json:"default_value" gorm:"column:dflt_value"`
	PrimaryKey   int    `json:"primary_key" gorm:"column:pk"`
}

type sqliteCatalogIndex struct {
	TableName string   `json:"table_name"`
	Seq       int      `json:"seq" gorm:"column:seq"`
	Name      string   `json:"name" gorm:"column:name"`
	Unique    int      `json:"unique" gorm:"column:unique"`
	Origin    string   `json:"origin" gorm:"column:origin"`
	Partial   int      `json:"partial" gorm:"column:partial"`
	Columns   []string `json:"columns"`
	SQL       string   `json:"sql"`
}

type sqliteIndexColumn struct {
	SeqNo int    `gorm:"column:seqno"`
	CID   int    `gorm:"column:cid"`
	Name  string `gorm:"column:name"`
}

func writeSQLiteCatalogSnapshot(db *gorm.DB) error {
	outputDir := strings.TrimSpace(os.Getenv("DEEIX_TEST_CATALOG_DIR"))
	if outputDir == "" {
		return nil
	}
	snapshot := sqliteCatalogSnapshot{}
	if err := db.Raw(`SELECT type, name, tbl_name AS table_name, COALESCE(sql, '') AS sql FROM sqlite_master WHERE type IN ('table', 'index') AND name NOT LIKE 'sqlite_%' ORDER BY type, name`).Scan(&snapshot.Objects).Error; err != nil {
		return fmt.Errorf("query sqlite_master: %w", err)
	}
	objectSQL := make(map[string]string, len(snapshot.Objects))
	for _, object := range snapshot.Objects {
		objectSQL[object.Name] = object.SQL
		if object.Type != "table" {
			continue
		}
		snapshot.Tables = append(snapshot.Tables, object.Name)
		quotedTable := quoteSQLiteIdentifier(object.Name)
		var columns []sqliteCatalogColumn
		if err := db.Raw("PRAGMA table_info(" + quotedTable + ")").Scan(&columns).Error; err != nil {
			return fmt.Errorf("query columns for %s: %w", object.Name, err)
		}
		for i := range columns {
			columns[i].TableName = object.Name
			snapshot.Columns = append(snapshot.Columns, columns[i])
		}
		var indexes []sqliteCatalogIndex
		if err := db.Raw("PRAGMA index_list(" + quotedTable + ")").Scan(&indexes).Error; err != nil {
			return fmt.Errorf("query indexes for %s: %w", object.Name, err)
		}
		for i := range indexes {
			indexes[i].TableName = object.Name
			indexes[i].SQL = objectSQL[indexes[i].Name]
			var parts []sqliteIndexColumn
			if err := db.Raw("PRAGMA index_info(" + quoteSQLiteIdentifier(indexes[i].Name) + ")").Scan(&parts).Error; err != nil {
				return fmt.Errorf("query index columns for %s: %w", indexes[i].Name, err)
			}
			for _, part := range parts {
				indexes[i].Columns = append(indexes[i].Columns, part.Name)
			}
			snapshot.Indexes = append(snapshot.Indexes, indexes[i])
		}
	}
	payload, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return fmt.Errorf("encode catalog: %w", err)
	}
	if err = os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("create catalog output directory: %w", err)
	}
	if err = os.WriteFile(filepath.Join(outputDir, "sqlite-catalog.json"), payload, 0o644); err != nil {
		return fmt.Errorf("write catalog snapshot: %w", err)
	}
	return nil
}

func quoteSQLiteIdentifier(value string) string {
	return "\"" + strings.ReplaceAll(value, "\"", "\"\"") + "\""
}
