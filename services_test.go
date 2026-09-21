package fluidpay

import (
	"net/http"
	"testing"
	"time"
)

func TestTerminals(t *testing.T) {
	g, c := newGateway(t)
	g.respond("GET", "/api/terminals", 200, terminalsFixture)
	g.respond("POST", "/api/terminal/term1/settle", 200, `{"status":"success","msg":"success","data":null}`)

	list, err := c.Terminals.List(ctx())
	mustNoError(t, err)
	assertRequest(t, g, "GET", "/api/terminals")
	equal(t, "total", list.TotalCount, 1)
	term := list.Data[0]
	equal(t, "ID", term.ID, "1ucio551tlv85l7moe5s")
	equal(t, "TPN", term.TPN, "1811000XXXX")
	equal(t, "Active", term.Active(), true)
	equal(t, "AutoSettle", term.AutoSettle, true)
	equal(t, "zero updated_at", term.UpdatedAt.IsZero(), true)

	mustNoError(t, c.Terminals.Settle(ctx(), "term1"))
	assertRequest(t, g, "POST", "/api/terminal/term1/settle")
	mustError(t, c.Terminals.Settle(ctx(), ""), "terminal id is required")
}

func TestSettlements_SearchBatches(t *testing.T) {
	g, c := newGateway(t)
	g.respond("POST", "/api/settlement/batch/search", 200, settlementSearchFixture)

	day := time.Date(2018, 12, 4, 0, 0, 0, 0, time.UTC)
	res, err := c.Settlements.SearchBatches(ctx(), &SettlementBatchSearchRequest{BatchDate: Day(day), Limit: 10})
	mustNoError(t, err)
	assertRequest(t, g, "POST", "/api/settlement/batch/search")
	equal(t, "body", string(g.last().Body), `{"batch_date":{"start_date":"2018-12-04T00:00:00Z","end_date":"2018-12-04T23:59:59Z"},"limit":10}`)
	equal(t, "TotalCount", res.TotalCount, 1)
	equal(t, "summary", res.Summary[0].Captured, 10000)
	equal(t, "result id", res.Results[0].ID, "cccccccccccccccccccc")
	equal(t, "net deposit", res.Results[0].NetDeposit, 10000)
	equal(t, "batch date", res.Results[0].BatchDate.Format(time.RFC3339), "2018-12-04T13:31:02Z")

	_, err = c.Settlements.SearchBatches(ctx(), nil)
	mustNoError(t, err)
	equal(t, "nil body", string(g.last().Body), "{}")
}

func TestLookup_BIN(t *testing.T) {
	g, c := newGateway(t)
	// The documented response carries no envelope.
	g.respond("POST", "/api/lookup/bin/protected", 200, binLookupFixture)

	res, err := c.Lookup.BIN(ctx(), &BINLookupRequest{BIN: "424242", Country: "US", State: "IL"})
	mustNoError(t, err)
	assertRequest(t, g, "POST", "/api/lookup/bin/protected")
	equal(t, "body", string(g.last().Body), `{"bin":"424242","country":"US","state":"IL"}`)
	equal(t, "brand", res.CardBrand, "Visa")
	equal(t, "surchargeable", res.IsSurchargeable, true)
	equal(t, "correlation", res.LastResponse.CorrelationID, testCorrelationID)

	// An enveloped response is handled as well.
	g.respond("POST", "/api/lookup/bin/protected", 200, wrap(binLookupFixture))
	res, err = c.Lookup.BIN(ctx(), &BINLookupRequest{TempToken: "tok"})
	mustNoError(t, err)
	equal(t, "enveloped brand", res.CardBrand, "Visa")

	_, err = c.Lookup.BIN(ctx(), nil)
	mustError(t, err, "must not be nil")
	_, err = c.Lookup.BIN(ctx(), &BINLookupRequest{Country: "US"})
	mustError(t, err, "requires one of")
}

