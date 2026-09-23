package fluidpay

import (
	"testing"
	"time"
)

func TestAddOns(t *testing.T) {
	g, c := newGateway(t)
	g.respond("POST", "/api/recurring/addon", 200, addOnFixture)
	g.respond("GET", "/api/recurring/addon/a1", 200, addOnFixture)
	g.respond("GET", "/api/recurring/addons", 200, listEnvelope(`[{"id":"a1","name":"x"},{"id":"a2"}]`, 2))
	g.respond("POST", "/api/recurring/addon/a1", 200, addOnFixture)
	g.respond("DELETE", "/api/recurring/addon/a1", 200, `{"status":"success","msg":"success","data":null}`)

	addOn, err := c.AddOns.Create(ctx(), &AdjustmentRequest{Name: "test addon", Description: "just a simple add-on", Amount: 100})
	mustNoError(t, err)
	assertRequest(t, g, "POST", "/api/recurring/addon")
	equal(t, "body", string(g.last().Body), `{"name":"test addon","description":"just a simple add-on","amount":100,"duration":0}`)
	equal(t, "ID", addOn.ID, "b89ffdqj8m0o735i19i0")
	equal(t, "Amount", addOn.Amount, 100)
	equal(t, "null percentage", addOn.Percentage, 0)
	equal(t, "CreatedAt", addOn.CreatedAt.Year(), 2017)

	_, err = c.AddOns.Get(ctx(), "a1")
	mustNoError(t, err)
	assertRequest(t, g, "GET", "/api/recurring/addon/a1")

	list, err := c.AddOns.List(ctx())
	mustNoError(t, err)
	assertRequest(t, g, "GET", "/api/recurring/addons")
	equal(t, "total", list.TotalCount, 2)
	equal(t, "name", list.Data[0].Name, "x")

	_, err = c.AddOns.Update(ctx(), "a1", &AdjustmentRequest{Name: "n", Percentage: 43440, Duration: 3})
	mustNoError(t, err)
	assertRequest(t, g, "POST", "/api/recurring/addon/a1")
	equal(t, "update body", string(g.last().Body), `{"name":"n","percentage":43440,"duration":3}`)

	mustNoError(t, errOf(c.AddOns.Delete(ctx(), "a1")))
	assertRequest(t, g, "DELETE", "/api/recurring/addon/a1")

	_, err = c.AddOns.Create(ctx(), nil)
	mustError(t, err, "must not be nil")
	_, err = c.AddOns.Get(ctx(), "")
	mustError(t, err, "add-on id is required")
	_, err = c.AddOns.Update(ctx(), "a1", nil)
	mustError(t, err, "must not be nil")
	_, err = c.AddOns.Update(ctx(), "", &AdjustmentRequest{})
	mustError(t, err, "add-on id is required")
	mustError(t, errOf(c.AddOns.Delete(ctx(), "")), "add-on id is required")
}

func TestDiscounts(t *testing.T) {
	g, c := newGateway(t)
	g.respond("POST", "/api/recurring/discount", 200, discountFixture)
	g.respond("GET", "/api/recurring/discount/d1", 200, discountFixture)
	g.respond("GET", "/api/recurring/discounts", 200, listEnvelope(`[{"id":"d1"}]`, 1))
	g.respond("POST", "/api/recurring/discount/d1", 200, discountFixture)
	g.respond("DELETE", "/api/recurring/discount/d1", 200, `{"status":"success","msg":"success","data":null}`)

	d, err := c.Discounts.Create(ctx(), &AdjustmentRequest{Name: "test discount", Percentage: 10000, Duration: 3})
	mustNoError(t, err)
	assertRequest(t, g, "POST", "/api/recurring/discount")
	equal(t, "ID", d.ID, "b89flfqj8m0o735i19ig")
	equal(t, "Percentage", d.Percentage, 10000)
	equal(t, "null amount", d.Amount, 0)
	equal(t, "Duration", d.Duration, 3)

	_, err = c.Discounts.Get(ctx(), "d1")
	mustNoError(t, err)
	assertRequest(t, g, "GET", "/api/recurring/discount/d1")

	list, err := c.Discounts.List(ctx())
	mustNoError(t, err)
	assertRequest(t, g, "GET", "/api/recurring/discounts")
	equal(t, "len", len(list.Data), 1)

	_, err = c.Discounts.Update(ctx(), "d1", &AdjustmentRequest{Name: "n", Amount: 5})
	mustNoError(t, err)
	assertRequest(t, g, "POST", "/api/recurring/discount/d1")

	mustNoError(t, errOf(c.Discounts.Delete(ctx(), "d1")))
	assertRequest(t, g, "DELETE", "/api/recurring/discount/d1")

	_, err = c.Discounts.Create(ctx(), nil)
	mustError(t, err, "must not be nil")
	_, err = c.Discounts.Get(ctx(), "")
	mustError(t, err, "discount id is required")
	_, err = c.Discounts.Update(ctx(), "d1", nil)
	mustError(t, err, "must not be nil")
	mustError(t, errOf(c.Discounts.Delete(ctx(), "")), "discount id is required")
}

