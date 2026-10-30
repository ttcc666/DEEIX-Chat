# 03C — API, Go route and Web caller coverage

## Counts

| Inventory | Count |
|---|---:|
| Swagger paths | 275 |
| Swagger operations | 339 |
| Go route registrations | 344 |
| Swagger operations with exact Go method/path/handler mapping | 339 |
| Request-body schema references | 145 |
| Response-schema maps | 339 |
| Go authenticated + admin-only base contexts | 168 |
| Go authenticated base contexts | 148 |
| Go public/rate-limited base contexts | 22 |
| Go local-mode-only contexts | 1 |
| Web API caller source rows | 329 |
| Operations with direct Web path-literal match | 321 |
| Operations with indirect browser-image caller | 2 |
| Operations without direct or indirect Web caller evidence | 16 |
| Go route registrations not represented in Swagger | 5 |
| Web-used Go routes missing from Swagger | 3 |

The detailed join is in [03-API-INVENTORY.csv](03-API-INVENTORY.csv). The source call sites are in [03A-WEB-API-CALLERS.csv](03A-WEB-API-CALLERS.csv); Go registrations are in [03B-GO-ROUTE-HANDLERS.csv](03B-GO-ROUTE-HANDLERS.csv). The three Web-used undocumented operations are listed in [03D-WEB-ONLY-UNSWAGGERED-ROUTES.csv](03D-WEB-ONLY-UNSWAGGERED-ROUTES.csv).

The Go auth context is the base middleware inherited from the server router: public routes use the public auth/rate-limit group, authenticated routes use JWT authentication and rate limiting, and admin routes additionally use `AdminOnly`. This does not replace handler-level ownership or permission-group checks.

## Five Go registrations not represented in Swagger

| Method | Route | Handler | Web caller | Review |
|---|---|---|---|---|
| GET | `/api/v1/billing/payments/epay/notify` | `EPayNotify` | Provider callback; no Web API caller | Determine whether to document as provider callback or retire |
| GET | `/api/v1/billing/config` | `GetBillingConfig` | `apps/web/shared/api/billing.ts:19` | Add to canonical OpenAPI contract |
| GET | `/api/v1/admin/billing/plans` | `ListPlans` | `apps/web/features/admin/api/billing.ts:41` | Add to canonical OpenAPI contract |
| GET | `/api/v1/version` | `buildinfo.Snapshot` | `apps/web/shared/api/app-version.ts:14` | Add to canonical OpenAPI contract |
| GET | `/manifest.webmanifest` | `GetBrandingManifest` | Browser manifest route; no API module caller | Document as static Web behavior outside `/api/v1` |

## Eighteen Swagger operations without a direct path-literal match in the inventoried Web API modules

| Group | Routes | Caller evidence / initial interpretation |
|---|---|---|
| Admin runtime controls (6) | RapidOCR and Tika start/stop/restart | No direct caller found in current Admin API modules |
| Desktop/local auth (1) | `POST /auth/local/exchange` | Likely Desktop-only; explicitly excluded from Web migration unless evidence says otherwise |
| OIDC callback (1) | `GET /auth/providers/{slug}/callback` | Browser/provider redirect callback, not an ordinary API client call |
| Payment callbacks (2) | EPay notify and Stripe webhook | Provider-to-server callbacks, not browser callers |
| Branding manifest (1) | `GET /branding/manifest.webmanifest` | Browser/static resource path |
| File processing (1) | `GET /files/{file_id}/processing` | No direct Web caller; check whether replaced by batch processing-status APIs |
| Model icon asset (1) | `GET /llm/icon-assets/{public_id}` | Indirect browser image caller confirmed: `apps/web/entities/model/components/model-icon.tsx:10` |
| User UI components (4) | `/ui-components/mine` GET/POST/PATCH/DELETE | No direct caller found in current Web API module; confirm whether feature is unused |
| User avatar (1) | `GET /users/{public_id}/avatar` | Indirect browser image caller confirmed: `apps/web/shared/lib/avatar.ts:185` |

“Without a direct path-literal match” is not equivalent to “unused”. Two of the 18 have confirmed indirect browser-image callers; the remaining 16 still need an explicit decision.

## Mapping limitations

- All Swagger operations include their Go handler and router source file.
- The operation matrix now records request parameters, body schema references, response schema references and base auth context.
- Web matching is path-only: the caller inventory records file/line, but this pass does not claim AST-verified HTTP-method equivalence.
- Handler-level DTO equivalence, fine-grained permissions, error/streaming contract verification and .NET implementation remain open.
- The API CSV intentionally keeps .NET implementation and compatibility status as `NOT_STARTED` / `NOT_MAPPED`.
