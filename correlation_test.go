package fluidpay

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"testing"
)

// These tests pin down every path through which the x-correlation-id of a
// gateway response is exposed to callers.

func TestCorrelationID_OnSuccessfulValues(t *testing.T) {
	g, c := newGateway(t)
	g.respond("POST", "/api/transaction", 200, cardSaleFixture)
	g.respond("GET", "/api/terminals", 200, terminalsFixture)
	g.respond("GET", "/api/vault/cust1", 200, vaultCustomerFixture)
	g.respond("POST", "/api/token-auth", 200, `{"status":"success","msg":"success","token":"t","sid":"s"}`)
	g.respond("POST", "/api/lookup/bin/protected", 200, binLookupFixture)

	tx, err := c.Transactions.Sale(ctx(), &TransactionRequest{Amount: 1})
	mustNoError(t, err)
	equal(t, "Transaction.CorrelationID()", tx.CorrelationID(), testCorrelationID)
	equal(t, "LastResponse.CorrelationID", tx.LastResponse.CorrelationID, testCorrelationID)
	equal(t, "LastResponse.Method", tx.LastResponse.Method, "POST")
	equal(t, "LastResponse.Path", tx.LastResponse.Path, "transaction")
	equal(t, "raw header", tx.LastResponse.Header.Get(CorrelationIDHeader), testCorrelationID)

	list, err := c.Terminals.List(ctx())
	mustNoError(t, err)
	equal(t, "list", list.CorrelationID(), testCorrelationID)
	equal(t, "list item", list.Data[0].CorrelationID(), testCorrelationID)

	// Types with custom decoding keep it too.
	cust, err := c.Customers.Get(ctx(), "cust1")
	mustNoError(t, err)
	equal(t, "customer", cust.CorrelationID(), testCorrelationID)

	jwt, err := c.Auth.ObtainJWT(ctx(), &JWTRequest{Username: "u", Password: "p"})
	mustNoError(t, err)
	equal(t, "jwt", jwt.CorrelationID(), testCorrelationID)

	bin, err := c.Lookup.BIN(ctx(), &BINLookupRequest{BIN: "424242"})
	mustNoError(t, err)
	equal(t, "bin", bin.CorrelationID(), testCorrelationID)

	// Values not produced by an API call report "" rather than panicking.
	equal(t, "zero value", (&Transaction{}).CorrelationID(), "")
	var nilRes *APIResource
	equal(t, "nil receiver", nilRes.CorrelationID(), "")
}

func TestCorrelationID_OnEmptyResponses(t *testing.T) {
	g, c := newGateway(t)
	g.respond("POST", "/api/transaction/tx1/void", 200, `{"status":"success","msg":"success","data":null}`)
	g.respond("DELETE", "/api/vault/cust1", 200, `{"status":"success","msg":"success"}`)
	g.respond("DELETE", "/api/recurring/plan/p1", 200, `{"status":"success","msg":"success","data":null}`)
	g.respond("POST", "/api/terminal/t1/settle", 200, `{"status":"success","msg":"success"}`)
	g.respond("GET", "/api/logout", 200, `{"status":"success","msg":"success"}`)
	g.respond("POST", "/api/customer/c1", 200, `{"status":"success","msg":"success","data":null}`)

	calls := []struct {
		name string
		call func() (*APIResponse, error)
		path string
	}{
		{"Void", func() (*APIResponse, error) { return c.Transactions.Void(ctx(), "tx1") }, "transaction/tx1/void"},
		{"Customers.Delete", func() (*APIResponse, error) { return c.Customers.Delete(ctx(), "cust1") }, "vault/cust1"},
		{"Plans.Delete", func() (*APIResponse, error) { return c.Plans.Delete(ctx(), "p1") }, "recurring/plan/p1"},
		{"Terminals.Settle", func() (*APIResponse, error) { return c.Terminals.Settle(ctx(), "t1") }, "terminal/t1/settle"},
		{"Auth.Logout", func() (*APIResponse, error) { return c.Auth.Logout(ctx()) }, "logout"},
		{"LegacyCustomers.Update", func() (*APIResponse, error) {
			return c.LegacyCustomers.Update(ctx(), "c1", &LegacyCustomerUpdateRequest{Description: "d"})
		}, "customer/c1"},
	}
	for _, tc := range calls {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := tc.call()
			mustNoError(t, err)
			if resp == nil {
				t.Fatal("response metadata is nil")
			}
			equal(t, "CorrelationID", resp.CorrelationID, testCorrelationID)
			equal(t, "StatusCode", resp.StatusCode, 200)
			equal(t, "Status", resp.Status, "success")
			equal(t, "Path", resp.Path, tc.path)
		})
	}

	// Validation failures happen before any request, so there is no
	// response and no correlation ID.
	resp, err := c.Transactions.Void(ctx(), "")
	mustError(t, err, "transaction id is required")
	if resp != nil {
		t.Error("expected nil response for a validation error")
	}
	equal(t, "CorrelationID(validation err)", CorrelationID(err), "")
}

