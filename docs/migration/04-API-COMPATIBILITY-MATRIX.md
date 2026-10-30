# 04 — API compatibility matrix

Source: `backend/docs/swagger.yaml` (Swagger 2.0, base path `/api/v1`). Detailed operation inventory: [03-API-INVENTORY.csv](03-API-INVENTORY.csv).

- Paths: 275.
- Operations (method + path): 339.
- Methods: GET 127, POST 125, PATCH 45, DELETE 33, PUT 9.

For each operation, map route registration, Handler, authentication/permissions, request/response DTO, required/optional/nullable semantics, status/error envelope, Cookie/Header, streaming/file behavior, actual Web caller, generated TS types, .NET endpoint/use case and tests.

All inventory rows are intentionally marked `web_caller_mapping=PENDING`, `dotnet_implementation=NOT_STARTED`, `compatibility_status=NOT_MAPPED`. This is a static inventory, not an implementation claim.
