# 04 — API compatibility matrix

Source: `backend/docs/swagger.yaml` (Swagger 2.0, base path `/api/v1`). Detailed operation matrix: [03-API-INVENTORY.csv](03-API-INVENTORY.csv). Source call sites: [03A-WEB-API-CALLERS.csv](03A-WEB-API-CALLERS.csv). Go route registrations: [03B-GO-ROUTE-HANDLERS.csv](03B-GO-ROUTE-HANDLERS.csv). Coverage analysis: [03C-ROUTE-COVERAGE-REPORT.md](03C-ROUTE-COVERAGE-REPORT.md). Non-Web routes and migration dispositions: [03E-WEB-ROUTE-REVIEW.csv](03E-WEB-ROUTE-REVIEW.csv).

## Verified route and contract mapping

- Swagger: 278 paths / 343 operations.
- All 343 Swagger operations have an exact HTTP-method + normalized-path match to a Go router registration and handler.
- All three Web-used routes previously missing from Swagger are now in the canonical contract: `GET /billing/config`, `GET /admin/billing/plans`, and `GET /version`.
- EPay's GET and POST callbacks are both documented. The only Go registration outside the `/api/v1` contract is the root `GET /manifest.webmanifest` static-resource alias.
- The operation matrix includes request parameter summaries, **145 request-body schema references**, and **343 response-schema maps**.
- Go base authentication contexts: 169 authenticated + admin-only, 149 authenticated, 23 public/rate-limited, 1 local-mode-only and 1 public without auth/rate-limit middleware.
- 329 Web caller source rows: 324 direct path matches and 2 indirect browser-image callers. The remaining 17 operations have explicit non-Web dispositions and source evidence in `03E-WEB-ROUTE-REVIEW.csv`.
- Web matching is deliberately marked with direct, indirect or reviewed non-Web classifications; caller rows record source file and line, but method-level AST verification remains a follow-up.

## Operation-level completion requirements

For each operation, preserve route registration, Handler, authentication and fine-grained permissions, request/response DTO, required/optional/nullable semantics, status/error envelope, Cookie/Header, streaming or file behavior, actual Web caller, generated TypeScript types, .NET endpoint/use case and tests.

The matrix's `go_auth_context` is the base router middleware context, not a complete permission decision. Handler-level authorization, resource ownership and permission-group checks still require review. Swagger schema references are contract evidence, not proof that the current Go handler returns identical payloads in every branch.

The API contract generator runs from Go annotations and updates `backend/docs/{docs.go,swagger.json,swagger.yaml}` and `packages/api-contract/src/types.generated.ts` together. Generated artifacts are verified by `pnpm api:check`.

The inventory records the Go handler and source router. `dotnet_implementation=NOT_STARTED` and `compatibility_status=NOT_MAPPED` are intentional; route/schema mapping is not .NET implementation or compatibility proof.
