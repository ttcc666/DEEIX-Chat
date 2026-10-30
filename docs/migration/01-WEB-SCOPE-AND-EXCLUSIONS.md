# 01 — Web scope and exclusions

## In scope
- Web-used Go HTTP API, application use cases, domain rules, repositories and infrastructure.
- Existing Next.js / React / TypeScript Web frontend; no page or interaction rewrite.
- `packages/api-contract/` generated TypeScript API contract.
- Web-consumed parts of `packages/core/` without removing Desktop/Mobile exports.
- Web deployment, configuration, static hosting, reverse proxy, health checks and CI.
- Auth/OIDC, model routing, Providers, MCP, files, OCR, RAG, Memory, billing, payments, admin and audit.

## Explicitly excluded
- `apps/desktop/`, Tauri 2, Rust, Cargo and Desktop UI.
- Go Desktop sidecar packaging/startup.
- Desktop tray, deep links, OS Keychain, updater, signing/notarization/installers.
- `.github/workflows/desktop-*.yml` release workflows.
- Desktop-only SQLite lifecycle and local process/port management.

**Do not delete or refactor `backend/`**: the existing Go backend is still required by the Desktop sidecar. Do not modify Tauri/Cargo/Desktop workflows. The target is only to remove Go from Web production after Phase 9 approval.
