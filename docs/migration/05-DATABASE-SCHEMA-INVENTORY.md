# 05 — Database schema inventory

- `backend/internal/infra/persistence/schema/schema.go` registers 65 GORM models through `Models()`.
- Table names are declared by model `TableName()` methods; mapping: [06-DATABASE-TABLE-MAPPING.csv](06-DATABASE-TABLE-MAPPING.csv).
- `schema.Migrate()` creates missing tables, runs `AutoMigrate`, invalidates unsigned file embeddings, backfills context artifact message IDs and usage ledger billing timestamps.
- `CleanupRemovedColumns()` drops legacy columns in prompt presets, skills and billing plans. This is destructive and must not be ported blindly.
- PostgreSQL adds custom baseline migrations, indexes/constraints and vector baseline; SQLite uses separate sqlite-vec migration.

## Vector invariants
- Maximum embedding width: 4096.
- PostgreSQL stores native dimensions; 4000-dimension halfvec HNSW candidates are re-ranked with full vectors.
- SQLite sqlite-vec uses 4096-dimensional fixed slots and pads only at the persistence boundary.
- Vector virtual tables: `file_chunk_vectors`, `chat_message_chunk_vectors`, `user_memory_vectors`.

G0 still requires a column/type/nullability/default/PK/FK/index/custom-SQL comparison against isolated PostgreSQL and SQLite databases. No production DSN may be used.
