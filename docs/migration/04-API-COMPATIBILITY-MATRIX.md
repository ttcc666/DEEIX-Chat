# 04 — API compatibility matrix

Source: `backend/docs/swagger.yaml` (Swagger 2.0, base path `/api/v1`). Detailed operation matrix: [03-API-INVENTORY.csv](03-API-INVENTORY.csv). Source call sites: [03A-WEB-API-CALLERS.csv](03A-WEB-API-CALLERS.csv). Go route registrations: [03B-GO-ROUTE-HANDLERS.csv](03B-GO-ROUTE-HANDLERS.csv). Coverage analysis: [03C-ROUTE-COVERAGE-REPORT.md](03C-ROUTE-COVERAGE-REPORT.md).

## Verified route mapping

- Swagger: 275 paths / 339 operations.
- Go route registrations: 344.
- All 339 Swagger operations have an exact HTTP-method + normalized-path match to a Go router registration and handler.
- Five Go routes are not represented as Swagger operations; they are listed in the coverage report for explicit inclusion/exclusion decisions.
- 327 Web API call-site rows were inventoried. 321 Swagger operations have a direct Web path-literal match; 18 remain in a review bucket because they are callback, browser/static, indirect, Desktop-only or currently unused routes.
- Web matching is deliberately marked `PATH_MATCH_ONLY`; the caller CSV records source file and line, but method-level AST verification remains a follow-up.

## Operation-level completion requirements

For each operation, preserve route registration, Handler, authentication and permissions, request/response DTO, required/optional/nullable semantics, status/error envelope, Cookie/Header, streaming or file behavior, actual Web caller, generated TypeScript types, .NET endpoint/use case and tests.

The inventory records the Go handler and source router. `dotnet_implementation=NOT_STARTED` and `compatibility_status=NOT_MAPPED` are intentional; route mapping is not .NET implementation or compatibility proof.
