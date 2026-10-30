package sqlite

import (
	"path/filepath"
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
