# 12 — Web deployment analysis

- Next.js static export output: `apps/web/out`.
- Public API base: `NEXT_PUBLIC_API_BASE_URL`.
- Existing root Dockerfile builds static Web plus Go API runtime.
- Existing Compose variants cover app-only, PostgreSQL/Redis and SQLite/in-process cache.
- Existing image workflow publishes only on `dev`, `main` and version tags; migration branch must not publish images.

Target: keep Web static export; build API/Worker with .NET 10 after G0; support separated static hosting/API and same-origin proxy; preserve OAuth callbacks, refresh cookie, CORS, health/readiness/version, NDJSON flush and file streaming. API/health routes must never be swallowed by SPA fallback; CDN must not cache API. Preserve old environment variables or document compatibility mapping.

No Docker image publication, deployment, DB migration or traffic switch is part of this phase.
