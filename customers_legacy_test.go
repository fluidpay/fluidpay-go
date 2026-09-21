package fluidpay

import "testing"

func TestLegacyCustomers_CRUD(t *testing.T) {
	g, c := newGateway(t)
	g.respond("POST", "/api/customer", 200, legacyCustomerFixture)
	g.respond("GET", "/api/customer/cust1", 200, legacyCustomerFixture)
	g.respond("POST", "/api/customer/search", 200, `{"status":"success","msg":"success","data":[`+legacyCustomerData+`],"count":1}`)
	g.respond("POST", "/api/customer/cust1", 200, `{"status":"success","msg":"success","data":null}`)
	g.respond("DELETE", "/api/customer/cust1", 200, `{"status":"success","msg":"success"}`)

	cust, err := c.LegacyCustomers.Create(ctx(), &LegacyCustomerRequest{
		Validate:    true,
		Description: "test description",
		PaymentMethod: &LegacyPaymentMethod{Card: &LegacyCardInput{
			CardNumber: "4111111111111111", ExpirationDate: "12/20",
		}},
		BillingAddress: &Address{FirstName: "John", AddressLine1: "123 Some St"},
	})
	mustNoError(t, err)
	r := assertRequest(t, g, "POST", "/api/customer")
	equal(t, "validate query", r.Query, "validate=true")
	equal(t, "body", string(r.Body), `{"description":"test description","payment_method":{"card":{"card_number":"4111111111111111","expiration_date":"12/20"}},"billing_address":{"first_name":"John","address_line_1":"123 Some St"}}`)
	equal(t, "ID", cust.ID, "b798ls2q9qq646ksu070")
	equal(t, "card", cust.PaymentMethod.Card.MaskedCard, "411111******1111")
	equal(t, "billing", cust.BillingAddress.AddressLine1, "123 Some St")
	equal(t, "billing id", cust.BillingAddress.ID, "b798ls2q9qq646ksu07g")
	if cust.ShippingAddress != nil {
		t.Error("null shipping address should be nil")
	}

	cust, err = c.LegacyCustomers.Get(ctx(), "cust1")
	mustNoError(t, err)
	assertRequest(t, g, "GET", "/api/customer/cust1")
	equal(t, "ID", cust.ID, "b798ls2q9qq646ksu070")

	list, err := c.LegacyCustomers.Search(ctx(), &LegacyCustomerSearchRequest{ID: Equal("b798ls2q9qq646ksu070"), Limit: 10})
	mustNoError(t, err)
	assertRequest(t, g, "POST", "/api/customer/search")
	equal(t, "body", string(g.last().Body), `{"id":{"operator":"=","value":"b798ls2q9qq646ksu070"},"limit":10}`)
	equal(t, "count fallback", list.TotalCount, 1)
	equal(t, "len", len(list.Data), 1)

	mustNoError(t, c.LegacyCustomers.Update(ctx(), "cust1", &LegacyCustomerUpdateRequest{PaymentMethod: PaymentMethodTypeCard, PaymentMethodID: "pm"}))
	assertRequest(t, g, "POST", "/api/customer/cust1")
	equal(t, "update body", string(g.last().Body), `{"payment_method":"card","payment_method_id":"pm"}`)

	mustNoError(t, c.LegacyCustomers.Delete(ctx(), "cust1"))
	assertRequest(t, g, "DELETE", "/api/customer/cust1")

	_, err = c.LegacyCustomers.Create(ctx(), nil)
	mustError(t, err, "must not be nil")
	_, err = c.LegacyCustomers.Get(ctx(), "")
	mustError(t, err, "customer id is required")
	mustError(t, c.LegacyCustomers.Update(ctx(), "cust1", nil), "must not be nil")
	mustError(t, c.LegacyCustomers.Delete(ctx(), ""), "customer id is required")
	_, err = c.LegacyCustomers.Search(ctx(), nil)
	mustNoError(t, err)
}

// legacyCustomerData is the data member of legacyCustomerFixture.
const legacyCustomerData = `{"id":"b798ls2q9qq646ksu070","description":"test description","payment_method":{"card":{"id":"b798ls2q9qq646ksu080","masked_card":"411111******1111"}},"created_at":"2017-10-02T18:52:32Z","updated_at":"2017-10-02T18:52:32Z"}`

