# 02 — Web dependency map

- `apps/web/`: Next.js App Router / React / TypeScript Web UI.
- `apps/web/shared/api/`: shared API clients.
- `apps/web/features/admin/api/`: Admin API clients.
- `packages/api-contract/`: OpenAPI-generated TypeScript contracts.
- `packages/core/`: shared client logic; preserve non-Web consumers.
- `apps/web/next.config.ts`: static export to `apps/web/out`.
- `NEXT_PUBLIC_API_BASE_URL`: frontend API base URL.
- `packages/api-contract/scripts/generate.mjs`: current generator invokes Go `swag`; Phase 8 must move canonical OpenAPI generation to ASP.NET Core without manually editing generated types.

A per-operation join between Swagger path, actual Web caller, Go handler, DTO, auth policy and .NET endpoint remains pending; do not infer every Swagger route is Web-required.
