// Package fluidpay is the official Go SDK for the FluidPay payment gateway
// REST API (https://sandbox.fluidpay.com/docs/).
//
// # Getting started
//
// Create a Client with your private API key (it starts with "api_") and the
// environment you want to talk to. Public keys ("pub_") are client-side only
// and are rejected by NewClient.
//
//	client, err := fluidpay.NewClient("api_your_secret_key", fluidpay.WithSandbox())
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	tx, err := client.Transactions.Sale(ctx, &fluidpay.TransactionRequest{
//		Amount: 1299, // cents: $12.99
//		PaymentMethod: fluidpay.PaymentMethod{
//			Card: &fluidpay.CardPayment{
//				Number:         "4111111111111111",
//				ExpirationDate: "12/30",
//				CVC:            "123",
//			},
//		},
//	})
//	if err != nil {
//		// Transport failures, HTTP errors and gateway "failed" envelopes
//		// all surface here, usually as a *fluidpay.Error.
//		log.Fatal(err)
//	}
//	if !tx.Approved() {
//		// A decline is a successful API call: inspect the response code.
//		log.Printf("declined: %d %s", tx.ResponseCode, tx.Response)
//	}
//
// Alternatively, NewClientFromEnv reads FLUIDPAY_API_KEY, FLUIDPAY_ENVIRONMENT
// and FLUIDPAY_BASE_URL so credentials never have to live in source code.
//
// # Services
//
// The Client exposes one service per API area:
//
//   - Transactions: sale, authorize, capture, void, refund, verification, credit, get, search
//   - Customers: the Customer Vault (customers, stored addresses and payment methods)
//   - LegacyCustomers: the deprecated /customer endpoints, kept for existing integrations
//   - AddOns, Discounts, Plans, Subscriptions: recurring billing
//   - Terminals: physical terminal listing and settlement
//   - Settlements: settlement batch search
//   - Lookup: BIN lookup
//   - Users, APIKeys, Auth: user administration and JWT authentication
//
// Every method takes a context.Context so callers control cancellation and
// deadlines. Amounts are always integers in the smallest currency unit
// (cents for USD).
//
// # Errors
//
// Requests that fail at the transport level return the underlying error.
// Requests the gateway rejects (a non-2xx status or a "failed" envelope)
// return a *Error carrying the HTTP status, the gateway message and the
// x-correlation-id you should quote when contacting FluidPay support.
//
// # Correlation IDs
//
// Every gateway response carries an x-correlation-id header. Results expose
// it through the embedded APIResource (tx.CorrelationID()), methods with no
// other result (Void, Delete, ...) return the *APIResponse directly, errors
// expose it through CorrelationID(err), and WithResponseHook lets you log it
// for every request.
//
// # Webhooks
//
// VerifyWebhookSignature and ParseWebhook help you authenticate and decode
// webhook deliveries sent by the gateway.
package fluidpay
