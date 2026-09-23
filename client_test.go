package fluidpay

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestNewClient_Defaults(t *testing.T) {
	c, err := NewClient("api_abc")
	mustNoError(t, err)
	equal(t, "BaseURL", c.BaseURL(), ProductionURL)
	equal(t, "timeout", c.httpClient.Timeout, DefaultTimeout)
	equal(t, "authorization", c.authorization, "api_abc")
	equal(t, "user agent", c.userAgent, "fluidpay-go/"+Version)
	for name, svc := range map[string]any{
		"Transactions": c.Transactions, "Customers": c.Customers, "LegacyCustomers": c.LegacyCustomers,
		"AddOns": c.AddOns, "Discounts": c.Discounts, "Plans": c.Plans, "Subscriptions": c.Subscriptions,
		"Terminals": c.Terminals, "Settlements": c.Settlements, "Lookup": c.Lookup,
		"Users": c.Users, "APIKeys": c.APIKeys, "Auth": c.Auth,
	} {
		if svc == nil {
			t.Errorf("service %s is nil", name)
		}
	}
}

func TestNewClient_Validation(t *testing.T) {
	tests := []struct {
		name string
		key  string
		opts []Option
		want string
	}{
		{"public key", "pub_abc", nil, "client-side only"},
		{"empty key", "", nil, "API key or bearer token is required"},
		{"whitespace key", "   ", nil, "API key or bearer token is required"},
		{"bad environment", "api_abc", []Option{WithEnvironment("staging")}, "unknown environment"},
		{"bad base url", "api_abc", []Option{WithBaseURL("not a url")}, "invalid base URL"},
		{"relative base url", "api_abc", []Option{WithBaseURL("/api")}, "scheme and host"},
		{"nil http client", "api_abc", []Option{WithHTTPClient(nil)}, "must not be nil"},
		{"zero timeout", "api_abc", []Option{WithTimeout(0)}, "must be positive"},
		{"empty bearer", "", []Option{WithBearerToken("  ")}, "bearer token must not be empty"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewClient(tt.key, tt.opts...)
			mustError(t, err, tt.want)
		})
	}
}

func TestNewClient_Options(t *testing.T) {
	hc := &http.Client{Timeout: time.Second}
	c, err := NewClient("api_abc",
		nil, // nil options are ignored
		WithSandbox(),
		WithHTTPClient(hc),
		WithTimeout(10*time.Second),
		WithUserAgent("my-shop/2.3"),
	)
	mustNoError(t, err)
	equal(t, "BaseURL", c.BaseURL(), SandboxURL)
	if c.httpClient != hc {
		t.Error("custom http client not used")
	}
	equal(t, "timeout", hc.Timeout, 10*time.Second)
	equal(t, "user agent", c.userAgent, "my-shop/2.3 fluidpay-go/"+Version)

	c, err = NewClient("api_abc", WithEnvironment(Production))
	mustNoError(t, err)
	equal(t, "production", c.BaseURL(), ProductionURL)

	c, err = NewClient("api_abc", WithBaseURL("https://gateway.example.com/api/"))
	mustNoError(t, err)
	equal(t, "trailing slash trimmed", c.BaseURL(), "https://gateway.example.com/api")

	c, err = NewClient("", WithBearerToken("jwt-token"))
	mustNoError(t, err)
	equal(t, "bearer", c.authorization, "Bearer jwt-token")
}

func TestEnvironment_BaseURL(t *testing.T) {
	for env, want := range map[Environment]string{
		"": ProductionURL, Production: ProductionURL, "PROD": ProductionURL, "live": ProductionURL,
		Sandbox: SandboxURL, " Sandbox ": SandboxURL, "test": SandboxURL,
	} {
		got, err := env.baseURL()
		mustNoError(t, err)
		equal(t, string(env), got, want)
	}
}

