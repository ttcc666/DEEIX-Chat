# 07 — Authentication and crypto compatibility

## Source areas
- `backend/internal/application/auth/`
- `backend/internal/transport/http/auth/`
- `backend/internal/transport/http/auth/token_delivery.go`
- `backend/internal/pkg/token/token.go`
- `backend/internal/pkg/secretbox/secretbox.go`
- `backend/internal/infra/config/config.go`

## Compatibility requirements extracted from source audit
- Password hashes: bcrypt cost 12; old hashes must verify without forced reset.
- JWT: HS256; preserve actual issuer/audience/claims/JTI/expiry from Go implementation.
- Default access-token TTL: 24 hours.
- Default refresh-token TTL: 720 hours; hash at rest; rotation, previous-token reuse detection and revocation.
- Refresh cookie: `deeix_chat_refresh_token`; Path `/api/v1/auth`; HttpOnly; SameSite=Lax; Secure follows HTTPS/environment.
- TOTP: 6 digits, 30-second step, ±1 window; recovery codes stored hashed.
- OAuth2/OIDC and PKCE; preserve callback and one-time authorization exchange.
- AES-256-GCM; SHA-256-derived key; envelope `v1:` + Base64(nonce || ciphertext/tag); some values bind AAD.
- Current key encrypts; current and previous keys decrypt; rotation must be atomic.

G2 tests: Go-generated bcrypt/AES/JWT golden vectors; old ciphertext decrypts in .NET; wrong AAD/tampering fails; refresh rotation/reuse/concurrency/revoke; Cookie attributes; TOTP/recovery; OAuth state/nonce/PKCE/callback. No cross-language .NET tests have run yet.