func TestUsers(t *testing.T) {
	g, c := newGateway(t)
	g.respond("GET", "/api/user", 200, userFixture)
	g.respond("GET", "/api/user/u1", 200, userFixture)
	g.respond("GET", "/api/users", 200, listEnvelope(`[{"id":"u1","username":"a"},{"id":"u2","username":"b"}]`, 2))
	g.respond("POST", "/api/user", 200, userFixture)
	g.respond("POST", "/api/user/u1", 200, userFixture)
	g.respond("DELETE", "/api/user/u1", 200, `{"status":"success","msg":"success","data":null}`)
	g.respond("POST", "/api/user/change-password", 200, `{"status":"success","msg":"success"}`)

	me, err := c.Users.Current(ctx())
	mustNoError(t, err)
	assertRequest(t, g, "GET", "/api/user")
	equal(t, "username", me.Username, "will_test")
	equal(t, "Can(process_sale)", me.Can("process_sale"), true)
	equal(t, "Can(process_credit)", me.Can("process_credit"), false)
	equal(t, "Can(unknown)", me.Can("nope"), false)
	equal(t, "flags", me.Flags["password_expired"], "true")
	equal(t, "api key", me.APIKey, "api_xxxxxxxxxxxxxxxxxxxx")
	equal(t, "raw defaults kept", string(me.Defaults), `{"processor_id": ""}`)

	_, err = c.Users.Get(ctx(), "u1")
	mustNoError(t, err)
	assertRequest(t, g, "GET", "/api/user/u1")

	list, err := c.Users.List(ctx())
	mustNoError(t, err)
	assertRequest(t, g, "GET", "/api/users")
	equal(t, "total", list.TotalCount, 2)

	_, err = c.Users.Create(ctx(), &UserCreateRequest{
		Username: "will_test", Name: "test user", Phone: "5555555555", Email: "user@fluidpay.com",
		Timezone: "America/Chicago", Password: "s3cret!!", Status: UserStatusActive, Role: UserRoleAdmin,
		CreateAPIKey: true, Permissions: map[string]bool{"process_sale": true},
	})
	mustNoError(t, err)
	assertRequest(t, g, "POST", "/api/user")
	equal(t, "body", string(g.last().Body), `{"username":"will_test","name":"test user","phone":"5555555555","email":"user@fluidpay.com","timezone":"America/Chicago","password":"s3cret!!","status":"active","role":"admin","create_api_key":true,"permissions":{"process_sale":true}}`)

	_, err = c.Users.Update(ctx(), "u1", &UserUpdateRequest{Status: UserStatusDisabled})
	mustNoError(t, err)
	assertRequest(t, g, "POST", "/api/user/u1")
	equal(t, "update body", string(g.last().Body), `{"status":"disabled"}`)

	mustNoError(t, c.Users.Delete(ctx(), "u1"))
	assertRequest(t, g, "DELETE", "/api/user/u1")

	mustNoError(t, c.Users.ChangePassword(ctx(), &ChangePasswordRequest{Username: "will_test", CurrentPassword: "a", NewPassword: "b"}))
	assertRequest(t, g, "POST", "/api/user/change-password")

	_, err = c.Users.Get(ctx(), "")
	mustError(t, err, "user id is required")
	_, err = c.Users.Create(ctx(), nil)
	mustError(t, err, "must not be nil")
	_, err = c.Users.Update(ctx(), "u1", nil)
	mustError(t, err, "must not be nil")
	_, err = c.Users.Update(ctx(), "", &UserUpdateRequest{})
	mustError(t, err, "user id is required")
	mustError(t, c.Users.Delete(ctx(), ""), "user id is required")
	mustError(t, c.Users.ChangePassword(ctx(), nil), "must not be nil")
}

func TestAPIKeys(t *testing.T) {
	g, c := newGateway(t)
	g.respond("POST", "/api/user/apikey", 200, apiKeyFixture)
	g.respond("GET", "/api/user/apikeys", 200, listEnvelope(`[{"id":"key1","type":"api","name":"backend"},{"id":"key2","type":"pub","name":"web"}]`, 2))
	g.respond("DELETE", "/api/user/apikey/key1", 200, `{"status":"success","msg":"success","data":null}`)

	key, err := c.APIKeys.Create(ctx(), &APIKeyCreateRequest{Type: APIKeyTypePrivate, Name: "backend"})
	mustNoError(t, err)
	assertRequest(t, g, "POST", "/api/user/apikey")
	equal(t, "body", string(g.last().Body), `{"type":"api","name":"backend"}`)
	equal(t, "Key", key.Key, "api_2EXPCBb5hEnyIa79UhE4Pa30WmH")
	equal(t, "ID", key.ID, "key00000000000000001")

	list, err := c.APIKeys.List(ctx())
	mustNoError(t, err)
	assertRequest(t, g, "GET", "/api/user/apikeys")
	equal(t, "total", list.TotalCount, 2)
	equal(t, "public key type", list.Data[1].Type, APIKeyTypePublic)
	equal(t, "secret not listed", list.Data[0].Key, "")

	mustNoError(t, c.APIKeys.Delete(ctx(), "key1"))
	assertRequest(t, g, "DELETE", "/api/user/apikey/key1")

	_, err = c.APIKeys.Create(ctx(), nil)
	mustError(t, err, "must not be nil")
	mustError(t, c.APIKeys.Delete(ctx(), ""), "api key id is required")
}