func TestNewClientFromEnv(t *testing.T) {
	t.Setenv(EnvAPIKey, "")
	t.Setenv(EnvEnvironment, "")
	t.Setenv(EnvBaseURL, "")
	_, err := NewClientFromEnv()
	mustError(t, err, EnvAPIKey+" is not set")

	t.Setenv(EnvAPIKey, "api_from_env")
	c, err := NewClientFromEnv()
	mustNoError(t, err)
	equal(t, "default env", c.BaseURL(), ProductionURL)
	equal(t, "key", c.authorization, "api_from_env")

	t.Setenv(EnvEnvironment, "sandbox")
	c, err = NewClientFromEnv()
	mustNoError(t, err)
	equal(t, "sandbox env", c.BaseURL(), SandboxURL)

	t.Setenv(EnvBaseURL, "http://localhost:8001/api")
	c, err = NewClientFromEnv()
	mustNoError(t, err)
	equal(t, "base url env", c.BaseURL(), "http://localhost:8001/api")

	// Explicit options win over the environment.
	c, err = NewClientFromEnv(WithBaseURL("http://127.0.0.1:9/api"))
	mustNoError(t, err)
	equal(t, "explicit option", c.BaseURL(), "http://127.0.0.1:9/api")
}

func TestClient_RequestHeaders(t *testing.T) {
	g, c := newGateway(t, WithUserAgent("shop/1"))
	g.respond("POST", "/api/transaction", 200, cardSaleFixture)
	g.respond("GET", "/api/terminals", 200, terminalsFixture)

	_, err := c.Transactions.Sale(ctx(), &TransactionRequest{Amount: 100})
	mustNoError(t, err)
	r := g.last()
	equal(t, "Authorization", r.Header.Get("Authorization"), testAPIKey)
	equal(t, "Content-Type", r.Header.Get("Content-Type"), "application/json")
	equal(t, "Accept", r.Header.Get("Accept"), "application/json")
	equal(t, "User-Agent", r.Header.Get("User-Agent"), "shop/1 fluidpay-go/"+Version)

	_, err = c.Terminals.List(ctx())
	mustNoError(t, err)
	r = g.last()
	equal(t, "GET has no body", len(r.Body), 0)
	equal(t, "GET has no Content-Type", r.Header.Get("Content-Type"), "")
}

func TestClient_BearerTokenHeader(t *testing.T) {
	g, _ := newGateway(t)
	c, err := NewClient("", WithBaseURL(g.server.URL+"/api"), WithBearerToken("jwt"))
	mustNoError(t, err)
	g.respond("GET", "/api/user", 200, userFixture)
	_, err = c.Users.Current(ctx())
	mustNoError(t, err)
	equal(t, "Authorization", g.last().Header.Get("Authorization"), "Bearer jwt")
}

func TestClient_HTTPErrorBecomesError(t *testing.T) {
	g, c := newGateway(t)
	g.respond("POST", "/api/transaction", 401, `{"status":"failed","msg":"unauthorized"}`)

	_, err := c.Transactions.Sale(ctx(), &TransactionRequest{Amount: 100})
	apiErr, ok := AsError(err)
	if !ok {
		t.Fatalf("expected *Error, got %T: %v", err, err)
	}
	equal(t, "StatusCode", apiErr.StatusCode, 401)
	equal(t, "Status", apiErr.Status, "failed")
	equal(t, "Msg", apiErr.Msg, "unauthorized")
	equal(t, "CorrelationID", apiErr.CorrelationID, testCorrelationID)
	equal(t, "Method", apiErr.Method, "POST")
	equal(t, "Path", apiErr.Path, "transaction")
	if !apiErr.IsUnauthorized() || !IsUnauthorized(err) {
		t.Error("expected unauthorized")
	}
	want := "fluidpay: POST /transaction: 401 Unauthorized: unauthorized (correlation id corr-1234)"
	equal(t, "Error()", err.Error(), want)
}