func TestPlans(t *testing.T) {
	g, c := newGateway(t)
	g.respond("POST", "/api/recurring/plan", 200, planFixture)
	g.respond("GET", "/api/recurring/plan/p1", 200, planFixture)
	g.respond("GET", "/api/recurring/plans", 200, listEnvelope(`[{"id":"p1","name":"test plan"}]`, 1))
	g.respond("POST", "/api/recurring/plan/p1", 200, planFixture)
	g.respond("DELETE", "/api/recurring/plan/p1", 200, `{"status":"success","msg":"success","data":null}`)

	plan, err := c.Plans.Create(ctx(), &PlanRequest{
		Name:                 "test plan",
		Amount:               100,
		BillingCycleInterval: 1,
		BillingFrequency:     BillingTwiceMonthly,
		BillingDays:          "1,15",
		AddOns:               []PlanAdjustment{{ID: "b75cvl51tlv38t0o7o30", Amount: 100}},
		Discounts:            []PlanAdjustment{{ID: "b89flfqj8m0o735i19ig"}},
	})
	mustNoError(t, err)
	assertRequest(t, g, "POST", "/api/recurring/plan")
	equal(t, "body", string(g.last().Body), `{"name":"test plan","amount":100,"billing_cycle_interval":1,"billing_frequency":"twice_monthly","billing_days":"1,15","duration":0,"add_ons":[{"id":"b75cvl51tlv38t0o7o30","amount":100}],"discounts":[{"id":"b89flfqj8m0o735i19ig"}]}`)
	equal(t, "ID", plan.ID, "b89g35qj8m0o735i19jg")
	equal(t, "TotalAddOns", plan.TotalAddOns, 100)
	equal(t, "add-on", plan.AddOns[0].Name, "test_addon")
	equal(t, "add-on null created_at", plan.AddOns[0].CreatedAt.IsZero(), true)
	equal(t, "discount", plan.Discounts[0].Amount, 50)

	_, err = c.Plans.Get(ctx(), "p1")
	mustNoError(t, err)
	assertRequest(t, g, "GET", "/api/recurring/plan/p1")

	list, err := c.Plans.List(ctx())
	mustNoError(t, err)
	assertRequest(t, g, "GET", "/api/recurring/plans")
	equal(t, "name", list.Data[0].Name, "test plan")

	_, err = c.Plans.Update(ctx(), "p1", &PlanRequest{Name: "update plan", Amount: 200, ChargeOnDay: true})
	mustNoError(t, err)
	assertRequest(t, g, "POST", "/api/recurring/plan/p1")
	equal(t, "charge_on_day", g.bodyJSON()["charge_on_day"], true)

	mustNoError(t, errOf(c.Plans.Delete(ctx(), "p1")))
	assertRequest(t, g, "DELETE", "/api/recurring/plan/p1")

	_, err = c.Plans.Create(ctx(), nil)
	mustError(t, err, "must not be nil")
	_, err = c.Plans.Get(ctx(), "")
	mustError(t, err, "plan id is required")
	_, err = c.Plans.Update(ctx(), "p1", nil)
	mustError(t, err, "must not be nil")
	mustError(t, errOf(c.Plans.Delete(ctx(), "")), "plan id is required")
}

