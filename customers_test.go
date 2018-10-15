package fluidpay

import (
	"context"
	"net/http"
	"testing"
)

var creCusReq = CreateCustomerRequest{
	Description: "test description",
	PaymentMethod: CustomerPaymentMethodRequest{
		Card: CustomerPaymentRequest{
			CardNumber:     "4111111111111111",
			ExpirationDate: "12/20",
			EntryType:      "keyed",
		},
	},
	BillingAddress: Address{
		FirstName:    "John",
		LastName:     "Smith",
		Company:      "Some Business",
		AddressLine1: "123 Some St",
		AddressLine2: "STE 1",
		City:         "Some Town",
		State:        "IL",
		PostalCode:   "60185",
		Country:      "US",
		Email:        "user@somesite.com",
		Phone:        "5555555555",
		Fax:          "555555555",
	},
	ShippingAddress: Address{
		FirstName:    "John",
		LastName:     "Smith",
		Company:      "Some Business",
		AddressLine1: "123 Some St",
		AddressLine2: "STE 1",
		City:         "Some Town",
		State:        "IL",
		PostalCode:   "60185",
		Country:      "US",
		Email:        "user@somesite.com",
		Phone:        "5555555555",
		Fax:          "555555555",
	},
}

var creCusAdrReg = CustomerAddressRequest{
	FirstName:    "John",
	LastName:     "Smith",
	Company:      "Some Business",
	AddressLine1: "123 Some St",
	AddressLine2: "STE 1",
	City:         "Some Town",
	State:        "IL",
	PostalCode:   "60185",
	Country:      "US",
	Email:        "user@somesite.com",
	Phone:        "5555555555",
	Fax:          "555555555",
}

var updCusAdrReg = CustomerAddressRequest{
	FirstName:    "Jane",
	LastName:     "Taylor",
	Company:      "Some Business",
	AddressLine1: "123 Some St",
	AddressLine2: "STE 1",
	City:         "Some Town",
	State:        "IL",
	PostalCode:   "60185",
	Country:      "US",
	Email:        "user@somesite.com",
	Phone:        "5555555555",
	Fax:          "555555555",
}

var creCusPayReq = CustomerPaymentRequest{
	CardNumber:     "4111111111111111",
	ExpirationDate: "12/20",
	EntryType:      "keyed",
}

var updCusPayReq = CustomerPaymentRequest{
	CardNumber:     "4111111111111555",
	ExpirationDate: "12/21",
	EntryType:      "keyed",
}

func TestCustomer(t *testing.T) {
	fluidpay := Fluidpay{
		APIKey:   TestAPIKey,
		Client:   new(http.Client),
		Ctx:      context.Background(),
		LocalDev: true,
	}

	creCusRes, err := CreateCustomer(fluidpay, creCusReq)
	if err != nil {
		t.Error(err)
	}
	if creCusRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", creCusRes.Msg)
	}

	cusID := creCusRes.Data.ID

	getCusRes, err := GetCustomer(fluidpay, cusID)
	if err != nil {
		t.Error(err)
	}
	if getCusRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", getCusRes)
	}

	// Address token section.

	creCusAdrRes, err := CreateCustomerAddress(fluidpay, creCusAdrReg, cusID)
	if err != nil {
		t.Error(err)
	}
	if creCusAdrRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", creCusAdrRes.Msg)
	}
	if creCusAdrRes.Data.FirstName != "John" {
		t.Errorf("FirstName should be 'John', instead of '%v'", creCusAdrRes.Data.FirstName)
	}

	adrID := creCusAdrRes.Data.ID

	getAllCusAdrRes, err := GetCustomerAddresses(fluidpay, cusID)
	if err != nil {
		t.Error(err)
	}
	if getAllCusAdrRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", getAllCusAdrRes.Msg)
	}
	if getAllCusAdrRes.TotalCount == 0 {
		t.Errorf("Count shouldn't be '%v'", getAllCusAdrRes.TotalCount)
	}

	getCusAdrRes, err := GetCustomerAddress(fluidpay, cusID, adrID)
	if err != nil {
		t.Error(err)
	}
	if getCusAdrRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", getCusAdrRes.Msg)
	}

	updCusAdrRes, err := UpdateCustomerAddress(fluidpay, updCusAdrReg, cusID, adrID)
	if err != nil {
		t.Error(err)
	}
	if updCusAdrRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", updCusAdrRes.Msg)
	}

	// Payment token section.

	creCusPayRes, err := CreateCustomerPayment(fluidpay, creCusPayReq, cusID, "card")
	if err != nil {
		t.Error(err)
	}
	if creCusPayRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", creCusPayRes.Msg)
	}

	payID := creCusPayRes.Data.Card.ID

	getAllCusPayRes, err := GetCustomerPayments(fluidpay, cusID, "card")
	if err != nil {
		t.Error(err)
	}
	if getAllCusPayRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", getAllCusPayRes.Msg)
	}
	if getAllCusPayRes.TotalCount == 0 {
		t.Errorf("Count shouldn't be '%v'", getAllCusPayRes.TotalCount)
	}

	getCusPayRes, err := GetCustomerPayment(fluidpay, cusID, "card", payID)
	if err != nil {
		t.Error(err)
	}
	if getCusPayRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", getCusPayRes)
	}

	// @TODO find out how the Update request works.

	/* updCusPayRes, err := UpdateCustomerPayment(fluidpay, updCusPayReq, cusID, "card", payID)
	if err != nil {
		t.Error(err)
	}
	if updCusPayRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", updCusPayRes.Msg)
	}
	*/

	// Customer update.

	updCusReq := UpdateCustomerRequest{
		Description:       "test update description",
		PaymentMethod:     "card",
		PaymentMethodID:   payID,
		BillingAddressID:  adrID,
		ShippingAddressID: adrID,
	}

	updCusRes, err := UpdateCustomer(fluidpay, updCusReq, cusID)
	if err != nil {
		t.Error(err)
	}
	if updCusRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", updCusRes.Msg)
	}

	// Token cleanup.

	// Token's cannot be deleted, because they are being in use.
	/* delCusAdrRes, err := DeleteCustomerAddress(fluidpay, cusID, adrID)
	if err != nil {
		t.Error(err)
	}
	if delCusAdrRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", delCusAdrRes.Msg)
	}

	delCusPayRes, err := DeleteCustomerPayment(fluidpay, cusID, "card", payID)
	if err != nil {
		t.Error(err)
	}
	if delCusPayRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", delCusPayRes.Msg)
	}
	*/

	delCusRes, err := DeleteCustomer(fluidpay, cusID)
	if err != nil {
		t.Error(err)
	}
	if delCusRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", delCusRes.Msg)
	}
}
