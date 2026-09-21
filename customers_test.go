package fluidpay

import (
	"encoding/json"
	"net/url"
	"testing"
)

func TestCustomers_Create(t *testing.T) {
	g, c := newGateway(t)
	g.respond("POST", "/api/vault/customer", 200, vaultCustomerFixture)

	req := &CustomerCreateRequest{
		VerificationOptions: VerificationOptions{Validate: true, BypassRuleEngine: true},
		Description:         "test description",
		Flags:               []string{CustomerFlagSurchargeExempt},
		DefaultPayment: &VaultPaymentMethod{Card: &VaultCard{
			Number: "4111111111111111", ExpirationDate: "12/30",
		}},
		DefaultBillingAddress: &VaultAddress{FirstName: "John", Line1: "123 Some St", PostalCode: "60187"},
	}
	cust, err := c.Customers.Create(ctx(), req)
	mustNoError(t, err)

	r := assertRequest(t, g, "POST", "/api/vault/customer")
	q, _ := url.ParseQuery(r.Query)
	equal(t, "validate", q.Get("validate"), "true")
	equal(t, "bypass", q.Get("bypass_rule_engine"), "true")
	equal(t, "authorize absent", q.Has("authorize"), false)
	body := g.bodyJSON()
	equal(t, "keys", joinKeys(body), "default_billing_address,default_payment,description,flags")
	card := body["default_payment"].(map[string]any)["card"].(map[string]any)
	equal(t, "card keys", joinKeys(card), "expiration_date,number")
	equal(t, "line_1", body["default_billing_address"].(map[string]any)["line_1"], "123 Some St")

	equal(t, "ID", cust.ID, "952f250d-fa85-40ac-a45e-3886f98032c6")
	equal(t, "OwnerID", cust.OwnerID, "testmerchant12345678")
	equal(t, "Description", cust.Description, "test description")
	equal(t, "Notes", cust.Notes, "vip")
	equal(t, "Flags", cust.Flags[0], "surcharge_exempt")
	equal(t, "Addresses", len(cust.Addresses), 1)
	equal(t, "Address line", cust.Addresses[0].Line1, "123 Some St")
	equal(t, "Address hash", cust.Addresses[0].Hash, "cV9qS6UvnnRDJCkgfitHqGrKSRfTARDsLy5zfkNNTcM=")
	equal(t, "Defaults", cust.Defaults.PaymentMethodID, "bl6nqk69ku6897l45fq0")
	equal(t, "Cards", cust.Cards[0].MaskedNumber, "411111******1111")
	equal(t, "ACH", cust.ACH[0].MaskedAccountNumber, "XXXXX1111")
	equal(t, "CreatedAt", cust.CreatedAt.IsZero(), false)
	equal(t, "correlation", cust.LastResponse.CorrelationID, testCorrelationID)

	// Helpers.
	equal(t, "Address()", cust.Address("bl6nqk69ku6897l45fp0").City, "Some Town")
	if cust.Address("nope") != nil || cust.Card("nope") != nil {
		t.Error("lookups for unknown ids should be nil")
	}
	equal(t, "DefaultCard", cust.DefaultCard().ID, "bl6nqk69ku6897l45fq0")
	cust.Defaults.PaymentMethodType = PaymentMethodTypeACH
	if cust.DefaultCard() != nil {
		t.Error("DefaultCard should be nil for ach default")
	}

	// Nil request creates an empty customer.
	_, err = c.Customers.Create(ctx(), nil)
	mustNoError(t, err)
	equal(t, "empty body", string(g.last().Body), "{}")
	equal(t, "no query", g.last().Query, "")
}

func TestCustomer_MarshalRoundTrip(t *testing.T) {
	var cust Customer
	mustNoError(t, json.Unmarshal([]byte(`{"id":"c1","data":{"customer":{"description":"d","payments":{"cards":[{"id":"k"}]}}}}`), &cust))
	equal(t, "id", cust.ID, "c1")
	equal(t, "card", cust.Cards[0].ID, "k")
	// Marshaling a Customer produces the flattened shape.
	b, err := json.Marshal(cust)
	mustNoError(t, err)
	var m map[string]any
	mustNoError(t, json.Unmarshal(b, &m))
	equal(t, "flattened", m["description"], "d")
}

