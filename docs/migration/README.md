# DEEIX-Chat Web → .NET 10 migration

- Repository: `ttcc666/DEEIX-Chat`
- Base: `dev@cd58a60eeff1334015ec11ae06315a9ef91a557d` (version 0.4.6)
- Working branch: `migration/dotnet-web-phase0`
- Gate: **G0 BLOCKED**; continue Phase 0 before creating the .NET solution.

The branch contains the source-level API, Go route, Web caller, database field/index, custom DDL and vector-storage inventories. The canonical OpenAPI contract now covers all 343 `/api/v1` operations; all 17 operations without direct/indirect Web API callers have explicit scope dispositions. See [PHASE-0-GATE-REPORT.md](PHASE-0-GATE-REPORT.md) for verified counts, remaining evidence and the next gate criteria.
