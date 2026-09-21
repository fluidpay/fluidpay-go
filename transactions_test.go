package fluidpay

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func jsonMarshal(v any) ([]byte, error) { return json.Marshal(v) }

func TestTransactions_Sale(t *testing.T) {
	g, c := newGateway(t)
	g.respond("POST", "/api/transaction", 200, cardSaleFixture)

	req := &TransactionRequest{
		Amount:         1299,
		OrderID:        "someOrderID",
		IdempotencyKey: "7df4d862-1a3d-44c4-b3df-536aadf307b0",
		PaymentMethod: PaymentMethod{Card: &CardPayment{
			Number: "4111111111111111", ExpirationDate: "12/30", CVC: "123",
		}},
		BillingAddress: &Address{PostalCode: "60187", Country: "US"},
		CustomFields:   flatMap{"cf_1": "value"}.toMulti(),
	}
	tx, err := c.Transactions.Sale(ctx(), req)
	mustNoError(t, err)

	assertRequest(t, g, "POST", "/api/transaction")
	body := g.bodyJSON()
	equal(t, "type", body["type"], "sale")
	equal(t, "amount", body["amount"], float64(1299))
	equal(t, "idempotency_key", body["idempotency_key"], "7df4d862-1a3d-44c4-b3df-536aadf307b0")
	card := body["payment_method"].(map[string]any)["card"].(map[string]any)
	equal(t, "card.number", card["number"], "4111111111111111")
	equal(t, "card.cvc", card["cvc"], "123")
	if _, present := body["base_amount"]; present {
		t.Error("base_amount should be omitted when zero")
	}
	equal(t, "request type untouched", req.Type, "")

	equal(t, "ID", tx.ID, "d3bfel2cunifubaumn00")
	equal(t, "Approved", tx.Approved(), true)
	equal(t, "Declined", tx.Declined(), false)
	equal(t, "Status", tx.Status, StatusPendingSettlement)
	equal(t, "ResponseCode", tx.ResponseCode, ResponseCodeApproved)
	equal(t, "AmountCaptured", tx.AmountCaptured, 1299)
	equal(t, "Features", tx.Features[0], "avs")
	equal(t, "CustomFields", tx.CustomFields["cf_1"][0], "value")
	equal(t, "BillingAddress", tx.BillingAddress.City, "Wheaton")
	equal(t, "CreatedAt", tx.CreatedAt.UTC().Format(time.RFC3339Nano), "2025-09-26T20:30:09.693858Z")
	if tx.CapturedAt == nil || tx.SettledAt != nil {
		t.Errorf("CapturedAt=%v SettledAt=%v", tx.CapturedAt, tx.SettledAt)
	}
	if tx.ResponseBody.Card == nil {
		t.Fatal("card response missing")
	}
	equal(t, "AuthCode", tx.ResponseBody.Card.AuthCode, "TAS000")
	equal(t, "AVS", tx.ResponseBody.Card.AVSResponseCode, "Y")
	equal(t, "processor_specific string", string(tx.ResponseBody.Card.ProcessorSpecific), `""`)
	if tx.ResponseBody.ACH != nil || tx.ResponseBody.Terminal != nil {
		t.Error("other payment responses should be nil")
	}
	equal(t, "correlation id", tx.LastResponse.CorrelationID, testCorrelationID)
}

// toMulti turns a flat map into the custom_fields shape.
type flatMap map[string]string

func (m flatMap) toMulti() map[string][]string {
	out := make(map[string][]string, len(m))
	for k, v := range m {
		out[k] = []string{v}
	}
	return out
}

