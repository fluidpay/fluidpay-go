package fluidpay

import (
	"context"
	"net/http"
	"testing"
)

var addr = Address{
	FirstName:    "John",
	LastName:     "Smith",
	Company:      "Test Company",
	AddressLine1: "123 Some St",
	City:         "West Chicago",
	State:        "IL",
	PostalCode:   "60185",
	Country:      "US",
	Phone:        "5555555555",
	Fax:          "5555555555",
	Email:        "help@website.com",
}

var cardTransReq = TransactionRequest{
	Type:           "authorize",
	Amount:         1112,
	TaxAmount:      100,
	ShippingAmount: 100,
	Currency:       "USD",
	Description:    "test transaction",
	OrderID:        "someOrderID",
	PoNumber:       "somePONumber",
	IPAddress:      "4.2.2.2",
	EmailReciept:   true,
	EmailAddress:   "user@home.com",
	PaymentMethod: PaymentMethodsRequest{
		Card: &CardRequest{
			EntryType:      "keyed",
			Number:         "4012000098765439",
			ExpirationDate: "12/20",
			Cvc:            "999",
			CardholderAuthentication: CardholderAuthentication{
				Condition: "",
				Eci:       "",
				Cavv:      "",
				Xid:       "",
			},
		},
	},
	BillingAddress:  addr,
	ShippingAddress: addr,
}

var queTransReq = TransactionQueryRequest{
	TransactionID: StringQuery{
		Operator: "!=",
		Value:    "b84vgb2j8m0jujadi4v0",
	},
	Amount: IntQuery{
		Operator: "<",
		Value:    1,
	},
}

func TestTransaction(t *testing.T) {
	fluidpay := Fluidpay{
		APIKey:   TestAPIKey,
		Client:   new(http.Client),
		Ctx:      context.Background(),
		LocalDev: true,
	}

	/* cardTransRes, err := DoTransaction(fluidpay, cardTransReq)
	if err != nil {
		t.Error(err)
	}
	if cardTransRes.Data.PaymentMethod != "card" {
		t.Errorf("PaymentMethod should be 'card', instead of '%v'", cardTransRes.Data.PaymentMethod)
	} */

	cusRes, err := CreateCustomer(fluidpay, creCusReq)
	cusID := cusRes.Data.ID

	cusAdrRes, err := CreateCustomerAddress(fluidpay, creCusAdrReg, cusID)
	cusAdrID := cusAdrRes.Data.ID

	cusPayRes, err := CreateCustomerPayment(fluidpay, creCusPayReq, cusID, "card")
	cusPayID := cusPayRes.Data.Card.ID

	var cusTransReq = TransactionRequest{
		Type:           "authorize",
		Amount:         1112,
		TaxAmount:      100,
		ShippingAmount: 100,
		Currency:       "USD",
		Description:    "test transaction",
		OrderID:        "someOrderID",
		PoNumber:       "somePONumber",
		IPAddress:      "4.2.2.2",
		EmailReciept:   true,
		EmailAddress:   "user@home.com",
		PaymentMethod: PaymentMethodsRequest{
			Customer: &CustomerRequest{
				ID:                cusID,
				PaymentMethodType: "card",
				PaymentMethodID:   cusPayID,
				BillingAddressID:  cusAdrID,
				ShippingAddressID: cusAdrID,
			},
		},
		BillingAddress:  addr,
		ShippingAddress: addr,
	}

	cusTransRes, err := DoTransaction(fluidpay, cusTransReq)
	if err != nil {
		t.Error(err)
	}
	if cusTransRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", cusTransRes.Msg)
	}
	cusTransID := cusTransRes.Data.ID

	termRes, err := GetTerminals(fluidpay)
	termID := termRes.Data[0].ID

	var termTransReq = TransactionRequest{
		Type:           "sale",
		Amount:         1112,
		TaxExempt:      true,
		TaxAmount:      100,
		ShippingAmount: 100,
		Currency:       "USD",
		Description:    "test transaction",
		OrderID:        "someOrderID",
		PoNumber:       "somePONumber",
		IPAddress:      "4.2.2.2",
		EmailReciept:   true,
		EmailAddress:   "user@home.com",
		PaymentMethod: PaymentMethodsRequest{
			Terminal: &TerminalRequest{
				ID:                termID,
				ExpirationDate:    "12/20",
				Cvc:               "555",
				PrintReceipt:      "both",
				SignatureRequired: true,
			},
		},
		BillingAddress:  addr,
		ShippingAddress: addr,
	}

	termTransRes, err := DoTransaction(fluidpay, termTransReq)
	if err != nil {
		t.Error(err)
	}
	if termTransRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", termTransRes.Msg)
	}
	termTransID := termTransRes.Data.ID

	/* queTransRes, err := QueryTransactions(fluidpay, queTransReq)
	if err != nil {
		t.Error(err)
	}
	if queTransRes.Status != "success" {
		t.Errorf("Status should be 'success', instead of '%v'", queTransRes.Status)
	} */

	voidTransRes, err := VoidTransaction(fluidpay, termTransID)
	if err != nil {
		t.Error(err)
	}
	if voidTransRes.Status != "success" {
		t.Errorf("Status should be 'success', instead of '%v'", voidTransRes.Status)
	}

	data := cusTransRes.Data
	capTransReq := TransactionCaptureRequest{
		Amount:         data.Amount,
		TaxAmount:      data.TaxAmount,
		TaxExempt:      data.TaxExempt,
		ShippingAmount: data.ShippingAmount,
		OrderID:        data.OrderID,
		PoNumber:       data.PoNumber,
		IPAddress:      data.IPAddress,
	}
	capTransRes, err := CaptureTransaction(fluidpay, capTransReq, cusTransID)
	if err != nil {
		t.Error(err)
	}
	if capTransRes.Status != "success" {
		t.Errorf("Status should be 'success', instead of '%v'", capTransRes.Status)
	}

	/* refTransReq := TransactionRefundRequest{
		Amount: data.Amount,
	}
	refTransRes, err := RefundTransaction(fluidpay, refTransReq, data.ID)
	if err != nil {
		t.Error(err)
	}
	if refTransRes.Status != "success" {
		t.Errorf("Status should be 'success', instead of '%v'", refTransRes.Status)
	} */
}