func TestClient_FailedEnvelopeWith200(t *testing.T) {
	g, c := newGateway(t)
	g.respond("POST", "/api/transaction", 200, `{"status":"failed","msg":"bad request error: invalid Postal Code"}`)

	_, err := c.Transactions.Sale(ctx(), &TransactionRequest{Amount: 100})
	apiErr, ok := AsError(err)
	if !ok {
		t.Fatalf("expected *Error, got %T", err)
	}
	equal(t, "StatusCode", apiErr.StatusCode, 200)
	equal(t, "Msg", apiErr.Msg, "bad request error: invalid Postal Code")
	if apiErr.IsBadRequest() {
		t.Error("200 should not report bad request")
	}
}

func TestClient_ErrorWithNonJSONBody(t *testing.T) {
	g, c := newGateway(t)
	g.handle("GET", "/api/terminals", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte("<html>bad gateway</html>"))
	})
	_, err := c.Terminals.List(ctx())
	apiErr, ok := AsError(err)
	if !ok {
		t.Fatalf("expected *Error, got %T", err)
	}
	equal(t, "StatusCode", apiErr.StatusCode, 502)
	equal(t, "Body", string(apiErr.Body), "<html>bad gateway</html>")
	if !strings.Contains(err.Error(), "502 Bad Gateway") {
		t.Errorf("unexpected message %q", err.Error())
	}
}

func TestClient_NotFound(t *testing.T) {
	g, c := newGateway(t)
	g.respond("GET", "/api/transaction/missing", 404, `{"status":"failed","msg":"not found"}`)
	_, err := c.Transactions.Get(ctx(), "missing")
	if !IsNotFound(err) {
		t.Fatalf("expected not found, got %v", err)
	}
	if !IsNotFound(ErrNotFound) {
		t.Error("ErrNotFound should be reported as not found")
	}
	if IsNotFound(errors.New("other")) || IsUnauthorized(errors.New("other")) {
		t.Error("unrelated errors must not match")
	}
}

func TestClient_MalformedSuccessBody(t *testing.T) {
	g, c := newGateway(t)
	g.respond("GET", "/api/terminals", 200, `{"status":"success","data":[`)
	_, err := c.Terminals.List(ctx())
	mustError(t, err, "decoding response")
}

func TestClient_EmptySuccessBody(t *testing.T) {
	g, c := newGateway(t)
	g.respond("POST", "/api/transaction/abc/void", 200, "")
	mustNoError(t, errOf(c.Transactions.Void(ctx(), "abc")))
}

func TestClient_ContextCancellation(t *testing.T) {
	g, c := newGateway(t)
	g.handle("GET", "/api/terminals", func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	})
	cctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	_, err := c.Terminals.List(cctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected deadline exceeded, got %v", err)
	}
}

func TestClient_NilContext(t *testing.T) {
	_, c := newGateway(t)
	//nolint:staticcheck // deliberately passing nil to exercise the guard
	_, err := c.Terminals.List(nil)
	mustError(t, err, "context must not be nil")
}

func TestClient_TransportError(t *testing.T) {
	c, err := NewClient("api_abc", WithBaseURL("http://127.0.0.1:1/api"), WithTimeout(time.Second))
	mustNoError(t, err)
	_, err = c.Terminals.List(ctx())
	if err == nil {
		t.Fatal("expected a transport error")
	}
	if _, ok := AsError(err); ok {
		t.Error("transport failures must not be *Error")
	}
}

