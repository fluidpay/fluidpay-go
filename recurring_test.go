package fluidpay

import (
	"context"
	"net/http"
	"testing"
)

var creAonReq = RecurrenceRequest{
	Name:        "asdQWE123",
	Description: "simple add on",
	Amount:      12,
	Duration:    23,
}

var updAonReq = RecurrenceRequest{
	Name:        "321EWQdsa",
	Description: "update add on",
	Amount:      23,
	Duration:    34,
}

var creDisReq = RecurrenceRequest{
	Name:        "jlkUIO789",
	Description: "simple discount",
	Amount:      98,
	Duration:    87,
}

var updDisReq = RecurrenceRequest{
	Name:        "789UIOjkl",
	Description: "update discount",
	Amount:      87,
	Duration:    76,
}

func TestAddOn(t *testing.T) {
	fluidpay := Fluidpay{
		APIKey:   TestAPIKey,
		Client:   new(http.Client),
		Ctx:      context.Background(),
		LocalDev: true,
	}

	creAonRes, err := CreateAddOn(fluidpay, creAonReq)
	if err != nil {
		t.Error(err)
	}
	if creAonRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", creAonRes.Msg)
	}

	id := creAonRes.Data.ID

	getAonRes, err := GetAddOn(fluidpay, id)
	if err != nil {
		t.Error(err)
	}
	if getAonRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", getAonRes.Msg)
	}

	getAonsRes, err := GetAddOns(fluidpay)
	if err != nil {
		t.Error(err)
	}
	if getAonsRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", getAonsRes.Msg)
	}
	if getAonsRes.TotalCount == 0 {
		t.Errorf("TotalCount shouldn't be '%v'", getAonsRes.TotalCount)
	}

	updAonRes, err := UpdateAddOn(fluidpay, updAonReq, id)
	if err != nil {
		t.Error(err)
	}
	if updAonRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", updAonRes.Msg)
	}

	delAonRes, err := DeleteAddOn(fluidpay, id)
	if err != nil {
		t.Error(err)
	}
	if delAonRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", delAonRes.Msg)
	}
}

func TestDiscount(t *testing.T) {
	fluidpay := Fluidpay{
		APIKey:   TestAPIKey,
		Client:   new(http.Client),
		Ctx:      context.Background(),
		LocalDev: true,
	}

	creDisRes, err := CreateDiscount(fluidpay, creDisReq)
	if err != nil {
		t.Error(err)
	}
	if creDisRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", creDisRes.Msg)
	}

	id := creDisRes.Data.ID

	getDisRes, err := GetDiscount(fluidpay, id)
	if err != nil {
		t.Error(err)
	}
	if getDisRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", getDisRes.Msg)
	}

	getDissRes, err := GetDiscounts(fluidpay)
	if err != nil {
		t.Error(err)
	}
	if getDissRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", getDissRes.Msg)
	}
	if getDissRes.TotalCount == 0 {
		t.Errorf("TotalCount shouldn't be '%v'", getDissRes.TotalCount)
	}

	updDisRes, err := UpdateDiscount(fluidpay, updDisReq, id)
	if err != nil {
		t.Error(err)
	}
	if updDisRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", updDisRes.Msg)
	}

	delDisRes, err := DeleteDiscount(fluidpay, id)
	if err != nil {
		t.Error(err)
	}
	if delDisRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", delDisRes.Msg)
	}
}

