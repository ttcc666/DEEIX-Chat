# 05 — Database schema inventory

## Registered model and table baseline

- `backend/internal/infra/persistence/schema/schema.go` registers **65 GORM models** through `Models()`.
- All **65 explicit `TableName()` mappings** have been checked against their model source files; see [06-DATABASE-TABLE-MAPPING.csv](06-DATABASE-TABLE-MAPPING.csv).
- [06B-DATABASE-FIELD-INVENTORY.csv](06B-DATABASE-FIELD-INVENTORY.csv) contains **886 flattened field rows** across the 65 models, including inherited `BaseModel` and `ControlPlaneModel` fields.
- [06C-DATABASE-GORM-INDEX-MEMBERS.csv](06C-DATABASE-GORM-INDEX-MEMBERS.csv) records **430 GORM index-tag members**.
- [06D-POSTGRES-CUSTOM-DDL.csv](06D-POSTGRES-CUSTOM-DDL.csv) records **87 PostgreSQL custom DDL statements**, including 35 index-related statements.
- [06E-VECTOR-STORAGE-INVENTORY.csv](06E-VECTOR-STORAGE-INVENTORY.csv) records the three SQLite sqlite-vec virtual tables.

## Migration behavior and destructive operations

- `schema.Migrate()` creates missing tables, runs `AutoMigrate`, invalidates unsigned file embeddings, backfills context artifact message IDs and backfills usage-ledger billing timestamps.
- PostgreSQL additionally applies custom baseline DDL, partial/composite indexes, comments, pgvector extension/index maintenance and identity/session/conversation/billing/announcement upgrades in `backend/internal/infra/persistence/postgres/postgres.go`.
- SQLite calls `schema.CleanupRemovedColumns()` and creates/updates sqlite-vec virtual tables in `backend/internal/infra/persistence/sqlitevec/sqlitevec.go`.
- `CleanupRemovedColumns()` drops legacy columns in prompt presets, skills and billing plans. This is destructive and must not be ported blindly.

## Vector invariants

- Maximum embedding width: 4096.
- PostgreSQL stores native dimensions; candidate retrieval uses 4000-dimensional halfvec HNSW and re-ranks with full vectors.
- SQLite sqlite-vec uses 4096-dimensional fixed slots and pads only at the persistence boundary.
- Vector virtual tables: `file_chunk_vectors`, `chat_message_chunk_vectors`, `user_memory_vectors`.

## G0 database work still required

The inventories above are source-derived, not a live database catalog comparison. Run migrations against isolated PostgreSQL (with pgvector) and SQLite databases, export actual columns/types/nullability/defaults/PK/FK/indexes, and compare the resulting catalogs with these source inventories. Do not use a production DSN. The migration branch CI now has a PostgreSQL integration job; its result is required before closing the database portion of G0.
