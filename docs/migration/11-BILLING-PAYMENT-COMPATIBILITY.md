# 11 — Billing and payment compatibility

Source areas: `backend/internal/application/billing/`, `backend/internal/infra/payment/stripe/`, `backend/internal/infra/payment/epay/`, `backend/internal/infra/persistence/models/billing.go`, PostgreSQL persistence/transaction code.

Preserve nanousd monetary precision, payment-order snapshots/state machine, Stripe webhook signature/timestamp validation, EPay submit.php MD5 canonicalization and constant-time signature comparison, webhook idempotency, subscription activation, atomic balance+transaction+ledger writes, usage-ledger idempotency by user/reference, reservation/settlement/release/reconciliation and PostgreSQL row locking.

G7 tests: golden Stripe/EPay signatures, duplicate/reordered webhooks, concurrent top-ups, double settlement, reservation expiry, ledger idempotency, rollback and reconciliation. Only isolated fixtures/payment sandbox; no production payment service.