func TestPlanAndSubscription(t *testing.T) {
	fluidpay := Fluidpay{
		APIKey:   TestAPIKey,
		Client:   new(http.Client),
		Ctx:      context.Background(),
		LocalDev: true,
	}

	creAonRes, err := CreateAddOn(fluidpay, creAonReq)
	creDisRes, err := CreateDiscount(fluidpay, creDisReq)
	creCusRes, err := CreateCustomer(fluidpay, creCusReq)

	aonID := creAonRes.Data.ID
	disID := creDisRes.Data.ID
	cusID := creCusRes.Data.ID

	crePlaReq := PlanRequest{
		Name:                 "test plan",
		Description:          "just a simple test plan",
		Amount:               88,
		BillingCycleInterval: 1,
		BillingFrequency:     "twice_monthly",
		BillingDays:          "1,15",
		Duration:             0,
		AddOns: []OptionalRecurringRequest{
			{
				ID:          aonID,
				Description: "this will add to the cost of the subscription",
				Amount:      33,
				Duration:    0,
			},
		},
		Discounts: []OptionalRecurringRequest{
			{
				ID:          disID,
				Description: "this will discount the cost of the subscription",
				Amount:      22,
				Duration:    0,
			},
		},
	}

	updPlaReq := PlanRequest{
		Name:                 "update plan",
		Description:          "just a simple test plan",
		Amount:               100,
		BillingCycleInterval: 1,
		BillingFrequency:     "twice_monthly",
		BillingDays:          "1,15",
		Duration:             0,
		AddOns: []OptionalRecurringRequest{
			{
				ID:          aonID,
				Description: "this will update the cost of the subscription",
				Amount:      100,
				Duration:    0,
			},
		},
		Discounts: []OptionalRecurringRequest{
			{
				ID:          disID,
				Description: "this will update the cost of the subscription",
				Amount:      50,
				Duration:    0,
			},
		},
	}

	crePlaRes, err := CreatePlan(fluidpay, crePlaReq)
	if err != nil {
		t.Error(err)
	}
	if crePlaRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", crePlaRes.Msg)
	}

	plaID := crePlaRes.Data.ID

	getPlaRes, err := GetPlan(fluidpay, plaID)
	if err != nil {
		t.Error(err)
	}
	if getPlaRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", getPlaRes.Msg)
	}

	getPlasRes, err := GetPlans(fluidpay)
	if err != nil {
		t.Error(err)
	}
	if getPlasRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", getPlasRes.Msg)
	}
	if getPlasRes.TotalCount == 0 {
		t.Errorf("TotalCount shouldn't be '%v'", getPlasRes.TotalCount)
	}

	updPlaRes, err := UpdatePlan(fluidpay, updPlaReq, plaID)
	if err != nil {
		t.Error(err)
	}
	if updPlaRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", updPlaRes.Msg)
	}

	// Subscription section.

	creSubReq := SubscriptionRequest{
		PlanID:      plaID,
		Description: "some description to describe the subscription",
		Customer: IDData{
			ID: cusID,
		},
		Amount:               100,
		BillingCycleInterval: 1,
		BillingFrequency:     "twice_monthly",
		BillingDays:          "1,15",
		Duration:             2,
		NextBillDate:         "2018-11-21",
		Discounts: []OptionalRecurringRequest{
			{
				ID: disID,
			},
		},
	}

	updSubReq := SubscriptionRequest{
		PlanID:      plaID,
		Description: "some update to the subscription",
		Customer: IDData{
			ID: cusID,
		},
		Amount:               100,
		BillingCycleInterval: 1,
		BillingFrequency:     "twice_monthly",
		BillingDays:          "1,15",
		Duration:             3,
		NextBillDate:         "2018-11-25",
		Discounts: []OptionalRecurringRequest{
			{
				ID:          disID,
				Name:        "test discount",
				Description: "this will discount the cost of the subscription",
				Amount:      50,
				Duration:    5,
			},
		},
	}

	creSubRes, err := CreateSubscription(fluidpay, creSubReq)
	if err != nil {
		t.Error(err)
	}
	if creSubRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", creSubRes.Msg)
	}

	subID := creSubRes.Data.ID

	getSubRes, err := GetSubscription(fluidpay, subID)
	if err != nil {
		t.Error(err)
	}
	if getSubRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", getSubRes.Msg)
	}

	updSubRes, err := UpdateSubscription(fluidpay, updSubReq, subID)
	if err != nil {
		t.Error(err)
	}
	if updSubRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", updSubRes.Msg)
	}

	delSubRes, err := DeleteSubscription(fluidpay, subID)
	if err != nil {
		t.Error(err)
	}
	if delSubRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", delSubRes.Msg)
	}

	delPlaRes, err := DeletePlan(fluidpay, plaID)
	if err != nil {
		t.Error(err)
	}
	if delPlaRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", delPlaRes.Msg)
	}

	DeleteAddOn(fluidpay, aonID)
	DeleteDiscount(fluidpay, disID)
	DeleteCustomer(fluidpay, cusID)
}
