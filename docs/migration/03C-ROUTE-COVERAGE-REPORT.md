# 03C — API, Go route and Web caller coverage

## Counts

| Inventory | Count |
|---|---:|
| Swagger paths | 275 |
| Swagger operations | 339 |
| Go route registrations | 344 |
| Swagger operations with exact Go method/path/handler mapping | 339 |
| Web API caller source rows | 327 |
| Operations with direct Web path-literal match | 321 |
| Operations without direct Web path-literal match | 18 |
| Go route registrations not represented in Swagger | 5 |

The detailed join is in [03-API-INVENTORY.csv](03-API-INVENTORY.csv). The source call sites are in [03A-WEB-API-CALLERS.csv](03A-WEB-API-CALLERS.csv); Go registrations are in [03B-GO-ROUTE-HANDLERS.csv](03B-GO-ROUTE-HANDLERS.csv).

## Five Go registrations not represented in Swagger

| Method | Route | Handler | Review |
|---|---|---|---|
| GET | `/api/v1/billing/payments/epay/notify` | `EPayNotify` | Determine whether this is a legacy/provider callback that should be documented or retired |
| GET | `/api/v1/billing/config` | `GetBillingConfig` | Decide whether the Web-facing billing config contract is missing from Swagger |
| GET | `/api/v1/admin/billing/plans` | `ListPlans` | Decide whether the Admin billing plans read contract is missing from Swagger |
| GET | `/api/v1/version` | `buildinfo.Snapshot` | Web calls this endpoint; add to canonical OpenAPI contract |
| GET | `/manifest.webmanifest` | `GetBrandingManifest` | Browser manifest route outside `/api/v1`; document as static Web behavior, not an API operation |

## Eighteen Swagger operations without a direct path-literal match in the inventoried Web API modules

| Group | Routes | Initial interpretation; confirm before excluding |
|---|---|---|
| Admin runtime controls (6) | RapidOCR and Tika start/stop/restart | No direct caller found in current Admin API modules |
| Desktop/local auth (1) | `POST /auth/local/exchange` | Likely Desktop-only; explicitly excluded from Web migration unless evidence says otherwise |
| OIDC callback (1) | `GET /auth/providers/{slug}/callback` | Browser/provider redirect callback, not an ordinary API client call |
| Payment callbacks (2) | EPay notify and Stripe webhook | Provider-to-server callbacks, not browser callers |
| Branding manifest (1) | `GET /branding/manifest.webmanifest` | Browser/static resource path |
| File processing (1) | `GET /files/{file_id}/processing` | Check whether replaced by batch processing-status APIs |
| Model icon asset (1) | `GET /llm/icon-assets/{public_id}` | May be loaded as an image URL rather than through an API module |
| User UI components (4) | `/ui-components/mine` GET/POST/PATCH/DELETE | No direct caller found in current Web API module; confirm whether feature is unused |
| User avatar (1) | `GET /users/{public_id}/avatar` | May be loaded as an image URL; search Web image/avatar components |

“Without a direct path-literal match” is not equivalent to “unused”. Browser resources, dynamically assembled URLs, redirects and Desktop-only flows need an explicit decision.

## Mapping limitations

- All Swagger operations now include their Go handler and router source file.
- Web matching is path-only: the caller inventory records file/line, but this pass does not claim AST-verified HTTP-method equivalence.
- Request/response DTO, permission middleware and error/streaming contract mapping remain G0 work.
- The API CSV intentionally keeps .NET implementation and compatibility status as `NOT_STARTED` / `NOT_MAPPED`.
