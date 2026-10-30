# 14 — Risk register

| Risk | Severity | Control | State |
|---|---|---|---|
| API/DTO/permission omissions | Critical | 339-operation route + Web caller + handler + .NET matrix | Open |
| DB schema mismatch | Critical | 65-table column/index diff on isolated PG/SQLite | Open |
| bcrypt/AES/JWT incompatibility | Critical | Go/.NET golden vectors and old-session tests | Open |
| Payment double credit / ledger corruption | Critical | transactions, idempotency, concurrency, reconciliation | Open |
| MCP privilege escalation / SSRF | Critical | signed context, outbound policy, DNS/redirect tests | Open |
| File/object/vector loss | Critical | hashes, object keys, quotas, signatures and backfill | Open |
| Provider native protocol loss | High | 16-adapter protocol/streaming fixtures | Open |
| NDJSON buffering/cancellation | High | flush, disconnect cancellation and event equivalence | Open |
| OAuth/PKCE callback regression | High | sandbox callback/state/nonce tests | Open |
| Static routes/CORS/cookies regression | High | reverse-proxy and browser E2E | Open |
| Worker duplicate/stale writes | Critical | durable lease, idempotency and fencing | Open |
| Desktop sidecar regression | Critical | preserve Go/Cargo/Tauri and sidecar-tag CI | Guarded |
| CI absent on migration branches | Medium | branch-specific migration workflow | Added; Actions run pending |
