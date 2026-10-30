# 13 — Test inventory and CI plan

- Existing Go source files: 939; Go test files: 323.
- Existing `.github/workflows/ci.yml` includes API lint/build/test, Web check/test/build, Desktop Rust/sidecar checks, Python, Docker and dependency audits.
- Existing CI push branches are `dev` and `main`; the new migration branch needs its own push-triggered workflow.

The new `.github/workflows/dotnet-web-migration.yml` runs migration artifact validation, Go default and Desktop-sidecar-tag lint/build/tests, Web checks/package tests/static build, and conditional .NET 10 restore/build/test when a solution exists under `backend-dotnet/`.

G0/G1+ must add API contract tests, isolated PostgreSQL/SQLite compatibility tests, cross-language auth/crypto vectors, Provider streaming fixtures, MCP/security tests, file/RAG/vector/Memory tests, payment idempotency/reconciliation, Web E2E and Docker build (no push). The first Phase 0 Actions run is a regression baseline, not .NET migration proof.