func TestTransactions_TypedHelpers(t *testing.T) {
	g, c := newGateway(t)
	g.respond("POST", "/api/transaction", 200, cardSaleFixture)

	base := &TransactionRequest{Amount: 100, Type: "ignored"}
	calls := []struct {
		name string
		call func() (*Transaction, error)
		want string
	}{
		{"Sale", func() (*Transaction, error) { return c.Transactions.Sale(ctx(), base) }, TransactionTypeSale},
		{"Authorize", func() (*Transaction, error) { return c.Transactions.Authorize(ctx(), base) }, TransactionTypeAuthorize},
		{"Verify", func() (*Transaction, error) { return c.Transactions.Verify(ctx(), base) }, TransactionTypeVerification},
		{"Credit", func() (*Transaction, error) { return c.Transactions.Credit(ctx(), base) }, TransactionTypeCredit},
		{"Create", func() (*Transaction, error) { return c.Transactions.Create(ctx(), base) }, "ignored"},
	}
	for _, tc := range calls {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.call()
			mustNoError(t, err)
			equal(t, "type", g.bodyJSON()["type"], tc.want)
		})
	}
	equal(t, "caller's request not mutated", base.Type, "ignored")

	for name, fn := range map[string]func() (*Transaction, error){
		"Sale":   func() (*Transaction, error) { return c.Transactions.Sale(ctx(), nil) },
		"Create": func() (*Transaction, error) { return c.Transactions.Create(ctx(), nil) },
	} {
		_, err := fn()
		mustError(t, err, "must not be nil")
		_ = name
	}
}

func TestTransactionRequest_MinimalJSON(t *testing.T) {
	req := TransactionRequest{
		Type:   TransactionTypeSale,
		Amount: 1299,
		PaymentMethod: PaymentMethod{Card: &CardPayment{
			Number: "4111111111111111", ExpirationDate: "12/30",
		}},
	}
	b, err := json.Marshal(req)
	mustNoError(t, err)
	want := `{"type":"sale","amount":1299,"payment_method":{"card":{"number":"4111111111111111","expiration_date":"12/30"}}}`
	equal(t, "json", string(b), want)
}

func TestTransactionRequest_FullJSONKeys(t *testing.T) {
	req := TransactionRequest{
		Type: "sale", Amount: 1, BaseAmount: 2, Currency: "USD", ProcessorID: "p",
		PaymentMethod:  PaymentMethod{Token: "tok"},
		IdempotencyKey: "k", IdempotencyTime: 60,
		OrderID: "o", PoNumber: "po", Description: "d", IPAddress: "1.1.1.1", VendorID: "v",
		TaxAmount: 1, TaxExempt: true, ShippingAmount: 1, DiscountAmount: 1, TipAmount: 1,
		PaymentAdjustment:    &PaymentAdjustment{Type: AdjustmentFlat, Value: 199},
		BillingAddress:       &Address{City: "x"},
		ShippingAddress:      &Address{City: "y"},
		LineItems:            []LineItem{{Name: "n", Quantity: 1.5, UnitPrice: 100}},
		SummaryCommodityCode: "1234", ShipFromPostalCode: "60601",
		CustomFields: map[string][]string{"a": {"b"}}, GroupName: "g",
		CreateVaultRecord: true, CreateVaultRecordFor: "cust",
		EmailReceipt: true, EmailAddress: "e@x.com",
		Descriptor:          &Descriptor{Name: "SHOP"},
		AllowPartialPayment: true, SplitTransactionAmount: 50,
		CardOnFileIndicator: "C", InitiatedBy: "merchant", InitialTransactionID: "t", StoredCredentialIndicator: "used", BillingMethod: "recurring",
		IIASStatus:        "verified",
		AdditionalAmounts: &AdditionalAmounts{HSA: &HSAAmounts{Total: 100}},
		ProcessorSpecific: map[string]any{"paysafe_direct": map[string]any{"subscription_amount": 1}},
	}
	b, err := json.Marshal(req)
	mustNoError(t, err)
	var m map[string]any
	mustNoError(t, json.Unmarshal(b, &m))
	want := "additional_amounts,allow_partial_payment,amount,base_amount,billing_address,billing_method,card_on_file_indicator,create_vault_record,create_vault_record_for,currency,custom_fields,description,descriptor,discount_amount,email_address,email_receipt,group_name,idempotency_key,idempotency_time,iias_status,initial_transaction_id,initiated_by,ip_address,line_items,order_id,payment_adjustment,payment_method,po_number,processor_id,processor_specific,ship_from_postal_code,shipping_address,shipping_amount,split_transaction_amount,stored_credential_indicator,summary_commodity_code,tax_amount,tax_exempt,tip_amount,type,vendor_id"
	equal(t, "keys", joinKeys(m), want)
	equal(t, "email_receipt spelled correctly", m["email_receipt"], true)
	equal(t, "token", m["payment_method"].(map[string]any)["token"], "tok")
}