func TestClient_ArrayWrappedSingleRecord(t *testing.T) {
	g, c := newGateway(t)
	g.respond("GET", "/api/transaction/one", 200, getTransactionFixture)
	tx, err := c.Transactions.Get(ctx(), "one")
	mustNoError(t, err)
	equal(t, "ID", tx.ID, "b7kgflt1tlv51er0fts0")

	g.respond("GET", "/api/transaction/none", 200, `{"status":"success","msg":"","data":[],"total_count":0}`)
	_, err = c.Transactions.Get(ctx(), "none")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestClient_NullData(t *testing.T) {
	g, c := newGateway(t)
	g.respond("GET", "/api/transaction/null", 200, `{"status":"success","msg":"success","data":null}`)
	tx, err := c.Transactions.Get(ctx(), "null")
	mustNoError(t, err)
	equal(t, "zero value", tx.ID, "")
	equal(t, "status", tx.LastResponse.Status, "success")
}

func TestClient_LastResponseMetadata(t *testing.T) {
	g, c := newGateway(t)
	g.respond("GET", "/api/terminals", 200, terminalsFixture)
	list, err := c.Terminals.List(ctx())
	mustNoError(t, err)
	equal(t, "list correlation", list.LastResponse.CorrelationID, testCorrelationID)
	equal(t, "list status code", list.LastResponse.StatusCode, 200)
	equal(t, "list total", list.LastResponse.TotalCount, 1)
	equal(t, "item correlation", list.Data[0].LastResponse.CorrelationID, testCorrelationID)
	equal(t, "header", list.LastResponse.Header.Get("Content-Type"), "application/json")
}

func TestClient_PathEscaping(t *testing.T) {
	g, c := newGateway(t)
	g.respond("GET", "/api/transaction/a/b c", 200, cardSaleFixture)
	_, err := c.Transactions.Get(ctx(), "a/b c")
	mustNoError(t, err)
	equal(t, "decoded path", g.last().Path, "/api/transaction/a/b c")
	equal(t, "wire path", g.last().RawPath, "/api/transaction/a%2Fb%20c")

	// Plain identifiers are sent untouched.
	g.respond("GET", "/api/transaction/plain-id_1", 200, cardSaleFixture)
	_, err = c.Transactions.Get(ctx(), "plain-id_1")
	mustNoError(t, err)
	equal(t, "plain wire path", g.last().RawPath, "/api/transaction/plain-id_1")
}

func TestClient_BaseURLWithPathPrefix(t *testing.T) {
	g, _ := newGateway(t)
	c, err := NewClient("api_abc", WithBaseURL(g.server.URL+"/proxy/api/"))
	mustNoError(t, err)
	g.respond("GET", "/proxy/api/terminals", 200, terminalsFixture)
	_, err = c.Terminals.List(ctx())
	mustNoError(t, err)
}

func TestJoinPath(t *testing.T) {
	equal(t, "plain", joinPath("recurring", "plan", "abc"), "recurring/plan/abc")
	equal(t, "escaped", joinPath("customer", "a b/c?d"), "customer/a%20b%2Fc%3Fd")
}

func TestRequireID(t *testing.T) {
	mustError(t, requireID("plan id", "  "), "plan id is required")
	mustNoError(t, requireID("plan id", "x"))
	mustError(t, requireIDs("a", "1", "b", ""), "b is required")
	mustNoError(t, requireIDs("a", "1", "b", "2"))
}

func TestError_Error(t *testing.T) {
	e := &Error{StatusCode: 400, Method: "POST", Path: "transaction", Status: "failed"}
	equal(t, "status only", e.Error(), "fluidpay: POST /transaction: 400 Bad Request: status failed")
	e = &Error{StatusCode: 599, Method: "GET", Path: "/x"}
	equal(t, "unknown status text", e.Error(), "fluidpay: GET /x: 599")
	if !(&Error{StatusCode: 400}).IsBadRequest() {
		t.Error("IsBadRequest")
	}
}

func TestAPIResource_JSONIgnored(t *testing.T) {
	// LastResponse must never be serialised into a request.
	tx := Transaction{APIResource: APIResource{LastResponse: &APIResponse{StatusCode: 200}}, ID: "x"}
	b, err := jsonMarshal(tx)
	mustNoError(t, err)
	if strings.Contains(string(b), "LastResponse") {
		t.Errorf("LastResponse leaked into JSON: %s", b)
	}
}
