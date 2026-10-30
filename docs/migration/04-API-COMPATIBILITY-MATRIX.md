# 04 — API compatibility matrix

Source: `backend/docs/swagger.yaml` (Swagger 2.0, base path `/api/v1`). Detailed operation matrix: [03-API-INVENTORY.csv](03-API-INVENTORY.csv). Source call sites: [03A-WEB-API-CALLERS.csv](03A-WEB-API-CALLERS.csv). Go route registrations: [03B-GO-ROUTE-HANDLERS.csv](03B-GO-ROUTE-HANDLERS.csv). Coverage analysis: [03C-ROUTE-COVERAGE-REPORT.md](03C-ROUTE-COVERAGE-REPORT.md).

## Verified route and contract mapping

- Swagger: 275 paths / 339 operations.
- All 339 Swagger operations have an exact HTTP-method + normalized-path match to a Go router registration and handler.
- Five Go route registrations are not represented as Swagger operations; see the coverage report for inclusion/exclusion decisions.
- The operation matrix includes request parameter summaries, **145 request-body schema references**, and **339 response-schema maps**.
- Go base authentication contexts are mapped from the actual router registration chain: 168 authenticated + admin-only, 148 authenticated, 22 public/rate-limited and 1 local-mode-only.
- 327 Web API call-site rows were inventoried. 321 Swagger operations have a direct Web path-literal match; 18 remain in a manual review bucket.
- Web matching is deliberately marked `PATH_MATCH_ONLY`; the caller CSV records source file and line, but method-level AST verification remains a follow-up.

## Operation-level completion requirements

For each operation, preserve route registration, Handler, authentication and fine-grained permissions, request/response DTO, required/optional/nullable semantics, status/error envelope, Cookie/Header, streaming or file behavior, actual Web caller, generated TypeScript types, .NET endpoint/use case and tests.

The matrix's `go_auth_context` is the base router middleware context, not a complete permission decision. Handler-level authorization, resource ownership and permission-group checks still require review. Swagger schema references are contract evidence, not proof that the current Go handler returns identical payloads in every branch.

The inventory records the Go handler and source router. `dotnet_implementation=NOT_STARTED` and `compatibility_status=NOT_MAPPED` are intentional; route/schema mapping is not .NET implementation or compatibility proof.
