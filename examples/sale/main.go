// Command sale charges a test card in the sandbox, then voids the charge.
//
//	FLUIDPAY_API_KEY=api_... FLUIDPAY_ENVIRONMENT=sandbox go run ./examples/sale
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/fluidpay/fluidpay-go"
)

func main() {
	// Log the correlation id of every request; FluidPay support uses it to
	// find the request in their logs.
	client, err := fluidpay.NewClientFromEnv(fluidpay.WithResponseHook(func(r *fluidpay.APIResponse) {
		log.Printf("fluidpay %s /%s -> %d (correlation id %s)", r.Method, r.Path, r.StatusCode, r.CorrelationID)
	}))
	if err != nil {
		log.Fatal(err)
	}

	// Authorizations can take a while; give the call the recommended budget.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	tx, err := client.Transactions.Sale(ctx, &fluidpay.TransactionRequest{
		Amount:         1299, // $12.99, always in cents
		IdempotencyKey: fluidpay.NewIdempotencyKey(),
		OrderID:        "order-1001",
		Description:    "fluidpay-go example",
		PaymentMethod: fluidpay.PaymentMethod{
			Card: &fluidpay.CardPayment{
				Number:         "4111111111111111", // sandbox: always approved
				ExpirationDate: "12/30",
				CVC:            "123",
			},
		},
		BillingAddress: &fluidpay.Address{
			FirstName: "Jane", LastName: "Doe",
			AddressLine1: "123 Some St", City: "Chicago", State: "IL",
			PostalCode: "60601", Country: "US",
		},
	})
	if err != nil {
		if apiErr, ok := fluidpay.AsError(err); ok {
			log.Fatalf("gateway rejected the request: %s (correlation id %s)", apiErr.Msg, apiErr.CorrelationID)
		}
		log.Fatalf("%v (correlation id %q)", err, fluidpay.CorrelationID(err))
	}

	if !tx.Approved() {
		log.Fatalf("declined: code %d (%s)", tx.ResponseCode, tx.ResponseBody.Card.ProcessorResponseText)
	}
	fmt.Printf("approved %s for %d cents, auth code %s, AVS %s (correlation id %s)\n",
		tx.ID, tx.AmountAuthorized, tx.ResponseBody.Card.AuthCode, tx.ResponseBody.Card.AVSResponseCode, tx.CorrelationID())

	resp, err := client.Transactions.Void(ctx, tx.ID)
	if err != nil {
		log.Fatalf("void: %v", err)
	}
	fmt.Printf("voided %s (correlation id %s)\n", tx.ID, resp.CorrelationID)
}