func TestPaymentMethod_JSONShapes(t *testing.T) {
	pm := PaymentMethod{
		Customer: &CustomerPayment{ID: "cust", PaymentMethodID: "pm", PaymentMethodType: "card"},
	}
	b, _ := json.Marshal(pm)
	equal(t, "customer", string(b), `{"customer":{"id":"cust","payment_method_id":"pm","payment_method_type":"card"}}`)

	pm = PaymentMethod{Terminal: &TerminalPayment{ID: "term", PrintReceipt: ReceiptBoth, SignatureRequired: true}}
	b, _ = json.Marshal(pm)
	equal(t, "terminal", string(b), `{"terminal":{"id":"term","print_receipt":"both","signature_required":true}}`)

	pm = PaymentMethod{ACH: &ACHPayment{RoutingNumber: "111111111", AccountNumber: "111111111", SecCode: SecCodeWeb, AccountType: AccountTypeChecking}}
	b, _ = json.Marshal(pm)
	equal(t, "ach", string(b), `{"ach":{"routing_number":"111111111","account_number":"111111111","sec_code":"web","account_type":"checking"}}`)

	pm = PaymentMethod{GooglePayToken: json.RawMessage(`{"signature":"s"}`)}
	b, _ = json.Marshal(pm)
	equal(t, "google pay passthrough", string(b), `{"google_pay_token":{"signature":"s"}}`)

	pm = PaymentMethod{ApplePayToken: &ApplePayToken{TemporaryToken: "abc"}}
	b, _ = json.Marshal(pm)
	equal(t, "apple pay", string(b), `{"apple_pay_token":{"temporary_token":"abc"}}`)

	pm = PaymentMethod{Card: &CardPayment{Number: "4", ExpirationDate: "1/1", CardholderAuthentication: &CardholderAuthentication{ECI: "05", CAVV: "c"}}}
	b, _ = json.Marshal(pm)
	equal(t, "3ds", string(b), `{"card":{"number":"4","expiration_date":"1/1","cardholder_authentication":{"eci":"05","cavv":"c"}}}`)
}

func TestTransactions_Declined(t *testing.T) {
	g, c := newGateway(t)
	g.respond("POST", "/api/transaction", 200, declinedSaleFixture)
	tx, err := c.Transactions.Sale(ctx(), &TransactionRequest{Amount: 1299})
	mustNoError(t, err)
	equal(t, "Approved", tx.Approved(), false)
	equal(t, "Declined", tx.Declined(), true)
	equal(t, "ResponseCode", tx.ResponseCode, ResponseCodeInsufficient)
	equal(t, "Status", tx.Status, StatusDeclined)
	equal(t, "processor text", tx.ResponseBody.Card.ProcessorResponseText, "INSUFFICIENT FUNDS")
	if tx.CapturedAt != nil {
		t.Error("CapturedAt should be nil")
	}
}

func TestTransaction_ResponseHelpers(t *testing.T) {
	tests := []struct {
		code                                       int
		amount, authorized                         int
		approved, partial, declined, gateway, proc bool
	}{
		{100, 100, 100, true, false, false, false, false},
		{110, 100, 50, true, true, false, false, false},
		{100, 100, 50, true, true, false, false, false},
		{101, 100, 100, true, false, false, false, false},
		{200, 100, 0, false, false, true, false, false},
		{301, 100, 0, false, false, false, true, false},
		{421, 100, 0, false, false, false, false, true},
		{0, 100, 0, false, false, false, false, false},
	}
	for _, tt := range tests {
		tx := &Transaction{ResponseCode: tt.code, Amount: tt.amount, AmountAuthorized: tt.authorized}
		equal(t, "Approved", tx.Approved(), tt.approved)
		equal(t, "PartiallyApproved", tx.PartiallyApproved(), tt.partial)
		equal(t, "Declined", tx.Declined(), tt.declined)
		equal(t, "GatewayDeclined", tx.GatewayDeclined(), tt.gateway)
		equal(t, "ProcessorError", tx.ProcessorError(), tt.proc)
	}
}