func TestSubscriptions(t *testing.T) {
	g, c := newGateway(t)
	g.respond("POST", "/api/recurring/subscription", 200, subscriptionFixture)
	g.respond("GET", "/api/recurring/subscription/s1", 200, subscriptionFixture)
	g.respond("POST", "/api/recurring/subscription/search", 200, listEnvelope(`[{"id":"s1","status":"failing"},{"id":"s2","status":"active"}]`, 2))
	g.respond("POST", "/api/recurring/subscription/s1", 200, subscriptionFixture)
	g.respond("DELETE", "/api/recurring/subscription/s1", 200, `{"status":"success","msg":"success","data":null}`)

	sub, err := c.Subscriptions.Create(ctx(), &SubscriptionRequest{
		PlanID:               "b89g35qj8m0o735i19jg",
		Description:          "some description",
		Customer:             SubscriptionCustomer{ID: "b81ko5qq9qq5v460r9i0"},
		Amount:               100,
		BillingCycleInterval: 1,
		BillingFrequency:     BillingTwiceMonthly,
		BillingDays:          "1,15",
		NextBillDate:         "2017-11-22",
		Discounts:            []PlanAdjustment{{ID: "b89flfqj8m0o735i19ig"}},
	})
	mustNoError(t, err)
	assertRequest(t, g, "POST", "/api/recurring/subscription")
	equal(t, "body", string(g.last().Body), `{"plan_id":"b89g35qj8m0o735i19jg","description":"some description","customer":{"id":"b81ko5qq9qq5v460r9i0"},"amount":100,"billing_cycle_interval":1,"billing_frequency":"twice_monthly","billing_days":"1,15","next_bill_date":"2017-11-22","discounts":[{"id":"b89flfqj8m0o735i19ig"}]}`)
	equal(t, "ID", sub.ID, "b89gftaj8m0oft7upk80")
	equal(t, "Status", sub.Status, SubscriptionActive)
	equal(t, "Customer", sub.Customer.PaymentMethodType, "card")
	equal(t, "NextBillDate", sub.NextBillDate, "2017-11-22")
	equal(t, "null add_ons", len(sub.AddOns), 0)
	equal(t, "discount", sub.Discounts[0].Amount, 50)
	equal(t, "NeedsAttention", sub.NeedsAttention(), false)

	_, err = c.Subscriptions.Get(ctx(), "s1")
	mustNoError(t, err)
	assertRequest(t, g, "GET", "/api/recurring/subscription/s1")

	list, err := c.Subscriptions.Search(ctx(), &SubscriptionSearchRequest{
		PlanID:   Equal("p1"),
		Customer: &SubscriptionCustomerQuery{ID: Equal("c1")},
		Status:   Equal(SubscriptionFailing),
		Limit:    10,
	})
	mustNoError(t, err)
	assertRequest(t, g, "POST", "/api/recurring/subscription/search")
	equal(t, "search body", string(g.last().Body), `{"plan_id":{"operator":"=","value":"p1"},"customer":{"id":{"operator":"=","value":"c1"}},"status":{"operator":"=","value":"failing"},"limit":10}`)
	equal(t, "total", list.TotalCount, 2)
	equal(t, "failing needs attention", list.Data[0].NeedsAttention(), true)
	equal(t, "active fine", list.Data[1].NeedsAttention(), false)

	_, err = c.Subscriptions.Search(ctx(), nil)
	mustNoError(t, err)
	equal(t, "nil search body", string(g.last().Body), "{}")

	_, err = c.Subscriptions.Update(ctx(), "s1", &SubscriptionRequest{Customer: SubscriptionCustomer{ID: "c"}, Amount: 200, LineItems: []LineItem{{Name: "Widget", Quantity: 2, UnitPrice: 100}}, DeriveAmountFromLineItems: true})
	mustNoError(t, err)
	assertRequest(t, g, "POST", "/api/recurring/subscription/s1")
	equal(t, "derive flag", g.bodyJSON()["derive_amount_from_line_items"], true)

	mustNoError(t, errOf(c.Subscriptions.Delete(ctx(), "s1")))
	assertRequest(t, g, "DELETE", "/api/recurring/subscription/s1")

	_, err = c.Subscriptions.Create(ctx(), nil)
	mustError(t, err, "must not be nil")
	_, err = c.Subscriptions.Get(ctx(), "")
	mustError(t, err, "subscription id is required")
	_, err = c.Subscriptions.Update(ctx(), "s1", nil)
	mustError(t, err, "must not be nil")
	mustError(t, errOf(c.Subscriptions.Delete(ctx(), "")), "subscription id is required")
}

func TestSubscriptions_Lifecycle(t *testing.T) {
	g, c := newGateway(t)
	for _, status := range []string{"paused", "past_due", "cancelled", "completed", "active"} {
		g.respond("GET", "/api/recurring/subscription/s1/status/"+status, 200, subscriptionFixture)
	}

	calls := []struct {
		name string
		call func() (*Subscription, error)
		path string
	}{
		{"Pause", func() (*Subscription, error) { return c.Subscriptions.Pause(ctx(), "s1") }, "/api/recurring/subscription/s1/status/paused"},
		{"MarkPastDue", func() (*Subscription, error) { return c.Subscriptions.MarkPastDue(ctx(), "s1") }, "/api/recurring/subscription/s1/status/past_due"},
		{"Cancel", func() (*Subscription, error) { return c.Subscriptions.Cancel(ctx(), "s1") }, "/api/recurring/subscription/s1/status/cancelled"},
		{"Complete", func() (*Subscription, error) { return c.Subscriptions.Complete(ctx(), "s1") }, "/api/recurring/subscription/s1/status/completed"},
		{"Activate", func() (*Subscription, error) { return c.Subscriptions.Activate(ctx(), "s1", time.Time{}) }, "/api/recurring/subscription/s1/status/active"},
	}
	for _, tc := range calls {
		t.Run(tc.name, func(t *testing.T) {
			sub, err := tc.call()
			mustNoError(t, err)
			r := assertRequest(t, g, "GET", tc.path)
			equal(t, "no query", r.Query, "")
			equal(t, "ID", sub.ID, "b89gftaj8m0oft7upk80")
		})
	}

	_, err := c.Subscriptions.Activate(ctx(), "s1", time.Date(2026, 10, 1, 15, 0, 0, 0, time.UTC))
	mustNoError(t, err)
	r := assertRequest(t, g, "GET", "/api/recurring/subscription/s1/status/active")
	equal(t, "next_bill_date", r.Query, "next_bill_date=2026-10-01")

	_, err = c.Subscriptions.Pause(ctx(), "")
	mustError(t, err, "subscription id is required")
}
