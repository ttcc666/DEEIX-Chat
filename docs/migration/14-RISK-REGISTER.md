# 14 — Risk register

| Risk | Severity | Control | State |
|---|---|---|---|
| API/DTO/permission omissions | Critical | 339 Swagger operations now map to Go handlers; request/response/auth policy review remains | Partially mapped |
| Web endpoint caller omissions | High | 321 direct path matches; 18 no-literal routes classified for manual review | Open |
| Undocumented Go routes | Medium | Five route registrations absent from Swagger require include/exclude decisions | Open |
| DB schema mismatch | Critical | 65 explicit tables, 886 fields, 430 index members and 87 custom DDL statements inventoried; physical catalog diff pending | Open |
| bcrypt/AES/JWT incompatibility | Critical | Go/.NET golden vectors and old-session tests in Phase 2 | Open |
| Payment double credit / ledger corruption | Critical | transactions, idempotency, concurrency, reconciliation | Open |
| MCP privilege escalation / SSRF | Critical | signed context, outbound policy, DNS/redirect tests | Open |
| File/object/vector loss | Critical | hashes, object keys, quotas, signatures and backfill | Open |
| Provider native protocol loss | High | 16-adapter protocol/streaming fixtures | Open |
| NDJSON buffering/cancellation | High | flush, disconnect cancellation and event equivalence | Open |
| OAuth/PKCE callback regression | High | sandbox callback/state/nonce tests | Open |
| Static routes/CORS/cookies regression | High | reverse-proxy and browser E2E | Open |
| Worker duplicate/stale writes | Critical | durable lease, idempotency and fencing | Open |
| Desktop sidecar regression | Critical | preserve Go/Cargo/Tauri and sidecar-tag CI | Guarded |
| PostgreSQL integration regressions | High | ephemeral pgvector Actions job | Added; run pending |
| CI absent on migration branches | Medium | branch-specific migration workflow | Added |
