// Command vault stores a customer in the Customer Vault, charges the stored
// card, then deletes the customer.
//
//	FLUIDPAY_API_KEY=api_... FLUIDPAY_ENVIRONMENT=sandbox go run ./examples/vault
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/fluidpay/fluidpay-go"
)

func main() {
	client, err := fluidpay.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	// 1. Store the customer with a default card and billing address.
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
	if err != nil {
		log.Fatalf("create customer: %v", err)
	}
	defer func() {
		if _, err := client.Customers.Delete(context.Background(), customer.ID); err != nil {
			log.Printf("delete customer: %v", err)
		}
	}()
	fmt.Printf("customer %s stored card %s\n", customer.ID, customer.DefaultCard().MaskedNumber)

	// 2. Add a second address. Reuse an existing one when the details match,
	// because the gateway rejects duplicate addresses.
	updated, err := client.Customers.CreateAddress(ctx, customer.ID, &fluidpay.VaultAddress{
		FirstName: "Jane", LastName: "Doe", Line1: "456 Other Ave",
		City: "Chicago", State: "IL", PostalCode: "60602", Country: "US",
	})
	if err != nil {
		log.Fatalf("create address: %v", err)
	}
	fmt.Printf("shipping address %s added (%d addresses on file)\n", updated.CreatedAddressID, len(updated.Addresses))

	// 3. Charge the stored default payment method.
	tx, err := client.Transactions.Sale(ctx, &fluidpay.TransactionRequest{
		Amount:         2500,
		IdempotencyKey: fluidpay.NewIdempotencyKey(),
		PaymentMethod: fluidpay.PaymentMethod{
			Customer: &fluidpay.CustomerPayment{
				ID:                customer.ID,
				ShippingAddressID: updated.CreatedAddressID,
			},
		},
	})
	if err != nil {
		log.Fatalf("sale: %v", err)
	}
	fmt.Printf("charged %d cents: %s (%s)\n", tx.Amount, tx.Response, tx.ID)
	_, _ = client.Transactions.Void(ctx, tx.ID)
}