func TestLegacyCustomers_Addresses(t *testing.T) {
	g, c := newGateway(t)
	g.respond("POST", "/api/customer/cust1/address", 200, legacyAddressListFixture)
	g.respond("GET", "/api/customer/cust1/address/addr1", 200, legacyAddressListFixture)
	g.respond("GET", "/api/customer/cust1/addresses", 200, legacyAddressListFixture)
	g.respond("POST", "/api/customer/cust1/address/addr1", 200, legacyAddressListFixture)
	g.respond("DELETE", "/api/customer/cust1/address/addr1", 200, `{"status":"success","msg":"success"}`)

	addr, err := c.LegacyCustomers.CreateAddress(ctx(), "cust1", &Address{FirstName: "John", City: "Some Town"})
	mustNoError(t, err)
	assertRequest(t, g, "POST", "/api/customer/cust1/address")
	equal(t, "body", string(g.last().Body), `{"first_name":"John","city":"Some Town"}`)
	equal(t, "ID", addr.ID, "b798ls2q9qq646ksu07g")
	equal(t, "CustomerID", addr.CustomerID, "b798ls2q9qq646ksu070")
	equal(t, "line", addr.AddressLine1, "123 Some St")

	addr, err = c.LegacyCustomers.GetAddress(ctx(), "cust1", "addr1")
	mustNoError(t, err)
	assertRequest(t, g, "GET", "/api/customer/cust1/address/addr1")
	equal(t, "array-wrapped single", addr.ID, "b798ls2q9qq646ksu07g")

	list, err := c.LegacyCustomers.ListAddresses(ctx(), "cust1")
	mustNoError(t, err)
	assertRequest(t, g, "GET", "/api/customer/cust1/addresses")
	equal(t, "count", list.TotalCount, 1)
	equal(t, "len", len(list.Data), 1)

	_, err = c.LegacyCustomers.UpdateAddress(ctx(), "cust1", "addr1", &Address{City: "Wheaton"})
	mustNoError(t, err)
	assertRequest(t, g, "POST", "/api/customer/cust1/address/addr1")

	mustNoError(t, c.LegacyCustomers.DeleteAddress(ctx(), "cust1", "addr1"))
	assertRequest(t, g, "DELETE", "/api/customer/cust1/address/addr1")

	_, err = c.LegacyCustomers.CreateAddress(ctx(), "cust1", nil)
	mustError(t, err, "must not be nil")
	_, err = c.LegacyCustomers.GetAddress(ctx(), "cust1", "")
	mustError(t, err, "address id is required")
	_, err = c.LegacyCustomers.ListAddresses(ctx(), "")
	mustError(t, err, "customer id is required")
	_, err = c.LegacyCustomers.UpdateAddress(ctx(), "cust1", "addr1", nil)
	mustError(t, err, "must not be nil")
	mustError(t, c.LegacyCustomers.DeleteAddress(ctx(), "cust1", ""), "address id is required")
}

