# 05 — Database schema inventory

## Registered model and table baseline

- `backend/internal/infra/persistence/schema/schema.go` registers **65 GORM models** through `Models()`.
- All **65 explicit `TableName()` mappings** have been checked against their model source files; see [06-DATABASE-TABLE-MAPPING.csv](06-DATABASE-TABLE-MAPPING.csv).
- [06B-DATABASE-FIELD-INVENTORY.csv](06B-DATABASE-FIELD-INVENTORY.csv) contains **985 flattened field rows** across the 65 models, including inherited `BaseModel` and `ControlPlaneModel` fields.
- [06C-DATABASE-GORM-INDEX-MEMBERS.csv](06C-DATABASE-GORM-INDEX-MEMBERS.csv) records **493 GORM index-tag members**.
- [06D-POSTGRES-CUSTOM-DDL.csv](06D-POSTGRES-CUSTOM-DDL.csv) records **88 PostgreSQL custom DDL statements**, including 36 index-related statements.
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


## Physical migration catalog verification

The migration branch's isolated CI job exports the post-migration physical catalog from a unique PostgreSQL schema and a temporary SQLite file. It compares these snapshots with the source inventories in:

- `06-DATABASE-TABLE-MAPPING.csv` — expected model tables.
- `06B-DATABASE-FIELD-INVENTORY.csv` — flattened fields and GORM column/nullability/default tags.
- `06C-DATABASE-GORM-INDEX-MEMBERS.csv` — named and composite GORM index members.
- `06D-POSTGRES-CUSTOM-DDL.csv` — PostgreSQL custom columns and indexes.
- `06E-VECTOR-STORAGE-INVENTORY.csv` — sqlite-vec virtual tables and embedding dimensions.

The comparison script is `scripts/compare-db-catalogs.py`. It currently verifies table/column presence, declared NOT NULL/default semantics, primary keys, named GORM indexes, composite index order, partial-index predicates, custom PostgreSQL indexes and sqlite-vec DDL. PostgreSQL type families, foreign-key semantics, vector operator classes and soft-delete query behavior remain separate review items; a passing name/index comparison alone does not close those risks.

 
The conversation-project public ID index is intentionally partial: `idx_chat_conversation_projects_public_id_active` enforces uniqueness only where `deleted_at IS NULL`. PostgreSQL and SQLite migrations now remove the legacy full unique index after the replacement partial index is created, so soft-deleted projects do not reserve public IDs forever.
