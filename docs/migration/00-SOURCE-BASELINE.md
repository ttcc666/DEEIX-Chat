# 00 — Source baseline

| Item | Baseline |
|---|---|
| Repository | `ttcc666/DEEIX-Chat` |
| Branch | `dev` |
| Commit | `cd58a60eeff1334015ec11ae06315a9ef91a557d` |
| Version | `0.4.6` |
| Go source files | 939 |
| Go test files | 323 |
| Swagger paths | 275 |
| Swagger operations | 339 |
| GORM models in `schema.Models()` | 65 |
| Explicit GORM table mappings | 65 |
| Provider adapter variants | 16 |
| Config env override variables | 57 |
| Web API TypeScript source files | 56 |
| Web Admin API TypeScript source files | 19 (included in 56) |
| Web page entrypoints | 35 (excluding Desktop-only route group) |

HTTP methods: GET 127, POST 125, PATCH 45, DELETE 33, PUT 9.

Primary evidence: `backend/docs/swagger.yaml`, `backend/internal/infra/persistence/schema/schema.go`, `backend/internal/infra/persistence/models/*.go`, `backend/internal/infra/llm/client.go`, `backend/internal/infra/config/config.go`, `.github/workflows/ci.yml`.

This audit uses the remote GitHub tree. No local clone was available, so local `git status --short` and user-local uncommitted changes could not be inspected; this remains a G0 blocker.
