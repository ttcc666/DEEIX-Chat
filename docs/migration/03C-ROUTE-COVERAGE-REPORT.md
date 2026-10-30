# 03C — API, Go route and Web caller coverage

## Counts

| Inventory | Count |
|---|---:|
| Swagger paths | 278 |
| Swagger operations | 343 |
| Go route registrations | 344 |
| Swagger operations with exact Go method/path/handler mapping | 343 |
| Request-body schema references | 145 |
| Response-schema maps | 343 |
| Go authenticated + admin-only base contexts | 169 |
| Go authenticated base contexts | 149 |
| Go public/rate-limited base contexts | 23 |
| Go local-mode-only contexts | 1 |
| Go public/no-auth/no-rate-limit contexts | 1 |
| Web API caller source rows | 329 |
| Operations with direct Web path-literal match | 324 |
| Operations with indirect browser-image caller | 2 |
| Operations with an explicit non-Web caller disposition | 17 |
| Web-used Go routes missing from Swagger | 0 |
| Go registrations outside the API contract | 1 |

The detailed operation join is in [03-API-INVENTORY.csv](03-API-INVENTORY.csv). Source call sites are in [03A-WEB-API-CALLERS.csv](03A-WEB-API-CALLERS.csv); Go registrations are in [03B-GO-ROUTE-HANDLERS.csv](03B-GO-ROUTE-HANDLERS.csv). The 17 operations without a direct or indirect Web API caller are explicitly classified, with evidence and migration disposition, in [03E-WEB-ROUTE-REVIEW.csv](03E-WEB-ROUTE-REVIEW.csv). [03D-WEB-ONLY-UNSWAGGERED-ROUTES.csv](03D-WEB-ONLY-UNSWAGGERED-ROUTES.csv) is now header-only: all three Web-used undocumented operations have been added to the canonical OpenAPI contract.

The Go auth context is the base middleware inherited from the server router: public routes use the public auth/rate-limit group, authenticated routes use JWT authentication and rate limiting, and admin routes additionally use `AdminOnly`. `GET /version` is registered directly on the API group before either auth group, so it is classified as public with no auth/rate-limit middleware. This does not replace handler-level ownership or permission-group checks.

## Route inventory reconciliation

- The three Web-used routes previously missing from Swagger are now documented: `GET /billing/config`, `GET /admin/billing/plans` and `GET /version`.
- The EPay GET notification callback is now documented separately from its POST callback. Both use the same payment notification validation and completion logic.
- All 343 Swagger operations map to a Go method/path and handler.
- The only Go registration outside the `/api/v1` contract is `GET /manifest.webmanifest`. It is a browser-facing static manifest route that reuses the handler documented at `GET /branding/manifest.webmanifest`; preserve the root route in the Web/static hosting layer.

## Seventeen operations without direct or indirect Web API callers

| Group | Count | Disposition |
|---|---:|---|
| RapidOCR and Tika runtime start/stop/restart | 6 | Web exposes runtime status probes and a test/refresh action only; exclude these action routes from the Web migration, but retain until external consumers are checked |
| Local auth exchange | 1 | Desktop/sidecar-only; excluded from Web scope |
| OIDC provider callback | 1 | Browser redirect callback; preserve 302 and state-validation semantics |
| Payment provider callbacks | 3 | EPay GET/POST and Stripe webhook; preserve provider protocols, not browser API calls |
| Branding manifest API alias | 1 | Static resource alias; preserve the root `/manifest.webmanifest` route |
| Single-file processing status | 1 | Web uses the batch POST endpoint; keep the single-file route until external consumers are checked |
| User-owned UI component routes | 4 | No current Web client wrapper/callsite; keep out of Web scope pending external-consumer review |

The model-icon and user-avatar routes are not in this no-caller group: both have confirmed indirect browser-image callers and are classified as `INDIRECT_IMAGE_MATCH`.

## Mapping limitations

- All Swagger operations include their Go handler and router source file.
- The operation matrix records request parameters, body schema references, response schema references and base auth context.
- Web matching remains path-based: caller rows record source file/line, but this pass does not claim AST-verified HTTP-method equivalence for every call.
- Handler-level DTO equivalence, fine-grained permissions, error/streaming contract verification and .NET implementation remain open.
- The API CSV intentionally keeps .NET implementation and compatibility status as `NOT_STARTED` / `NOT_MAPPED`.
