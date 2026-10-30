# 13 — Test inventory and CI plan

## Existing baseline

- Existing Go source files: 939; Go test files: 323.
- Existing `.github/workflows/ci.yml` includes API lint/build/test, Web check/test/build, Desktop Rust/sidecar checks, Python, Docker and dependency audits.
- The migration workflow is branch-scoped and does not publish images or deploy.

## Migration branch workflow

`.github/workflows/dotnet-web-migration.yml` runs:

1. Inventory validation: 339 Swagger operations, 344 Go route registrations, 327 Web caller rows, 65 explicit GORM table mappings, 886 flattened model fields, 430 GORM index members, 87 PostgreSQL custom DDL statements and 3 SQLite vector tables.
2. Go lint, default build/test and Desktop-sidecar-tag build/test.
3. Web dependency install, lint/type/architecture/API-contract checks, package tests and static build.
4. PostgreSQL persistence integration tests against an ephemeral pgvector service with `DEEIX_TEST_DATABASE_DSN`.
5. Conditional .NET 10 restore/build/test when a solution exists under `backend-dotnet/`.

The previous baseline Run `41006537076` passed audit, Go and Web checks. The PostgreSQL integration job is newly added and must complete on the next run.

## Remaining test work

- Export and compare physical PostgreSQL and SQLite catalogs after migration; no production DSN.
- Add explicit API contract tests for body/response schemas, permissions, cookies, streaming and error envelopes.
- Add cross-language bcrypt/AES/JWT golden vectors in Phase 2.
- Add Provider protocol/streaming fixtures in Phase 3.
- Add MCP signed-context/SSRF/replay tests in Phase 5.
- Add file/OCR/vector/RAG/Memory lifecycle tests in Phase 6.
- Add payment idempotency/reconciliation/concurrency tests in Phase 7.
- Add Web E2E and Docker build (no push) in Phase 8.
