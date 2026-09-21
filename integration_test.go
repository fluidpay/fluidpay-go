//go:build integration

package fluidpay_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/fluidpay/fluidpay-go"
)

// These tests run against the FluidPay sandbox and are opt-in:
//
//	FLUIDPAY_API_KEY=api_... go test -tags integration -run Integration ./...
//
// They only use documented sandbox test data and clean up what they create.

func integrationClient(t *testing.T) *fluidpay.Client {
	t.Helper()
	if os.Getenv(fluidpay.EnvAPIKey) == "" {
		t.Skipf("%s not set; skipping sandbox integration test", fluidpay.EnvAPIKey)
	}
	if os.Getenv(fluidpay.EnvEnvironment) == "" {
		os.Setenv(fluidpay.EnvEnvironment, string(fluidpay.Sandbox))
	}
	client, err := fluidpay.NewClientFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func TestIntegration_CurrentUser(t *testing.T) {
	client := integrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	me, err := client.Users.Current(ctx)
	if err != nil {
		t.Fatalf("Users.Current: %v", err)
	}
	if me.ID == "" || me.Username == "" {
		t.Fatalf("unexpected user %+v", me)
	}
	t.Logf("authenticated as %s (%s)", me.Username, me.AccountType)
}

func TestIntegration_SaleVoid(t *testing.T) {
	client := integrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	tx, err := client.Transactions.Sale(ctx, &fluidpay.TransactionRequest{
		Amount:         1299,
		IdempotencyKey: fluidpay.NewIdempotencyKey(),
		Description:    "fluidpay-go integration test",
		PaymentMethod: fluidpay.PaymentMethod{Card: &fluidpay.CardPayment{
			Number: "4111111111111111", ExpirationDate: "12/30", CVC: "123",
		}},
		BillingAddress: &fluidpay.Address{PostalCode: "60601", Country: "US"},
	})
	if err != nil {
		t.Fatalf("Sale: %v", err)
	}
	if !tx.Approved() {
		t.Fatalf("sandbox sale was not approved: %d %s", tx.ResponseCode, tx.Response)
	}

	got, err := client.Transactions.Get(ctx, tx.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.ID != tx.ID {
		t.Fatalf("Get returned %s, want %s", got.ID, tx.ID)
	}

	if _, err := client.Transactions.Void(ctx, tx.ID); err != nil {
		t.Fatalf("Void: %v", err)
	}

	declined, err := client.Transactions.Sale(ctx, &fluidpay.TransactionRequest{
		Amount: 1000,
		PaymentMethod: fluidpay.PaymentMethod{Card: &fluidpay.CardPayment{
			Number: "4000000000000002", ExpirationDate: "12/30", CVC: "123",
		}},
	})
	if err != nil {
		t.Fatalf("declined Sale: %v", err)
	}
	if declined.Approved() {
		t.Fatalf("decline trigger card was approved: %+v", declined)
	}
}

func TestIntegration_VaultLifecycle(t *testing.T) {
	client := integrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	customer, err := client.Customers.Create(ctx, &fluidpay.CustomerCreateRequest{
		Description: "fluidpay-go integration test",
		DefaultPayment: &fluidpay.VaultPaymentMethod{
			Card: &fluidpay.VaultCard{Number: "4111111111111111", ExpirationDate: "12/30"},
		},
		DefaultBillingAddress: &fluidpay.VaultAddress{
			FirstName: "Test", LastName: "Customer", Line1: "1 Main St",
			City: "Chicago", State: "IL", PostalCode: "60601", Country: "US",
		},
	})
	if err != nil {
		t.Fatalf("Customers.Create: %v", err)
	}
	defer func() {
		if _, err := client.Customers.Delete(context.Background(), customer.ID); err != nil {
			t.Errorf("Customers.Delete: %v", err)
		}
	}()

	if len(customer.Cards) != 1 || len(customer.Addresses) != 1 {
		t.Fatalf("unexpected customer shape: %+v", customer)
	}

	got, err := client.Customers.Get(ctx, customer.ID)
	if err != nil {
		t.Fatalf("Customers.Get: %v", err)
	}
	if got.ID != customer.ID {
		t.Fatalf("Get returned %s, want %s", got.ID, customer.ID)
	}

	tx, err := client.Transactions.Sale(ctx, &fluidpay.TransactionRequest{
		Amount:        500,
		PaymentMethod: fluidpay.PaymentMethod{Customer: &fluidpay.CustomerPayment{ID: customer.ID}},
	})
	if err != nil {
		t.Fatalf("Sale against vault: %v", err)
	}
	if !tx.Approved() {
		t.Fatalf("vault sale not approved: %d %s", tx.ResponseCode, tx.Response)
	}
	_, _ = client.Transactions.Void(ctx, tx.ID)
}

func TestIntegration_Terminals(t *testing.T) {
	client := integrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	list, err := client.Terminals.List(ctx)
	if err != nil {
		t.Fatalf("Terminals.List: %v", err)
	}
	t.Logf("%d terminals", list.TotalCount)
}
