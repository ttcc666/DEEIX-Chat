# 09 — MCP compatibility

## Source areas
- `backend/internal/application/mcp/`
- `backend/internal/infra/mcp/`
- `backend/internal/pkg/mcpauth/`
- `backend/internal/shared/security/outbound.go`
- `backend/internal/infra/outboundhttp/`

## Protocol and security invariants
- MCP Streamable HTTP / JSON-RPC; `tools/list` and `tools/call`.
- User context header: `X-Deeix-User-Context`.
- Signed context format: `v1.<base64url(payload)>.<base64url(hmac-sha256)>` using HMAC-SHA256.
- DEEIX owns registry, tool catalog, user authorization, execution, timeout/cancellation, concurrency and audit; the official SDK is protocol-only.
- SSRF controls validate hostname and resolved IP/CIDR, deny loopback/private/link-local/metadata targets unless explicitly allowed, and constrain redirects, response size and timeouts.

G5 requires signed-context replay/cross-user tests, permission tests, DNS rebinding/redirect/SSRF tests, bounded responses and cancellation. .NET MCP is not implemented yet.
