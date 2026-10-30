# 05B — Database index inventory

## Source-level inventory

- [06C-DATABASE-GORM-INDEX-MEMBERS.csv](06C-DATABASE-GORM-INDEX-MEMBERS.csv) contains **493 index-tag members** across all 65 registered GORM models, including composite and unique index membership.
- [06D-POSTGRES-CUSTOM-DDL.csv](06D-POSTGRES-CUSTOM-DDL.csv) contains **87 custom PostgreSQL DDL statements**, of which 35 create/drop indexes or define the dynamic vector HNSW index. This captures partial predicates and ordered composite expressions that cannot be reconstructed from struct tags alone.
- [06E-VECTOR-STORAGE-INVENTORY.csv](06E-VECTOR-STORAGE-INVENTORY.csv) records the SQLite sqlite-vec virtual tables separately from ordinary relational indexes.

## Important custom index behavior

- LLM routing indexes include uniqueness on upstream/model binding and a partial active-route routing index.
- Billing indexes include billable usage-ledger partial indexes and an idempotency unique index for usage reserve/refund transactions.
- Conversation indexes include descending sort order, soft-delete predicates and active-share filtering.
- File-object uniqueness is partial and applies only to active, non-deleted rows with a non-empty SHA-256.
- PostgreSQL vector indexes are generated dynamically as HNSW over a padded/subvector 4000-dimension halfvec expression; they are not ordinary GORM indexes.

## Remaining verification

This is a source inventory, not proof that an existing deployment's physical indexes match. G0 still requires catalog extraction from isolated PostgreSQL and SQLite databases after migration, including index validity, predicates, ordering, operator classes, foreign keys, unique constraints and vector extension/index definitions. No production DSN may be used.

 
The `chat_conversation_projects.public_id` uniqueness rule is active-row-only. The GORM tag and PostgreSQL custom DDL now create `idx_chat_conversation_projects_public_id_active` with `WHERE deleted_at IS NULL`; both database migrations remove the legacy full unique index name after creating the replacement.