func TestCustomers_Get(t *testing.T) {
	g, c := newGateway(t)
	g.respond("GET", "/api/vault/cust1", 200, vaultCustomerFixture)
	cust, err := c.Customers.Get(ctx(), "cust1")
	mustNoError(t, err)
	assertRequest(t, g, "GET", "/api/vault/cust1")
	equal(t, "ID", cust.ID, "952f250d-fa85-40ac-a45e-3886f98032c6")

	_, err = c.Customers.Get(ctx(), "")
	mustError(t, err, "customer id is required")
}

func TestCustomers_Search(t *testing.T) {
	g, c := newGateway(t)
	g.respond("POST", "/api/vault/customer/search", 200, vaultSearchFixture)
	list, err := c.Customers.Search(ctx(), &CustomerSearchRequest{
		Email: Equal("the.original.e@gmail.com"),
		Limit: 25,
	})
	mustNoError(t, err)
	assertRequest(t, g, "POST", "/api/vault/customer/search")
	equal(t, "body", string(g.last().Body), `{"email":{"operator":"=","value":"the.original.e@gmail.com"},"limit":25}`)
	equal(t, "TotalCount", list.TotalCount, 81)
	equal(t, "len", len(list.Data), 1)
	equal(t, "ID", list.Data[0].ID, "bgl4vq1erttokpdi1kj0")
	equal(t, "cards", len(list.Data[0].Cards), 2)
	equal(t, "second card", list.Data[0].Cards[1].MaskedNumber, "400551******0004")
	equal(t, "address", list.Data[0].Addresses[0].Line1, "ABC LANE")

	_, err = c.Customers.Search(ctx(), nil)
	mustNoError(t, err)
	equal(t, "nil body", string(g.last().Body), "{}")
}

func TestCustomers_Update(t *testing.T) {
	g, c := newGateway(t)
	g.respond("POST", "/api/vault/customer/cust1", 200, vaultCustomerFixture)
	_, err := c.Customers.Update(ctx(), "cust1", &CustomerUpdateRequest{
		Notes:    "vip",
		Defaults: &CustomerDefaults{PaymentMethodID: "pm", PaymentMethodType: PaymentMethodTypeCard},
	})
	mustNoError(t, err)
	assertRequest(t, g, "POST", "/api/vault/customer/cust1")
	equal(t, "body", string(g.last().Body), `{"notes":"vip","defaults":{"payment_method_id":"pm","payment_method_type":"card"}}`)

	_, err = c.Customers.Update(ctx(), "cust1", nil)
	mustError(t, err, "must not be nil")
	_, err = c.Customers.Update(ctx(), "", &CustomerUpdateRequest{})
	mustError(t, err, "customer id is required")
}

func TestCustomers_Delete(t *testing.T) {
	g, c := newGateway(t)
	g.respond("DELETE", "/api/vault/cust1", 200, `{"status":"success","msg":"success"}`)
	mustNoError(t, c.Customers.Delete(ctx(), "cust1"))
	assertRequest(t, g, "DELETE", "/api/vault/cust1")
	mustError(t, c.Customers.Delete(ctx(), ""), "customer id is required")
}

