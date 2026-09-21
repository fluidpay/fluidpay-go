# fluidpay-go

[![Go Reference](https://pkg.go.dev/badge/github.com/fluidpay/fluidpay-go.svg)](https://pkg.go.dev/github.com/fluidpay/fluidpay-go)
[![CI](https://github.com/fluidpay/fluidpay-go/actions/workflows/ci.yml/badge.svg)](https://github.com/fluidpay/fluidpay-go/actions/workflows/ci.yml)

The official Go SDK for the [FluidPay](https://www.fluidpay.com) payment gateway.
It follows the current API reference at
[sandbox.fluidpay.com/docs](https://sandbox.fluidpay.com/docs/) and covers
transactions, the Customer Vault, recurring billing, terminals, settlements,
BIN lookup, user and API key administration, and webhook verification.

- Idiomatic Go: one `Client`, one service per API area, `context.Context`
  on every call, typed requests and responses, no external dependencies.
- API-key first: private keys go in the `Authorization` header, public
  keys are rejected up front, and `NewClientFromEnv` keeps secrets out of
  source.
- Honest errors: a `*fluidpay.Error` carries the HTTP status, the gateway
  message and the `x-correlation-id` support will ask you for.
- Fully unit tested against an in-process gateway, with runnable examples
  and an opt-in sandbox integration suite.

## Contents

- [Install](#install)
- [Quick start](#quick-start)
- [Authentication and API keys](#authentication-and-api-keys)
- [Configuration](#configuration)
- [Transactions](#transactions)
- [Customer Vault](#customer-vault)
- [Recurring billing](#recurring-billing)
- [Terminals, settlements and BIN lookup](#terminals-settlements-and-bin-lookup)
- [Users, API keys and JWTs](#users-api-keys-and-jwts)
- [Errors and correlation IDs](#errors-and-correlation-ids)
- [Webhooks](#webhooks)
- [Testing](#testing)
- [Migrating from the previous SDK](#migrating-from-the-previous-sdk)
- [Contributing](#contributing)

## Install

```sh
go get github.com/fluidpay/fluidpay-go
```

Requires Go 1.21 or newer.

## Quick start

```go
package main

import (
	"context"
	"log"
	"time"

	"github.com/fluidpay/fluidpay-go"
)

func main() {
	// Use your private key (api_...). Public keys (pub_...) are for the browser.
	client, err := fluidpay.NewClient("api_your_secret_key", fluidpay.WithSandbox())
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	tx, err := client.Transactions.Sale(ctx, &fluidpay.TransactionRequest{
		Amount:         1299, // cents: $12.99
		IdempotencyKey: fluidpay.NewIdempotencyKey(),
		PaymentMethod: fluidpay.PaymentMethod{
			Card: &fluidpay.CardPayment{
				Number:         "4111111111111111",
				ExpirationDate: "12/30",
				CVC:            "123",
			},
		},
	})
	if err != nil {
		log.Fatal(err) // transport error, bad credentials, invalid request...
	}
	if !tx.Approved() {
		log.Fatalf("declined: %d %s", tx.ResponseCode, tx.Response)
	}
	log.Printf("approved %s, auth code %s", tx.ID, tx.ResponseBody.Card.AuthCode)
}
```

Three things to remember:

1. Amounts are integers in the smallest currency unit (cents for USD).
2. A declined card is **not** an error. The call succeeds and
   `tx.Approved()` is false. Errors mean the request itself failed.
3. Give transaction calls a generous deadline. FluidPay recommends at least
   three minutes because some authorizations take that long; the default
   client timeout is set accordingly.

## Authentication and API keys

FluidPay issues two kinds of keys from **Settings → API Keys** in the
control panel:

| Key | Prefix | Use |
| --- | --- | --- |
| Private (secret) | `api_` | Server-side calls, which is everything this SDK does |
| Public | `pub_` | Browser-side tooling such as the Tokenizer and Cart sessions |

`NewClient` sends the private key as the `Authorization` header on every
request. It refuses `pub_` keys, because the gateway would answer
`401 Unauthorized` anyway.

Keep keys out of source code with `NewClientFromEnv`:

```sh
export FLUIDPAY_API_KEY=api_your_secret_key
export FLUIDPAY_ENVIRONMENT=sandbox     # or production (default)
# export FLUIDPAY_BASE_URL=http://localhost:8001/api   # optional override
```

```go
client, err := fluidpay.NewClientFromEnv()
```

Options passed to `NewClientFromEnv` are applied after the environment and
take precedence.

If a key stops working, check the usual suspects from the FluidPay docs:
the key was deleted, it is a public key, it has IP or URL restrictions
(`localhost` is never a valid URL), it belongs to a partner rather than a
merchant account, or it belongs to the other environment.
`fluidpay.IsUnauthorized(err)` identifies this case.

You can also manage keys programmatically; see
[Users, API keys and JWTs](#users-api-keys-and-jwts).

## Configuration

| Option | Effect |
| --- | --- |
| `WithSandbox()` | Target `https://sandbox.fluidpay.com/api`. Simulated processing, no billing. |
| `WithEnvironment(env)` | `fluidpay.Sandbox` or `fluidpay.Production` (default). |
| `WithBaseURL(url)` | Any base URL including the `/api` prefix, for local gateways or proxies. |
| `WithHTTPClient(hc)` | Supply your own `*http.Client` (custom transport, tracing, retries). |
| `WithTimeout(d)` | Set the HTTP timeout. Default is `fluidpay.DefaultTimeout` (3 minutes). |
| `WithUserAgent(ua)` | Prefix the `User-Agent` with your application name. |
| `WithBearerToken(jwt)` | Authenticate with a JWT from `Auth.ObtainJWT` instead of an API key. |

A `*Client` is safe for concurrent use. Build one and share it.

Every result embeds `fluidpay.APIResource`, so the HTTP response that
produced it is one field away:

```go
tx.LastResponse.CorrelationID // quote this in support requests
tx.LastResponse.StatusCode
tx.LastResponse.Header
```

## Transactions

`client.Transactions` maps to `/transaction`:

| Method | Endpoint | Notes |
| --- | --- | --- |
| `Sale` | `POST /transaction` | Authorize and capture in one call |
| `Authorize` | `POST /transaction` | Hold funds; capture later |
| `Verify` | `POST /transaction` | Zero-dollar verification |
| `Credit` | `POST /transaction` | Send funds without a prior sale |
| `Create` | `POST /transaction` | Any type via `TransactionRequest.Type` |
| `Capture` | `POST /transaction/{id}/capture` | Nil request captures the full amount |
| `Void` | `POST /transaction/{id}/void` | Cancels a pending settlement (auth reversal where supported) |
| `Refund` | `POST /transaction/{id}/refund` | Partial refunds allowed up to the settled total |
| `Get` | `GET /transaction/{id}` | |
| `Search` | `POST /transaction/search` | Defaults to the prior four months |

### Payment methods

`PaymentMethod` holds exactly one of:

```go
fluidpay.PaymentMethod{Card: &fluidpay.CardPayment{Number: "...", ExpirationDate: "MM/YY", CVC: "..."}}
fluidpay.PaymentMethod{Token: "temporary-token-from-tokenizer"}
fluidpay.PaymentMethod{Customer: &fluidpay.CustomerPayment{ID: customerID}}       // Customer Vault
fluidpay.PaymentMethod{ACH: &fluidpay.ACHPayment{RoutingNumber: "...", AccountNumber: "...", SecCode: fluidpay.SecCodeWeb, AccountType: fluidpay.AccountTypeChecking}}
fluidpay.PaymentMethod{Terminal: &fluidpay.TerminalPayment{ID: terminalID, PrintReceipt: fluidpay.ReceiptBoth}}
fluidpay.PaymentMethod{ApplePayToken: &fluidpay.ApplePayToken{TemporaryToken: "..."}}  // production only
fluidpay.PaymentMethod{GooglePayToken: json.RawMessage(tokenFromGoogle)}               // production only
fluidpay.PaymentMethod{APM: &fluidpay.APMPayment{Type: "klarna", MerchantRedirectURL: "..."}}
```

`TransactionRequest` exposes every documented field: `BaseAmount` (let the
gateway add surcharges), tax/shipping/tip/discount amounts,
`PaymentAdjustment`, billing and shipping addresses, `LineItems` and the
other Level 3 fields, `CustomFields`, `Descriptor`, `CreateVaultRecord`,
`AllowPartialPayment`, `SplitTransactionAmount`, CIT/MIT indicators and
HSA/FSA amounts. Zero values are omitted from the request.

### Authorize, then capture

```go
auth, err := client.Transactions.Authorize(ctx, &fluidpay.TransactionRequest{
	Amount:        5000,
	PaymentMethod: fluidpay.PaymentMethod{Token: token},
})
// ... ship the goods ...
captured, err := client.Transactions.Capture(ctx, auth.ID, &fluidpay.CaptureRequest{Amount: 4500})
```

### Reading the result

```go
switch {
case tx.Approved():          // response codes 100-199
	if tx.PartiallyApproved() {
		// tx.AmountAuthorized < tx.Amount: collect the rest or void
	}
case tx.Declined():          // 200-299, issuer decline
case tx.GatewayDeclined():   // 300-399, duplicate, rule engine, blocked customer
case tx.ProcessorError():    // 400-499
}
```

The processor detail lives under `tx.ResponseBody.Card`, `.ACH` or
`.Terminal`, whichever applies, including AVS and CVV response codes.

### Searching

Query fields take small operator structs. Constructors keep it short:

```go
page, err := client.Transactions.Search(ctx, &fluidpay.TransactionSearchRequest{
	Status:     fluidpay.Equal(fluidpay.StatusSettled),
	Amount:     fluidpay.GreaterThan(10000),
	CreatedAt:  fluidpay.DateRange(since, time.Now()),   // or fluidpay.Day(day)
	CustomerID: fluidpay.NotEqual(""),
	Limit:      100,
	Offset:     0,
})
for _, tx := range page.Data { ... }
fmt.Println(page.TotalCount)
```

### Idempotency

Set `IdempotencyKey` (a UUID; `fluidpay.NewIdempotencyKey()` makes one) and
retry the same request safely for five minutes, or for `IdempotencyTime`
seconds.

## Customer Vault

`client.Customers` is the canonical customer API (`/vault/customer`).
Store a customer once, then charge with
`PaymentMethod{Customer: &CustomerPayment{ID: customer.ID}}`.

```go
customer, err := client.Customers.Create(ctx, &fluidpay.CustomerCreateRequest{
	VerificationOptions: fluidpay.VerificationOptions{Validate: true}, // $0 verification
	Description:         "Jane Doe",
	DefaultPayment: &fluidpay.VaultPaymentMethod{
		Card: &fluidpay.VaultCard{Number: "4111111111111111", ExpirationDate: "12/30"},
	},
	DefaultBillingAddress: &fluidpay.VaultAddress{
		FirstName: "Jane", LastName: "Doe", Line1: "123 Some St",
		City: "Chicago", State: "IL", PostalCode: "60601", Country: "US",
	},
})

customer.ID
customer.DefaultCard().MaskedNumber
customer.Addresses, customer.Cards, customer.ACH   // flattened from the gateway's nested shape
```

| Method | Endpoint |
| --- | --- |
| `Create`, `Get`, `Search`, `Update`, `Delete` | `/vault/customer`, `/vault/{id}`, `/vault/customer/search` |
| `CreateAddress`, `UpdateAddress`, `DeleteAddress` | `/vault/customer/{id}/address[/{addressId}]` |
| `CreateCard`, `UpdateCard`, `DeleteCard` | `/vault/customer/{id}/card[/{paymentMethodId}]` |
| `CreateACH`, `UpdateACH`, `DeleteACH` | `/vault/customer/{id}/ach[/{paymentMethodId}]` |
| `CreateToken`, `UpdateToken` | `/vault/customer/{id}/token[/{paymentMethodId}]` |
| `CreateApplePay`, `CreateGooglePay` | `/vault/customer/{id}/applepay`, `/googlepay` |

Create calls return the whole customer; the new record's ID is in
`CreatedAddressID` or `CreatedPaymentMethodID`. `VerificationOptions`
(`Validate`, `Authorize`, `BypassRuleEngine`) are sent as query parameters
on customer and payment method creation.

The gateway rejects an address identical to one already stored
(`invalid: would create a duplicate address`). Check `customer.Addresses`
and reuse the existing ID instead.

The older `/customer` endpoints are still available as
`client.LegacyCustomers`, marked `Deprecated`. FluidPay advises against
mixing the two families in one application.

## Recurring billing

Add-ons and discounts adjust a charge by a fixed amount or a percentage.
Plans define the schedule. Subscriptions apply a plan to a vault customer.

```go
plan, err := client.Plans.Create(ctx, &fluidpay.PlanRequest{
	Name:                 "Pro monthly",
	Amount:               4900,
	BillingCycleInterval: 1,
	BillingFrequency:     fluidpay.BillingMonthly,   // monthly, twice_monthly, daily
	BillingDays:          "1",                        // "1,15" for twice monthly, "0" for last day
})

sub, err := client.Subscriptions.Create(ctx, &fluidpay.SubscriptionRequest{
	PlanID:               plan.ID,
	Customer:             fluidpay.SubscriptionCustomer{ID: customer.ID},
	Amount:               plan.Amount,
	BillingCycleInterval: 1,
	BillingFrequency:     fluidpay.BillingMonthly,
	BillingDays:          "1",
	NextBillDate:         "2026-10-01",
})
```

Lifecycle methods on `client.Subscriptions`: `Pause`, `Activate` (with an
optional next bill date), `MarkPastDue`, `Cancel`, `Complete`, `Delete`,
plus `Search` and `Get`.

The billing engine moves a subscription through `active` → `failing` →
`failed` on consecutive declines and sets `error` on configuration
problems. `sub.NeedsAttention()` covers all three. Renewal charges are
ordinary transactions with `TransactionSource == fluidpay.SourceRecurring`
and `SubscriptionID` set.

## Terminals, settlements and BIN lookup

```go
terminals, err := client.Terminals.List(ctx)             // GET /terminals
err = client.Terminals.Settle(ctx, terminalID)           // POST /terminal/{id}/settle

batches, err := client.Settlements.SearchBatches(ctx, &fluidpay.SettlementBatchSearchRequest{
	BatchDate: fluidpay.Day(time.Now().AddDate(0, 0, -1)),
})                                                       // POST /settlement/batch/search
batches.Summary, batches.Results

bin, err := client.Lookup.BIN(ctx, &fluidpay.BINLookupRequest{
	BIN: "424242", Country: "US", State: "IL",
})                                                       // POST /lookup/bin/protected
bin.CardBrand, bin.CardType, bin.IsSurchargeable
```

To charge through a terminal, run a normal sale with
`PaymentMethod{Terminal: ...}`. Sandbox terminal behaviour is driven by
the terminal's TPN (see [Testing](#testing)).

## Users, API keys and JWTs

```go
me, err := client.Users.Current(ctx)         // who owns this key
me.Can("process_refund")

key, err := client.APIKeys.Create(ctx, &fluidpay.APIKeyCreateRequest{
	Type: fluidpay.APIKeyTypePrivate, Name: "backend-2026",
})
key.Key                                       // shown once; store it now
keys, err := client.APIKeys.List(ctx)
err = client.APIKeys.Delete(ctx, key.ID)      // revoke
```

`client.Users` also offers `Get`, `List`, `Create`, `Update`, `Delete` and
`ChangePassword`. For user-driven sessions, `client.Auth.ObtainJWT`
exchanges a username and password for a token you can pass to
`WithBearerToken`; `Auth.Logout`, `ForgotUsername`, `ForgotPassword` and
`ResetPassword` round out account recovery.

## Errors and correlation IDs

Every method returns a plain `error` for transport failures (DNS, TLS,
timeouts, cancelled contexts) and a `*fluidpay.Error` when the gateway
rejects the request with a non-2xx status or a `"status": "failed"`
envelope.

```go
tx, err := client.Transactions.Sale(ctx, req)
if err != nil {
	if apiErr, ok := fluidpay.AsError(err); ok {
		log.Printf("%d %s (correlation id %s)", apiErr.StatusCode, apiErr.Msg, apiErr.CorrelationID)
	}
	if fluidpay.IsUnauthorized(err) { /* rotate or fix the key */ }
	if fluidpay.IsNotFound(err)     { /* unknown id */ }
	return err
}
```

`Error.Error()` reads like
`fluidpay: POST /transaction: 400 Bad Request: bad request error: invalid Postal Code (correlation id ...)`.

The `x-correlation-id` header identifies the request in FluidPay's logs.
It is on every `*Error` and on every successful result under
`LastResponse.CorrelationID`. Include it in support tickets.

## Webhooks

Deliveries are signed with HMAC-SHA256 over the raw body, base64url without
padding, in the `Signature` header. Verify against the raw bytes exactly as
received.

```go
http.HandleFunc("/webhooks/fluidpay", func(w http.ResponseWriter, r *http.Request) {
	event, err := fluidpay.ParseWebhookRequest(r, os.Getenv("FLUIDPAY_WEBHOOK_SECRET"))
	if err != nil {
		http.Error(w, "invalid signature", http.StatusUnauthorized)
		return
	}
	if event.IsTest() {          // signed with fluidpay.WebhookTestSecret, verified automatically
		w.WriteHeader(http.StatusOK)
		return
	}
	// event.EventID is the idempotency key; store it before doing side effects.
	tx, _ := event.Transaction()
	log.Printf("%s: %s is %s", event.Type, tx.ID, tx.Status)
	w.WriteHeader(http.StatusOK)
})
```

`ParseWebhookRequest` also tries the `Signature-Previous` header, and you
may pass the previous secret as an extra argument during rotation. For
non-HTTP contexts use `ParseWebhook(body, signature, secret)` or the lower
level `VerifyWebhookSignature`.

## Testing

### Sandbox test data

Everything below is documented under
[Test Data](https://sandbox.fluidpay.com/docs/test_data/) and only works in
the sandbox. Any future expiration date and any CVC work; cards not listed
are approved.

| Card | Result |
| --- | --- |
| `4111111111111111` | Approved (Visa debit, not surchargeable) |
| `4005519200000004` | Approved (Visa credit, surchargeable; use with `BaseAmount`) |
| `4000000000000002` | Declined |
| `4000000000009995` | Insufficient funds |
| `4000000000000069` | Expired card |
| `4000000000000051` | Partial approval (half the amount) |
| `4000000000000010` | Auth succeeds, refund declines |

| Trigger | Effect |
| --- | --- |
| CVC `200` | CVV response `N` (no match) |
| Postal code `20000` | AVS response `N` (no match) |
| ACH routing `000000000` | Decline |
| ACH routing `000000001` / `000000002` | Return / late return during settlement |
| Terminal TPN `000000000001` / `2` / `3` / `4` | Success / decline / error / success with tip |

### Unit tests

The unit suite runs entirely against an in-process fake gateway and checks
every endpoint's method, path, query string, request body and response
decoding, using response fixtures taken from the API reference:

```sh
go test -race ./...
```

### Sandbox integration tests

An opt-in suite exercises the real sandbox with the test data above and
cleans up after itself:

```sh
FLUIDPAY_API_KEY=api_... FLUIDPAY_ENVIRONMENT=sandbox go test -tags integration -run TestIntegration ./...
```

The CI workflow runs it on pushes when a `FLUIDPAY_SANDBOX_API_KEY`
repository secret is configured.

### Examples

Runnable programs live under [`examples/`](examples):

| Program | Shows |
| --- | --- |
| `examples/sale` | A card sale and void with error handling |
| `examples/vault` | Storing a customer, adding an address, charging the stored card |
| `examples/subscription` | Plan and subscription creation, pause/activate/cancel, search |
| `examples/webhook` | An HTTP handler that verifies signatures and dedupes events |

```sh
FLUIDPAY_API_KEY=api_... FLUIDPAY_ENVIRONMENT=sandbox go run ./examples/sale
```

Godoc examples are in [`example_test.go`](example_test.go).

## Migrating from the previous SDK

The earlier package exposed free functions that took a `Fluidpay` struct
and returned envelope structs. Those are gone; the table shows the
replacements.

| Before | After |
| --- | --- |
| `Fluidpay{APIKey, Client, Ctx, Sandbox, LocalDev}` | `fluidpay.NewClient(key, fluidpay.WithSandbox())`; `WithBaseURL` replaces `LocalDev`; the context is passed per call |
| `DoTransaction(fp, req)` | `client.Transactions.Sale/Authorize/Verify/Credit/Create(ctx, req)` |
| `GetTransactionStatus(fp, id)` | `client.Transactions.Get(ctx, id)` |
| `QueryTransactions(fp, req)` | `client.Transactions.Search(ctx, req)` with `fluidpay.Equal`, `IntEqual`, `DateRange`... |
| `CaptureTransaction(fp, req, id)` | `client.Transactions.Capture(ctx, id, req)` |
| `VoidTransaction(fp, id)` | `client.Transactions.Void(ctx, id)` |
| `RefundTransaction(fp, req, id)` | `client.Transactions.Refund(ctx, id, req)` |
| `CreateCustomer` and friends (`/customer`) | `client.Customers.*` (Customer Vault) for new code; `client.LegacyCustomers.*` keeps the old endpoints |
| `CreateAddOn`, `GetAddOns`, ... | `client.AddOns.Create/Get/List/Update/Delete` |
| `CreateDiscount`, ... | `client.Discounts.*` |
| `CreatePlan`, ... | `client.Plans.*` |
| `CreateSubscription`, ... | `client.Subscriptions.*` plus `Pause`, `Activate`, `Cancel`, `Complete`, `MarkPastDue`, `Search` |
| `GetTerminals(fp)` | `client.Terminals.List(ctx)` |
| `CreateKey`, `GetKeys`, `DeleteKey` | `client.APIKeys.Create/List/Delete` |
| `CreateUser`, `GetCurrentUser`, ... | `client.Users.*` |
| `ObtainJWT`, `TokenLogout`, `Forgotten*`, `PasswordReset` | `client.Auth.*`; `NewAuth`/`SetAuth` become `WithBearerToken` |
| `TransactionRequest.EmailReciept` | `EmailReceipt` (the JSON tag was wrong before) |
| `PaymentMethodsRequest`, `CardRequest`, `AchRequest`, ... | `PaymentMethod`, `CardPayment`, `ACHPayment`, `CustomerPayment`, `TerminalPayment` |
| `TerminalResponse` / `TransactionResponse` | `*Transaction`; check `tx.Approved()` instead of comparing `Msg == "success"` |
| `uint` amounts | `int` amounts |
| `TestAPIKey` constant | Removed. Use `FLUIDPAY_API_KEY`. |

Behavioural changes worth knowing:

- Failed requests now return an error instead of an empty struct. Declines
  still return a `*Transaction`.
- The context you pass is honoured; the old code ignored it.
- Response `Data` envelopes are unwrapped for you. Lists return a `*XList`
  with `Data` and `TotalCount`.

## Contributing

```sh
gofmt -l .            # formatting
go vet ./...          # static checks
go test -race ./...   # unit tests
```

Please add a unit test with a fixture from the API reference for any new
endpoint, and keep the [CHANGELOG](CHANGELOG.md) current.
