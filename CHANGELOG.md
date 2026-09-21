# Changelog

All notable changes to this project are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

### Added

- `Client` with per-area services (`Transactions`, `Customers`, `AddOns`,
  `Discounts`, `Plans`, `Subscriptions`, `Terminals`, `Settlements`,
  `Lookup`, `Users`, `APIKeys`, `Auth`) and functional options
  (`WithSandbox`, `WithEnvironment`, `WithBaseURL`, `WithHTTPClient`,
  `WithTimeout`, `WithUserAgent`, `WithBearerToken`).
- `NewClientFromEnv` reading `FLUIDPAY_API_KEY`, `FLUIDPAY_ENVIRONMENT` and
  `FLUIDPAY_BASE_URL`.
- Customer Vault support (`/vault/customer`), the canonical customer API:
  customers, stored addresses, and card/ACH/token/Apple Pay/Google Pay
  payment methods, with `validate`, `authorize` and `bypass_rule_engine`
  verification options.
- Subscription lifecycle: `Pause`, `MarkPastDue`, `Cancel`, `Complete`,
  `Activate` (with optional next bill date) and `Search`.
- `Transactions.Verify` and `Transactions.Credit`, plus every current
  sale/auth request field: base amount, idempotency key, processor id,
  tokens, Apple Pay and Google Pay, APMs, line items, custom fields,
  descriptor, partial payments, split transactions, Level 3, CIT/MIT
  indicators and HSA/FSA amounts.
- `Transaction.Approved`, `PartiallyApproved`, `Declined`,
  `GatewayDeclined` and `ProcessorError` helpers, and constants for
  transaction types, statuses, response codes and sources.
- Settlement batch search, terminal settlement and BIN lookup.
- Webhook helpers: `VerifyWebhookSignature`, `ComputeWebhookSignature`,
  `ParseWebhook`, `ParseWebhookRequest` and the `WebhookEvent` type,
  including test-delivery and secret-rotation handling.
- `NewIdempotencyKey` for safe retries.
- Typed `*Error` carrying the HTTP status, gateway message and
  `x-correlation-id`, with `IsUnauthorized` / `IsNotFound` helpers.
- `APIResource.LastResponse` on every result with the correlation id,
  status code and headers.
- Unit tests for every endpoint against an in-process gateway, runnable
  examples, an opt-in sandbox integration suite and a CI workflow.

### Changed

- **Breaking:** package functions taking a `Fluidpay` struct are replaced by
  methods on `*Client` services that take a `context.Context`. See the
  migration guide in the README.
- Requests honour the caller's context for cancellation and deadlines.
- Non-2xx responses and `"status": "failed"` envelopes now return an error
  instead of a zero-value struct.
- The default HTTP timeout is three minutes, per FluidPay's guidance for
  long-running authorizations.
- Amounts are `int` rather than `uint`.
- The deprecated `/customer` endpoints live on `Client.LegacyCustomers` and
  are marked `Deprecated`.

### Fixed

- Wrong JSON tags that silently dropped data: `email_reciept`,
  `payment_merhod_id`, `referenced_transaction_is`, `trasaction_source`,
  `tay_state`, and a `customer_id` tag containing a non-ASCII character.
- `DoRequest` discarded the result of `req.WithContext`, so the context was
  never applied.
- A stray `fmt.Println` of the response body in `UpdateCustomerPayment`.
- Terminal and ACH response codes are decoded as integers, matching the
  gateway.

### Removed

- The hard-coded sandbox API key (`TestAPIKey`) and the `test-key` file.
  Supply keys through the environment instead.
- The `LocalDev` flag; use `WithBaseURL` to target a local gateway.