func TestTransactions_TerminalAndACHResponses(t *testing.T) {
	g, c := newGateway(t)
	g.respond("POST", "/api/transaction", 200, terminalSaleFixture)
	tx, err := c.Transactions.Sale(ctx(), &TransactionRequest{Amount: 500, PaymentMethod: PaymentMethod{Terminal: &TerminalPayment{ID: "t"}}})
	mustNoError(t, err)
	if tx.ResponseBody.Terminal == nil {
		t.Fatal("terminal response missing")
	}
	equal(t, "entry type", tx.ResponseBody.Terminal.EntryType, "swiped")
	equal(t, "cardholder", tx.ResponseBody.Terminal.CardholderName, "FDCS TEST CARD /MASTERCARD")
	var ps map[string]string
	mustNoError(t, json.Unmarshal(tx.ResponseBody.Terminal.ProcessorSpecific, &ps))
	equal(t, "BatchNum", ps["BatchNum"], "8")

	g.respond("POST", "/api/transaction", 200, achSaleFixture)
	tx, err = c.Transactions.Sale(ctx(), &TransactionRequest{Amount: 8912})
	mustNoError(t, err)
	if tx.ResponseBody.ACH == nil {
		t.Fatal("ach response missing")
	}
	equal(t, "sec code", tx.ResponseBody.ACH.SecCode, "web")
	equal(t, "response code", tx.ResponseBody.ACH.ResponseCode, 100)
	equal(t, "payment type", tx.PaymentType, "ach")
}

func TestTransactions_Get(t *testing.T) {
	g, c := newGateway(t)
	g.respond("GET", "/api/transaction/b7kgflt1tlv51er0fts0", 200, getTransactionFixture)
	tx, err := c.Transactions.Get(ctx(), "b7kgflt1tlv51er0fts0")
	mustNoError(t, err)
	assertRequest(t, g, "GET", "/api/transaction/b7kgflt1tlv51er0fts0")
	equal(t, "ID", tx.ID, "b7kgflt1tlv51er0fts0")
	equal(t, "Amount", tx.Amount, 1112)

	// An object under data works too.
	g.respond("GET", "/api/transaction/obj", 200, cardSaleFixture)
	tx, err = c.Transactions.Get(ctx(), "obj")
	mustNoError(t, err)
	equal(t, "ID", tx.ID, "d3bfel2cunifubaumn00")

	_, err = c.Transactions.Get(ctx(), "")
	mustError(t, err, "transaction id is required")
	equal(t, "no request for blank id", g.count(), 2)
}

func TestTransactions_Search(t *testing.T) {
	g, c := newGateway(t)
	g.respond("POST", "/api/transaction/search", 200, searchTransactionsFixture)

	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	list, err := c.Transactions.Search(ctx(), &TransactionSearchRequest{
		Status:         Equal(StatusPendingSettlement),
		Amount:         GreaterThan(1000),
		CreatedAt:      DateRange(start, start.Add(48*time.Hour)),
		BillingAddress: &AddressQuery{PostalCode: NotEqual("00000")},
		Limit:          50,
		Offset:         100,
	})
	mustNoError(t, err)
	assertRequest(t, g, "POST", "/api/transaction/search")
	body := g.bodyJSON()
	equal(t, "keys", joinKeys(body), "amount,billing_address,created_at,limit,offset,status")
	equal(t, "status op", body["status"].(map[string]any)["operator"], "=")
	equal(t, "amount op", body["amount"].(map[string]any)["operator"], ">")
	equal(t, "amount value", body["amount"].(map[string]any)["value"], float64(1000))
	equal(t, "start", body["created_at"].(map[string]any)["start_date"], "2024-01-01T00:00:00Z")
	equal(t, "end", body["created_at"].(map[string]any)["end_date"], "2024-01-03T00:00:00Z")
	equal(t, "postal op", body["billing_address"].(map[string]any)["postal_code"].(map[string]any)["operator"], "!=")

	equal(t, "TotalCount", list.TotalCount, 2)
	equal(t, "len", len(list.Data), 2)
	equal(t, "first id", list.Data[0].ID, "b84vgb2j8m0jujadi4v0")
	equal(t, "second source", list.Data[1].TransactionSource, SourceRecurring)
	equal(t, "second subscription", list.Data[1].SubscriptionID, "sub_1")
	if list.Data[0].CapturedAt != nil {
		t.Error("null captured_at should be nil")
	}
	equal(t, "item correlation", list.Data[0].LastResponse.CorrelationID, testCorrelationID)

	// A nil request sends an empty object.
	_, err = c.Transactions.Search(ctx(), nil)
	mustNoError(t, err)
	equal(t, "empty body", string(g.last().Body), "{}")
}

