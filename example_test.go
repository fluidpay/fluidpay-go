package fluidpay_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fluidpay/fluidpay-go"
)

func ExampleNewClient() {
	client, err := fluidpay.NewClient(os.Getenv("FLUIDPAY_API_KEY"), fluidpay.WithSandbox())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(client.BaseURL())
}

func ExampleNewClientFromEnv() {
	// FLUIDPAY_API_KEY=api_...  FLUIDPAY_ENVIRONMENT=sandbox
	client, err := fluidpay.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	_ = client
}

func ExampleTransactionsService_Sale() {
	client, err := fluidpay.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	tx, err := client.Transactions.Sale(ctx, &fluidpay.TransactionRequest{
		Amount:         1299, // $12.99
		IdempotencyKey: fluidpay.NewIdempotencyKey(),
		OrderID:        "order-1001",
		PaymentMethod: fluidpay.PaymentMethod{
			Card: &fluidpay.CardPayment{
				Number:         "4111111111111111",
				ExpirationDate: "12/30",
				CVC:            "123",
			},
		},
		BillingAddress: &fluidpay.Address{PostalCode: "60601", Country: "US"},
	})
	if err != nil {
		// The request itself failed: network, auth, validation.
		if apiErr, ok := fluidpay.AsError(err); ok {
			log.Fatalf("gateway rejected the request: %s (correlation id %s)", apiErr.Msg, apiErr.CorrelationID)
		}
		log.Fatal(err)
	}
	if !tx.Approved() {
		// A decline is a successful call with a decline response code.
		log.Printf("declined: %d %s", tx.ResponseCode, tx.ResponseBody.Card.ProcessorResponseText)
		return
	}
	fmt.Println("approved", tx.ID, tx.ResponseBody.Card.AuthCode)
}

func ExampleTransactionsService_Authorize() {
	client, _ := fluidpay.NewClientFromEnv()
	ctx := context.Background()

	auth, err := client.Transactions.Authorize(ctx, &fluidpay.TransactionRequest{
		Amount:        5000,
		PaymentMethod: fluidpay.PaymentMethod{Token: "temporary-token-from-tokenizer"},
	})
	if err != nil || !auth.Approved() {
		log.Fatal("authorization failed")
	}

	// Ship the goods, then capture. A nil request captures the full amount.
	captured, err := client.Transactions.Capture(ctx, auth.ID, &fluidpay.CaptureRequest{Amount: 4500})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(captured.Status) // pending_settlement
}

