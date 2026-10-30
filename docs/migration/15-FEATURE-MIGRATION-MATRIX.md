# 15 — Feature migration matrix

| Phase | Area | Status |
|---|---|---|
| 0 | Source/API/DB/Auth/Provider/MCP/Files/Billing/Web/CI audit | In progress — G0 BLOCKED |
| 1 | .NET solution/config/logging/DB/cache/storage/OpenAPI/health | Not started; blocked by G0 |
| 2 | Identity/sessions/MFA/OAuth/authorization | Not started |
| 3 | Providers/models/routes/circuit breakers | Not started |
| 4 | Conversations/context/NDJSON | Not started |
| 5 | MCP and native tools | Not started |
| 6 | Files/OCR/embeddings/vector/RAG/Memory | Not started |
| 7 | Billing/payments/admin/audit/operations | Not started |
| 8 | OpenAPI/TS contract/Web static deployment/Docker/CI | Not started |
| 9 | Shadow compare/cutover/Go Web API retirement | Not started; explicit approval required |

A phase is complete only when its Gate evidence and tests pass. Inventory-only commits do not mean the feature is migrated.