func TestAuth(t *testing.T) {
	g, c := newGateway(t)
	g.respond("POST", "/api/token-auth", 200, `{"status":"success","msg":"success","token":"eyJhbGciOi...","sid":"sid123"}`)
	g.respond("GET", "/api/logout", 200, `{"status":"success","msg":"success"}`)
	g.respond("POST", "/api/user/forgot-username", 200, `{"status":"success","msg":"success"}`)
	g.respond("POST", "/api/user/forgot-password", 200, `{"status":"success","msg":"success"}`)
	g.respond("POST", "/api/user/forgot-password/reset", 200, `{"status":"success","msg":"success"}`)

	jwt, err := c.Auth.ObtainJWT(ctx(), &JWTRequest{Username: "u", Password: "p"})
	mustNoError(t, err)
	assertRequest(t, g, "POST", "/api/token-auth")
	equal(t, "body", string(g.last().Body), `{"username":"u","password":"p"}`)
	equal(t, "Token", jwt.Token, "eyJhbGciOi...")
	equal(t, "SID", jwt.SID, "sid123")
	equal(t, "correlation", jwt.LastResponse.CorrelationID, testCorrelationID)

	// The token can be used to build a second client.
	c2, err := NewClient("", WithBaseURL(c.BaseURL()), WithBearerToken(jwt.Token))
	mustNoError(t, err)
	mustNoError(t, c2.Auth.Logout(ctx()))
	r := assertRequestWithAuth(t, g, "GET", "/api/logout", "Bearer eyJhbGciOi...")
	equal(t, "no body", len(r.Body), 0)

	mustNoError(t, c.Auth.ForgotUsername(ctx(), &ForgotUsernameRequest{Email: "e@x.com"}))
	assertRequest(t, g, "POST", "/api/user/forgot-username")
	mustNoError(t, c.Auth.ForgotPassword(ctx(), &ForgotPasswordRequest{Username: "u"}))
	assertRequest(t, g, "POST", "/api/user/forgot-password")
	mustNoError(t, c.Auth.ResetPassword(ctx(), &PasswordResetRequest{Username: "u", ResetCode: "123", Password: "new"}))
	assertRequest(t, g, "POST", "/api/user/forgot-password/reset")
	equal(t, "reset body", string(g.last().Body), `{"username":"u","reset_code":"123","password":"new"}`)

	// Failure modes.
	g.respond("POST", "/api/token-auth", 200, `{"status":"success","msg":"success"}`)
	_, err = c.Auth.ObtainJWT(ctx(), &JWTRequest{Username: "u", Password: "p"})
	mustError(t, err, "did not include a token")

	g.respond("POST", "/api/token-auth", 401, `{"status":"failed","msg":"invalid credentials"}`)
	_, err = c.Auth.ObtainJWT(ctx(), &JWTRequest{Username: "u", Password: "wrong"})
	if !IsUnauthorized(err) {
		t.Fatalf("expected unauthorized, got %v", err)
	}

	_, err = c.Auth.ObtainJWT(ctx(), nil)
	mustError(t, err, "must not be nil")
	mustError(t, c.Auth.ForgotUsername(ctx(), nil), "must not be nil")
	mustError(t, c.Auth.ForgotPassword(ctx(), nil), "must not be nil")
	mustError(t, c.Auth.ResetPassword(ctx(), nil), "must not be nil")
}

func assertRequestWithAuth(t *testing.T, g *gateway, method, path, auth string) recorded {
	t.Helper()
	r := g.last()
	if r.Method != method || r.Path != path {
		t.Fatalf("request = %s %s, want %s %s", r.Method, r.Path, method, path)
	}
	equal(t, "Authorization", r.Header.Get("Authorization"), auth)
	return r
}

func TestDecodeTopLevel(t *testing.T) {
	mustError(t, decodeTopLevel(nil, &struct{}{}), "empty response body")
	mustError(t, decodeTopLevel(&APIResponse{raw: []byte("  ")}, &struct{}{}), "empty response body")
	var out struct {
		Token string `json:"token"`
	}
	mustNoError(t, decodeTopLevel(&APIResponse{raw: []byte(`{"token":"t"}`)}, &out))
	equal(t, "token", out.Token, "t")
}

func TestStatusTextInErrors(t *testing.T) {
	for code, want := range map[int]string{
		http.StatusUnauthorized:    "401 Unauthorized",
		http.StatusNotFound:        "404 Not Found",
		http.StatusTooManyRequests: "429 Too Many Requests",
	} {
		e := &Error{StatusCode: code, Method: "GET", Path: "x"}
		mustError(t, e, want)
	}
}
