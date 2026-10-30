# Phase 0 Gate report

Repository: `ttcc666/DEEIX-Chat` · base `dev@cd58a60eeff1334015ec11ae06315a9ef91a557d` · branch `migration/dotnet-web-phase0`.

## Decision: G0 BLOCKED — continue Phase 0

API route coverage is now reconciled and all currently uncalled Web routes have explicit dispositions. This is not yet sufficient to start Phase 1 because physical database catalog comparison and fine-grained authorization/DTO compatibility review remain open.

## Verified source baseline

- Swagger: 278 paths / 343 operations.
- Go routing: 344 registrations; all 343 Swagger operations map to an exact Go HTTP method/path and handler. The only route outside the API contract is the root `GET /manifest.webmanifest` static-resource alias.
- API contracts: 145 request-body schema references and 343 response-schema maps extracted from Swagger.
- Go base auth context: 169 authenticated + admin-only, 149 authenticated, 23 public/rate-limited, 1 local-mode-only and 1 public without auth/rate-limit middleware.
- All three Web-used routes previously missing from Swagger are now documented: `/billing/config`, `/admin/billing/plans` and `/version`. EPay GET and POST callbacks are both documented.
- Web callers: 329 source rows; 324 direct path-literal matches and 2 indirect browser-image callers. The remaining 17 operations without direct/indirect Web API callers are classified with evidence and migration dispositions in [03E-WEB-ROUTE-REVIEW.csv](03E-WEB-ROUTE-REVIEW.csv). No Web-used Go route remains missing from Swagger.
- GORM: 65 registered models; all 65 explicit `TableName()` values verified against source.
- Database source inventory: 985 flattened field rows, 493 GORM index members, 88 PostgreSQL custom DDL statements and 3 SQLite sqlite-vec tables.
- Go: 939 source files / 323 test files.
- Provider adapter variants: 16.

## Remaining G0 blockers

1. Complete the new isolated PostgreSQL/SQLite catalog export and comparison in CI. The job now exports physical columns, indexes, constraints and SQLite virtual-table DDL, then compares table/column presence, declared nullability/defaults, PKs, GORM and custom indexes, composite order and partial predicates against the checked-in inventories. G0 remains blocked until this comparison finishes with no unexplained mismatches; data types, FK/operator-class and soft-delete semantics still need explicit review.
2. Compare Swagger schema references with actual Go handler DTOs and generated TypeScript contracts; map fine-grained permission checks, ownership rules, permission-group checks, feature gates, Cookie/Header behavior, status/error envelopes and streaming/file semantics for all 343 operations. Base router middleware is mapped, but not every handler authorization branch.
3. Review external-client ownership before omitting the 6 RapidOCR/Tika lifecycle action routes and 4 user-owned UI-component routes from the Web-only .NET scope. The local auth exchange is registered only in LocalMode and remains a Desktop/sidecar exclusion. Preserve the OIDC browser redirect, EPay/Stripe provider callbacks, root manifest route and payment callback protocols as recorded in 03E.
4. Preserve Go/.NET cross-language auth/crypto compatibility tests for Phase 2; no .NET implementation exists yet.

The absence of a local clone is a visibility limitation only; this work is scoped to the remote GitHub `dev` baseline and isolated migration branch. No claim is made about unstaged changes in any developer's local working tree.

## Safety

No production DSN, production payment endpoint or production deployment was used. No Go/Web API business behavior was changed beyond equivalent handler wrappers for route documentation. A targeted persistence migration correction now makes conversation-project public ID uniqueness apply only to non-deleted rows, matching the existing PostgreSQL DDL intent. No Desktop/Tauri/Rust/Cargo/sidecar behavior was modified. Persistence tests use only ephemeral PostgreSQL and temporary SQLite databases. No destructive operation was run against a real database. No image publishing is configured for this branch.

## Next steps

Complete the physical PostgreSQL/SQLite catalog export and diff, finish handler-level permission and DTO mapping, and confirm external consumers for the routes excluded from Web scope. Re-evaluate G0 only when those checks have evidence. Then start Phase 1 in a separate commit.