func TestLegacyCustomers_PaymentMethods(t *testing.T) {
	g, c := newGateway(t)
	ok := `{"status":"success","msg":"success","data":null}`
	g.respond("POST", "/api/customer/cust1/paymentmethod/card", 200, legacyCardCreatedFixture)
	g.respond("POST", "/api/customer/cust1/paymentmethod/ach", 200, `{"status":"success","msg":"success","data":{"ach":{"id":"ach1","sec_code":"web","masked_account_number":"XXXXX1111"}}}`)
	g.respond("POST", "/api/customer/cust1/paymentmethod/token", 200, legacyCardCreatedFixture)
	g.respond("GET", "/api/customer/cust1/paymentmethod/card/card1", 200, legacyCardListFixture)
	g.respond("GET", "/api/customer/cust1/paymentmethod/ach/ach1", 200, `{"status":"success","msg":"success","data":[{"id":"ach1","sec_code":"web"}],"count":1}`)
	g.respond("GET", "/api/customer/cust1/paymentmethod/card", 200, legacyCardListFixture)
	g.respond("GET", "/api/customer/cust1/paymentmethod/ach", 200, `{"status":"success","msg":"success","data":[{"id":"ach1"}],"count":1}`)
	g.respond("POST", "/api/customer/cust1/paymentmethod/card/card1", 200, ok)
	g.respond("POST", "/api/customer/cust1/paymentmethod/ach/ach1", 200, ok)
	g.respond("DELETE", "/api/customer/cust1/paymentmethod/card/card1", 200, ok)
	g.respond("DELETE", "/api/customer/cust1/paymentmethod/ach/ach1", 200, ok)

	pm, err := c.LegacyCustomers.CreateCard(ctx(), "cust1", &LegacyCardTokenRequest{Validate: true, Number: "4111111111111112", ExpirationDate: "12/20"})
	mustNoError(t, err)
	r := assertRequest(t, g, "POST", "/api/customer/cust1/paymentmethod/card")
	equal(t, "query", r.Query, "validate=true")
	equal(t, "body uses number", string(r.Body), `{"number":"4111111111111112","expiration_date":"12/20"}`)
	equal(t, "card id", pm.Card.ID, "b799g1iq9qq6dk5l39i0")
	equal(t, "last four", pm.Card.LastFour, "1112")

	pm, err = c.LegacyCustomers.CreateACH(ctx(), "cust1", &LegacyACHTokenRequest{AccountNumber: "1", RoutingNumber: "2", AccountType: "checking", SecCode: "web"})
	mustNoError(t, err)
	assertRequest(t, g, "POST", "/api/customer/cust1/paymentmethod/ach")
	equal(t, "ach id", pm.ACH.ID, "ach1")

	_, err = c.LegacyCustomers.CreateToken(ctx(), "cust1", &LegacyTokenRequest{Token: "tok"})
	mustNoError(t, err)
	r = assertRequest(t, g, "POST", "/api/customer/cust1/paymentmethod/token")
	equal(t, "token body", string(r.Body), `{"token":"tok"}`)

	card, err := c.LegacyCustomers.GetCard(ctx(), "cust1", "card1")
	mustNoError(t, err)
	assertRequest(t, g, "GET", "/api/customer/cust1/paymentmethod/card/card1")
	equal(t, "first of array", card.ID, "b798ls2q9qq646ksu080")

	ach, err := c.LegacyCustomers.GetACH(ctx(), "cust1", "ach1")
	mustNoError(t, err)
	equal(t, "ach", ach.SecCode, "web")

	cards, err := c.LegacyCustomers.ListCards(ctx(), "cust1")
	mustNoError(t, err)
	assertRequest(t, g, "GET", "/api/customer/cust1/paymentmethod/card")
	equal(t, "cards count", cards.TotalCount, 2)
	equal(t, "cards len", len(cards.Data), 2)

	achs, err := c.LegacyCustomers.ListACH(ctx(), "cust1")
	mustNoError(t, err)
	equal(t, "ach len", len(achs.Data), 1)

	mustNoError(t, c.LegacyCustomers.UpdateCard(ctx(), "cust1", "card1", &LegacyCardInput{CardNumber: "4", ExpirationDate: "1/1"}))
	r = assertRequest(t, g, "POST", "/api/customer/cust1/paymentmethod/card/card1")
	equal(t, "nested card", string(r.Body), `{"card":{"card_number":"4","expiration_date":"1/1"}}`)

	mustNoError(t, c.LegacyCustomers.UpdateACH(ctx(), "cust1", "ach1", &LegacyACHInput{AccountNumber: "1", RoutingNumber: "2", AccountType: "checking", SecCode: "web"}))
	r = assertRequest(t, g, "POST", "/api/customer/cust1/paymentmethod/ach/ach1")
	equal(t, "nested ach", string(r.Body), `{"ach":{"account_number":"1","routing_number":"2","account_type":"checking","sec_code":"web"}}`)

	mustNoError(t, c.LegacyCustomers.DeleteCard(ctx(), "cust1", "card1"))
	assertRequest(t, g, "DELETE", "/api/customer/cust1/paymentmethod/card/card1")
	mustNoError(t, c.LegacyCustomers.DeleteACH(ctx(), "cust1", "ach1"))
	assertRequest(t, g, "DELETE", "/api/customer/cust1/paymentmethod/ach/ach1")

	// Validation.
	_, err = c.LegacyCustomers.CreateCard(ctx(), "cust1", nil)
	mustError(t, err, "must not be nil")
	_, err = c.LegacyCustomers.CreateACH(ctx(), "", &LegacyACHTokenRequest{})
	mustError(t, err, "customer id is required")
	_, err = c.LegacyCustomers.CreateToken(ctx(), "cust1", nil)
	mustError(t, err, "must not be nil")
	_, err = c.LegacyCustomers.GetCard(ctx(), "cust1", "")
	mustError(t, err, "card id is required")
	_, err = c.LegacyCustomers.GetACH(ctx(), "", "a")
	mustError(t, err, "customer id is required")
	_, err = c.LegacyCustomers.ListCards(ctx(), "")
	mustError(t, err, "customer id is required")
	_, err = c.LegacyCustomers.ListACH(ctx(), "")
	mustError(t, err, "customer id is required")
	mustError(t, c.LegacyCustomers.UpdateCard(ctx(), "cust1", "card1", nil), "must not be nil")
	mustError(t, c.LegacyCustomers.UpdateACH(ctx(), "cust1", "", &LegacyACHInput{}), "ach id is required")
	mustError(t, c.LegacyCustomers.DeleteCard(ctx(), "", "card1"), "customer id is required")
	mustError(t, c.LegacyCustomers.DeleteACH(ctx(), "cust1", ""), "ach id is required")
}