func ExampleTransactionsService_Search() {
	client, _ := fluidpay.NewClientFromEnv()
	since := time.Now().AddDate(0, 0, -7)

	page, err := client.Transactions.Search(context.Background(), &fluidpay.TransactionSearchRequest{
		Status:    fluidpay.Equal(fluidpay.StatusSettled),
		Amount:    fluidpay.GreaterThan(10000),
		CreatedAt: fluidpay.DateRange(since, time.Now()),
		Limit:     100,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d of %d\n", len(page.Data), page.TotalCount)
}

func ExampleCustomersService_Create() {
	client, _ := fluidpay.NewClientFromEnv()
	ctx := context.Background()

	customer, err := client.Customers.Create(ctx, &fluidpay.CustomerCreateRequest{
		VerificationOptions: fluidpay.VerificationOptions{Validate: true},
		Description:         "Jane Doe",
		DefaultPayment: &fluidpay.VaultPaymentMethod{
			Card: &fluidpay.VaultCard{Number: "4111111111111111", ExpirationDate: "12/30"},
		},
		DefaultBillingAddress: &fluidpay.VaultAddress{
			FirstName: "Jane", LastName: "Doe", Line1: "123 Some St",
			City: "Chicago", State: "IL", PostalCode: "60601", Country: "US",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Charge the stored default payment method later.
	tx, err := client.Transactions.Sale(ctx, &fluidpay.TransactionRequest{
		Amount:        2500,
		PaymentMethod: fluidpay.PaymentMethod{Customer: &fluidpay.CustomerPayment{ID: customer.ID}},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tx.Approved())
}

func ExampleSubscriptionsService_Create() {
	client, _ := fluidpay.NewClientFromEnv()
	ctx := context.Background()

	plan, err := client.Plans.Create(ctx, &fluidpay.PlanRequest{
		Name:                 "Pro monthly",
		Amount:               4900,
		BillingCycleInterval: 1,
		BillingFrequency:     fluidpay.BillingMonthly,
		BillingDays:          "1",
	})
	if err != nil {
		log.Fatal(err)
	}

	sub, err := client.Subscriptions.Create(ctx, &fluidpay.SubscriptionRequest{
		PlanID:               plan.ID,
		Customer:             fluidpay.SubscriptionCustomer{ID: "customer-id"},
		Amount:               plan.Amount,
		BillingCycleInterval: 1,
		BillingFrequency:     fluidpay.BillingMonthly,
		BillingDays:          "1",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(sub.Status, sub.NextBillDate)
}

func ExampleParseWebhookRequest() {
	secret := os.Getenv("FLUIDPAY_WEBHOOK_SECRET")

	http.HandleFunc("/webhooks/fluidpay", func(w http.ResponseWriter, r *http.Request) {
		event, err := fluidpay.ParseWebhookRequest(r, secret)
		if err != nil {
			http.Error(w, "invalid signature", http.StatusUnauthorized)
			return
		}
		if event.IsTest() {
			w.WriteHeader(http.StatusOK)
			return
		}
		tx, err := event.Transaction()
		if err == nil {
			log.Printf("event %s: transaction %s is %s", event.EventID, tx.ID, tx.Status)
		}
		w.WriteHeader(http.StatusOK)
	})
}

func ExampleVerifyWebhookSignature() {
	secret := "12345678-1234-1234-1234-123456789012"
	body := []byte(`{"data":"this is test data"}`)
	signature := "JacUiw_ztpEZJWvOhhKoHTLBf4b-aZv9n_0YmJJxltc"

	fmt.Println(fluidpay.VerifyWebhookSignature(secret, body, signature))
	fmt.Println(fluidpay.VerifyWebhookSignature("wrong-secret", body, signature))
	// Output:
	// true
	// false
}

func ExampleError() {
	client, _ := fluidpay.NewClient("api_revoked_key")
	_, err := client.Terminals.List(context.Background())
	if fluidpay.IsUnauthorized(err) {
		fmt.Println("check the key, its environment and any IP restrictions")
	}
}

func ExampleWithResponseHook() {
	// Record the correlation id of every gateway response in your logs so
	// support requests can reference it.
	client, err := fluidpay.NewClientFromEnv(fluidpay.WithResponseHook(func(r *fluidpay.APIResponse) {
		log.Printf("fluidpay %s /%s -> %d correlation id %s", r.Method, r.Path, r.StatusCode, r.CorrelationID)
	}))
	if err != nil {
		log.Fatal(err)
	}
	_, _ = client.Terminals.List(context.Background())
}

func ExampleCorrelationID() {
	client, _ := fluidpay.NewClientFromEnv()
	ctx := context.Background()

	tx, err := client.Transactions.Sale(ctx, &fluidpay.TransactionRequest{
		Amount:        1000,
		PaymentMethod: fluidpay.PaymentMethod{Token: "temporary-token"},
	})
	if err != nil {
		// Works for gateway rejections and decoding failures, wrapped or not.
		log.Fatalf("sale failed: %v (correlation id %q)", err, fluidpay.CorrelationID(err))
	}
	fmt.Println(tx.ID, tx.CorrelationID())

	// Calls with no other result return the response metadata directly.
	resp, err := client.Transactions.Void(ctx, tx.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.CorrelationID)
}