func TestTransactions_Capture(t *testing.T) {
	g, c := newGateway(t)
	g.respond("POST", "/api/transaction/auth1/capture", 200, cardSaleFixture)

	tx, err := c.Transactions.Capture(ctx(), "auth1", &CaptureRequest{Amount: 1299, OrderID: "o1"})
	mustNoError(t, err)
	assertRequest(t, g, "POST", "/api/transaction/auth1/capture")
	equal(t, "body", string(g.last().Body), `{"amount":1299,"order_id":"o1"}`)
	equal(t, "captured", tx.AmountCaptured, 1299)

	_, err = c.Transactions.Capture(ctx(), "auth1", nil)
	mustNoError(t, err)
	equal(t, "nil request captures in full", string(g.last().Body), "{}")

	_, err = c.Transactions.Capture(ctx(), " ", nil)
	mustError(t, err, "transaction id is required")
}

func TestTransactions_Void(t *testing.T) {
	g, c := newGateway(t)
	g.respond("POST", "/api/transaction/tx1/void", 200, `{"status":"success","msg":"success","data":null}`)
	mustNoError(t, c.Transactions.Void(ctx(), "tx1"))
	r := assertRequest(t, g, "POST", "/api/transaction/tx1/void")
	equal(t, "no body", len(r.Body), 0)

	g.respond("POST", "/api/transaction/tx2/void", 200, `{"status":"failed","msg":"transaction already settled"}`)
	mustError(t, c.Transactions.Void(ctx(), "tx2"), "already settled")
	mustError(t, c.Transactions.Void(ctx(), ""), "transaction id is required")
}

func TestTransactions_Refund(t *testing.T) {
	g, c := newGateway(t)
	refund := `{"status":"success","msg":"success","data":{"id":"ref1","type":"refund","status":"pending_settlement","response":"approved","response_code":100,"amount":500,"referenced_transaction_id":"tx1","created_at":"2025-09-26T20:30:09Z","updated_at":"2025-09-26T20:30:09Z","captured_at":"2025-09-26T20:30:09Z","settled_at":null}}`
	g.respond("POST", "/api/transaction/tx1/refund", 200, refund)

	tx, err := c.Transactions.Refund(ctx(), "tx1", &RefundRequest{Amount: 500})
	mustNoError(t, err)
	assertRequest(t, g, "POST", "/api/transaction/tx1/refund")
	equal(t, "body", string(g.last().Body), `{"amount":500}`)
	equal(t, "type", tx.Type, "refund")
	equal(t, "referenced", tx.ReferencedTransactionID, "tx1")

	_, err = c.Transactions.Refund(ctx(), "tx1", nil)
	mustNoError(t, err)
	equal(t, "nil request refunds in full", string(g.last().Body), "{}")

	_, err = c.Transactions.Refund(ctx(), "", nil)
	mustError(t, err, "transaction id is required")
}

func TestTransactions_ServerRejectsRequest(t *testing.T) {
	g, c := newGateway(t)
	g.handle("POST", "/api/transaction", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusBadRequest, `{"status":"failed","msg":"bad request error: invalid Postal Code"}`)
	})
	_, err := c.Transactions.Sale(ctx(), &TransactionRequest{Amount: 1})
	apiErr, ok := AsError(err)
	if !ok || !apiErr.IsBadRequest() {
		t.Fatalf("expected bad request *Error, got %v", err)
	}
	equal(t, "msg", apiErr.Msg, "bad request error: invalid Postal Code")
}