func TestCustomers_Addresses(t *testing.T) {
	g, c := newGateway(t)
	g.respond("POST", "/api/vault/customer/cust1/address", 200, vaultAddressCreatedFixture)
	g.respond("POST", "/api/vault/customer/cust1/address/addr1", 200, vaultCustomerFixture)
	g.respond("DELETE", "/api/vault/customer/cust1/address/addr1", 200, `{"status":"success","msg":"success"}`)

	cust, err := c.Customers.CreateAddress(ctx(), "cust1", &VaultAddress{Line1: "1 Main", City: "Chicago"})
	mustNoError(t, err)
	assertRequest(t, g, "POST", "/api/vault/customer/cust1/address")
	equal(t, "body", string(g.last().Body), `{"line_1":"1 Main","city":"Chicago"}`)
	equal(t, "CreatedAddressID", cust.CreatedAddressID, "newaddr0000000000000")
	equal(t, "address present", cust.Address("newaddr0000000000000").Line1, "1 Main")

	_, err = c.Customers.UpdateAddress(ctx(), "cust1", "addr1", &VaultAddress{City: "Wheaton"})
	mustNoError(t, err)
	assertRequest(t, g, "POST", "/api/vault/customer/cust1/address/addr1")

	mustNoError(t, c.Customers.DeleteAddress(ctx(), "cust1", "addr1"))
	assertRequest(t, g, "DELETE", "/api/vault/customer/cust1/address/addr1")

	_, err = c.Customers.CreateAddress(ctx(), "cust1", nil)
	mustError(t, err, "must not be nil")
	_, err = c.Customers.UpdateAddress(ctx(), "cust1", "", &VaultAddress{})
	mustError(t, err, "address id is required")
	mustError(t, c.Customers.DeleteAddress(ctx(), "", "addr1"), "customer id is required")
}

func TestCustomers_PaymentMethods(t *testing.T) {
	g, c := newGateway(t)
	ok := `{"status":"success","msg":"success"}`
	for _, p := range []string{"card", "ach", "token", "applepay", "googlepay", "card/pm1", "ach/pm1", "token/pm1"} {
		g.respond("POST", "/api/vault/customer/cust1/"+p, 200, vaultPaymentCreatedFixture)
	}
	g.respond("DELETE", "/api/vault/customer/cust1/card/pm1", 200, ok)
	g.respond("DELETE", "/api/vault/customer/cust1/ach/pm1", 200, ok)

	cust, err := c.Customers.CreateCard(ctx(), "cust1", &VaultCard{
		VerificationOptions: VerificationOptions{Authorize: true},
		Number:              "4111111111111111", ExpirationDate: "12/30",
	})
	mustNoError(t, err)
	r := assertRequest(t, g, "POST", "/api/vault/customer/cust1/card")
	equal(t, "query", r.Query, "authorize=true")
	equal(t, "body", string(r.Body), `{"number":"4111111111111111","expiration_date":"12/30"}`)
	equal(t, "CreatedPaymentMethodID", cust.CreatedPaymentMethodID, "newcard0000000000000")
	equal(t, "card present", cust.Card("newcard0000000000000").MaskedNumber, "411111******1111")

	_, err = c.Customers.CreateACH(ctx(), "cust1", &VaultACH{AccountNumber: "111111111", RoutingNumber: "111111111", AccountType: AccountTypeChecking, SecCode: SecCodeWeb})
	mustNoError(t, err)
	r = assertRequest(t, g, "POST", "/api/vault/customer/cust1/ach")
	equal(t, "ach body", string(r.Body), `{"account_number":"111111111","routing_number":"111111111","account_type":"checking","sec_code":"web"}`)
	equal(t, "no query", r.Query, "")

	_, err = c.Customers.CreateToken(ctx(), "cust1", &VaultToken{VerificationOptions: VerificationOptions{Validate: true}, Token: "tok_1"})
	mustNoError(t, err)
	r = assertRequest(t, g, "POST", "/api/vault/customer/cust1/token")
	equal(t, "token body", string(r.Body), `{"token":"tok_1"}`)
	equal(t, "token query", r.Query, "validate=true")

	_, err = c.Customers.CreateApplePay(ctx(), "cust1", &VaultApplePay{KeyID: "k", TemporaryToken: "tmp"})
	mustNoError(t, err)
	r = assertRequest(t, g, "POST", "/api/vault/customer/cust1/applepay")
	equal(t, "apple body", string(r.Body), `{"key_id":"k","temporary_token":"tmp"}`)

	_, err = c.Customers.CreateGooglePay(ctx(), "cust1", json.RawMessage(`{"signature":"s","protocolVersion":"ECv2"}`), VerificationOptions{Validate: true})
	mustNoError(t, err)
	r = assertRequest(t, g, "POST", "/api/vault/customer/cust1/googlepay")
	equal(t, "google body passthrough", string(r.Body), `{"signature":"s","protocolVersion":"ECv2"}`)
	equal(t, "google query", r.Query, "validate=true")

	_, err = c.Customers.UpdateCard(ctx(), "cust1", "pm1", &VaultCard{Number: "4", ExpirationDate: "1/1", Flags: []string{"surcharge_exempt"}})
	mustNoError(t, err)
	r = assertRequest(t, g, "POST", "/api/vault/customer/cust1/card/pm1")
	equal(t, "update card body", string(r.Body), `{"number":"4","expiration_date":"1/1","flags":["surcharge_exempt"]}`)

	_, err = c.Customers.UpdateACH(ctx(), "cust1", "pm1", &VaultACH{AccountNumber: "1", RoutingNumber: "2", AccountType: "savings", SecCode: "ppd"})
	mustNoError(t, err)
	assertRequest(t, g, "POST", "/api/vault/customer/cust1/ach/pm1")

	_, err = c.Customers.UpdateToken(ctx(), "cust1", "pm1", &VaultToken{Token: "tok_2"})
	mustNoError(t, err)
	assertRequest(t, g, "POST", "/api/vault/customer/cust1/token/pm1")

	mustNoError(t, c.Customers.DeleteCard(ctx(), "cust1", "pm1"))
	assertRequest(t, g, "DELETE", "/api/vault/customer/cust1/card/pm1")
	mustNoError(t, c.Customers.DeleteACH(ctx(), "cust1", "pm1"))
	assertRequest(t, g, "DELETE", "/api/vault/customer/cust1/ach/pm1")
}

