// Command subscription creates a plan and subscribes a vault customer to it,
// then pauses, reactivates and finally cancels the subscription.
//
//	FLUIDPAY_API_KEY=api_... FLUIDPAY_ENVIRONMENT=sandbox go run ./examples/subscription
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
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	customer, err := client.Customers.Create(ctx, &fluidpay.CustomerCreateRequest{
		Description: "Subscriber",
		DefaultPayment: &fluidpay.VaultPaymentMethod{
			Card: &fluidpay.VaultCard{Number: "4111111111111111", ExpirationDate: "12/30"},
		},
	})
	if err != nil {
		log.Fatalf("create customer: %v", err)
	}
	defer client.Customers.Delete(context.Background(), customer.ID) //nolint:errcheck

	plan, err := client.Plans.Create(ctx, &fluidpay.PlanRequest{
		Name:                 "Pro monthly",
		Description:          "Billed on the 1st",
		Amount:               4900,
		BillingCycleInterval: 1,
		BillingFrequency:     fluidpay.BillingMonthly,
		BillingDays:          "1",
	})
	if err != nil {
		log.Fatalf("create plan: %v", err)
	}
	defer client.Plans.Delete(context.Background(), plan.ID) //nolint:errcheck
	fmt.Printf("plan %s: %d cents %s\n", plan.ID, plan.Amount, plan.BillingFrequency)

	sub, err := client.Subscriptions.Create(ctx, &fluidpay.SubscriptionRequest{
		PlanID:               plan.ID,
		Description:          "Jane's Pro subscription",
		Customer:             fluidpay.SubscriptionCustomer{ID: customer.ID},
		Amount:               plan.Amount,
		BillingCycleInterval: plan.BillingCycleInterval,
		BillingFrequency:     plan.BillingFrequency,
		BillingDays:          plan.BillingDays,
		NextBillDate:         time.Now().AddDate(0, 1, 0).Format("2006-01-02"),
	})
	if err != nil {
		log.Fatalf("create subscription: %v", err)
	}
	fmt.Printf("subscription %s is %s, next bill %s\n", sub.ID, sub.Status, sub.NextBillDate)

	if sub, err = client.Subscriptions.Pause(ctx, sub.ID); err != nil {
		log.Fatalf("pause: %v", err)
	}
	fmt.Println("paused:", sub.Status)

	if sub, err = client.Subscriptions.Activate(ctx, sub.ID, time.Now().AddDate(0, 2, 0)); err != nil {
		log.Fatalf("activate: %v", err)
	}
	fmt.Println("active again, next bill", sub.NextBillDate)

	if sub, err = client.Subscriptions.Cancel(ctx, sub.ID); err != nil {
		log.Fatalf("cancel: %v", err)
	}
	fmt.Println("cancelled:", sub.Status)

	// Find everything that needs attention on this plan.
	failing, err := client.Subscriptions.Search(ctx, &fluidpay.SubscriptionSearchRequest{
		PlanID: fluidpay.Equal(plan.ID),
		Limit:  100,
	})
	if err != nil {
		log.Fatalf("search: %v", err)
	}
	for _, s := range failing.Data {
		if s.NeedsAttention() {
			fmt.Printf("subscription %s needs attention: %s\n", s.ID, s.Status)
		}
	}
	_, _ = client.Subscriptions.Delete(ctx, sub.ID)
}