func TestCorrelationID_OnErrors(t *testing.T) {
	g, c := newGateway(t)
	g.respond("POST", "/api/transaction", 400, `{"status":"failed","msg":"bad request error: invalid Postal Code"}`)
	g.respond("GET", "/api/terminals", 200, `{"status":"success","msg":"","data":[{"id":`) // truncated body
	g.respond("GET", "/api/transaction/x", 200, `{"status":"success","msg":"","data":{"amount":"not a number"}}`)

	// Gateway rejection.
	_, err := c.Transactions.Sale(ctx(), &TransactionRequest{Amount: 1})
	equal(t, "CorrelationID(*Error)", CorrelationID(err), testCorrelationID)
	apiErr, _ := AsError(err)
	equal(t, "Error.CorrelationID", apiErr.CorrelationID, testCorrelationID)
	resp := apiErr.Response()
	equal(t, "Response().CorrelationID", resp.CorrelationID, testCorrelationID)
	equal(t, "Response().StatusCode", resp.StatusCode, 400)
	equal(t, "Response().Method", resp.Method, "POST")
	equal(t, "Response().Header", resp.Header.Get(CorrelationIDHeader), testCorrelationID)
	equal(t, "Response().Msg", resp.Msg, "bad request error: invalid Postal Code")

	// Wrapped errors still resolve.
	wrapped := fmt.Errorf("charging order 42: %w", err)
	equal(t, "CorrelationID(wrapped)", CorrelationID(wrapped), testCorrelationID)

	// Undecodable envelope.
	_, err = c.Terminals.List(ctx())
	mustError(t, err, "decoding response")
	mustError(t, err, "(correlation id "+testCorrelationID+")")
	equal(t, "CorrelationID(decode envelope)", CorrelationID(err), testCorrelationID)
	if _, isAPIErr := AsError(err); isAPIErr {
		t.Error("a decoding failure must not masquerade as a gateway rejection")
	}

	// Undecodable data member.
	_, err = c.Transactions.Get(ctx(), "x")
	mustError(t, err, "decoding response")
	equal(t, "CorrelationID(decode data)", CorrelationID(err), testCorrelationID)
	var jsonErr interface{ Unwrap() error }
	if !errors.As(err, &jsonErr) {
		t.Error("decode error should unwrap to the json error")
	}

	// Errors with no response attached.
	equal(t, "plain error", CorrelationID(errors.New("boom")), "")
	equal(t, "nil error", CorrelationID(nil), "")
	equal(t, "ErrNotFound", CorrelationID(ErrNotFound), "")
}

func TestCorrelationID_HeaderNameFallback(t *testing.T) {
	g, c := newGateway(t)
	g.handle("GET", "/api/terminals", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Correlation-Id", "short-form")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(terminalsFixture))
	})
	list, err := c.Terminals.List(ctx())
	mustNoError(t, err)
	equal(t, "fallback header", list.CorrelationID(), "short-form")

	g.handle("GET", "/api/terminals", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("x-correlation-id", "lower-case")
		w.Header().Set("Correlation-Id", "ignored")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(terminalsFixture))
	})
	list, err = c.Terminals.List(ctx())
	mustNoError(t, err)
	equal(t, "documented header wins", list.CorrelationID(), "lower-case")

	g.handle("GET", "/api/terminals", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte(terminalsFixture))
	})
	list, err = c.Terminals.List(ctx())
	mustNoError(t, err)
	equal(t, "absent header", list.CorrelationID(), "")
}

func TestWithResponseHook(t *testing.T) {
	var (
		mu   sync.Mutex
		seen []*APIResponse
	)
	hook := func(r *APIResponse) {
		mu.Lock()
		defer mu.Unlock()
		seen = append(seen, r)
	}
	var order []string
	g, c := newGateway(t, WithResponseHook(hook), WithResponseHook(func(r *APIResponse) {
		mu.Lock()
		defer mu.Unlock()
		order = append(order, "second:"+r.Path)
	}))
	g.respond("POST", "/api/transaction", 200, cardSaleFixture)
	g.respond("POST", "/api/transaction/tx1/void", 200, `{"status":"success","msg":"success","data":null}`)
	g.respond("GET", "/api/terminals", 401, `{"status":"failed","msg":"unauthorized"}`)
	g.respond("GET", "/api/user", 200, `{"status":"success",`)

	_, _ = c.Transactions.Sale(ctx(), &TransactionRequest{Amount: 1})
	_, _ = c.Transactions.Void(ctx(), "tx1")
	_, _ = c.Terminals.List(ctx())
	_, _ = c.Users.Current(ctx())
	_, _ = c.Transactions.Void(ctx(), "") // validation error: no request, no hook

	mu.Lock()
	defer mu.Unlock()
	equal(t, "hook count", len(seen), 4)
	equal(t, "second hook count", len(order), 4)
	equal(t, "second hook order", order[0], "second:transaction")

	want := []struct {
		method, path string
		status       int
	}{
		{"POST", "transaction", 200},
		{"POST", "transaction/tx1/void", 200},
		{"GET", "terminals", 401},
		{"GET", "user", 200},
	}
	for i, w := range want {
		equal(t, fmt.Sprintf("[%d] method", i), seen[i].Method, w.method)
		equal(t, fmt.Sprintf("[%d] path", i), seen[i].Path, w.path)
		equal(t, fmt.Sprintf("[%d] status", i), seen[i].StatusCode, w.status)
		equal(t, fmt.Sprintf("[%d] correlation", i), seen[i].CorrelationID, testCorrelationID)
	}
	// The failed call's envelope fields are visible to the hook as well.
	equal(t, "failed status", seen[2].Status, "failed")
	equal(t, "failed msg", seen[2].Msg, "unauthorized")

	_, err := NewClient("api_x", WithResponseHook(nil))
	mustError(t, err, "response hook must not be nil")
}

func TestWithResponseHook_NotCalledOnTransportFailure(t *testing.T) {
	called := false
	c, err := NewClient("api_x", WithBaseURL("http://127.0.0.1:1/api"), WithResponseHook(func(*APIResponse) { called = true }))
	mustNoError(t, err)
	_, err = c.Terminals.List(ctx())
	if err == nil || called {
		t.Fatalf("expected transport error without hook call; err=%v called=%v", err, called)
	}
	equal(t, "no correlation id", CorrelationID(err), "")
}