func TestCustomers_Validation(t *testing.T) {
	_, c := newGateway(t)
	var err error
	_, err = c.Customers.CreateCard(ctx(), "", &VaultCard{})
	mustError(t, err, "customer id is required")
	_, err = c.Customers.CreateCard(ctx(), "c", nil)
	mustError(t, err, "must not be nil")
	_, err = c.Customers.CreateACH(ctx(), "c", nil)
	mustError(t, err, "must not be nil")
	_, err = c.Customers.CreateToken(ctx(), "c", nil)
	mustError(t, err, "must not be nil")
	_, err = c.Customers.CreateApplePay(ctx(), "c", nil)
	mustError(t, err, "must not be nil")
	_, err = c.Customers.CreateGooglePay(ctx(), "c", nil, VerificationOptions{})
	mustError(t, err, "must not be nil")
	_, err = c.Customers.UpdateCard(ctx(), "c", "", &VaultCard{})
	mustError(t, err, "payment method id is required")
	_, err = c.Customers.UpdateCard(ctx(), "c", "p", nil)
	mustError(t, err, "must not be nil")
	_, err = c.Customers.UpdateACH(ctx(), "c", "p", nil)
	mustError(t, err, "must not be nil")
	_, err = c.Customers.UpdateToken(ctx(), "c", "p", nil)
	mustError(t, err, "must not be nil")
	mustError(t, c.Customers.DeleteCard(ctx(), "c", ""), "card id is required")
	mustError(t, c.Customers.DeleteACH(ctx(), "c", ""), "ach id is required")
}

func TestVerificationOptions_Query(t *testing.T) {
	if q := (VerificationOptions{}).query(); q != nil {
		t.Error("empty options should produce nil query")
	}
	q := VerificationOptions{Validate: true, Authorize: true, BypassRuleEngine: true}.query()
	equal(t, "encoded", q.Encode(), "authorize=true&bypass_rule_engine=true&validate=true")
}
