package fluidpay

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

const (
	testAPIKey        = "api_test_key_for_unit_tests"
	testCorrelationID = "corr-1234"
)

// recorded is the last request the fake gateway received.
type recorded struct {
	Method string
	Path   string
	// RawPath is the path exactly as sent on the wire.
	RawPath string
	Query   string
	Header  http.Header
	Body    []byte
}

// gateway is an in-process stand-in for the FluidPay API.
type gateway struct {
	t      *testing.T
	server *httptest.Server
	mu     sync.Mutex
	routes map[string]http.HandlerFunc
	calls  []recorded
}

// newGateway starts a fake gateway and returns it with a Client pointed at
// it. Routes are registered with respond or handle.
func newGateway(t *testing.T, opts ...Option) (*gateway, *Client) {
	t.Helper()
	g := &gateway{t: t, routes: map[string]http.HandlerFunc{}}
	g.server = httptest.NewServer(http.HandlerFunc(g.serve))
	t.Cleanup(g.server.Close)

	opts = append([]Option{WithBaseURL(g.server.URL + "/api")}, opts...)
	client, err := NewClient(testAPIKey, opts...)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return g, client
}

func (g *gateway) serve(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	g.mu.Lock()
	g.calls = append(g.calls, recorded{
		Method:  r.Method,
		Path:    r.URL.Path,
		RawPath: r.URL.EscapedPath(),
		Query:   r.URL.RawQuery,
		Header:  r.Header.Clone(),
		Body:    body,
	})
	handler := g.routes[r.Method+" "+r.URL.Path]
	g.mu.Unlock()

	if handler == nil {
		g.t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"status":"failed","msg":"not found"}`))
		return
	}
	// Hand the body back to the handler so it can decode it.
	r.Body = io.NopCloser(strings.NewReader(string(body)))
	handler(w, r)
}

// respond registers a static JSON response for METHOD path.
func (g *gateway) respond(method, path string, status int, body string) {
	g.handle(method, path, func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, status, body)
	})
}

// handle registers a handler for METHOD path.
func (g *gateway) handle(method, path string, fn http.HandlerFunc) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.routes[method+" "+path] = fn
}

// last returns the most recent request.
func (g *gateway) last() recorded {
	g.mu.Lock()
	defer g.mu.Unlock()
	if len(g.calls) == 0 {
		g.t.Fatal("no requests were made")
	}
	return g.calls[len(g.calls)-1]
}

// count returns how many requests were made.
func (g *gateway) count() int {
	g.mu.Lock()
	defer g.mu.Unlock()
	return len(g.calls)
}

// bodyJSON decodes the most recent request body.
func (g *gateway) bodyJSON() map[string]any {
	g.t.Helper()
	r := g.last()
	if len(r.Body) == 0 {
		return nil
	}
	var m map[string]any
	if err := json.Unmarshal(r.Body, &m); err != nil {
		g.t.Fatalf("request body is not a JSON object: %v\n%s", err, r.Body)
	}
	return m
}

func writeJSON(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Correlation-Id", testCorrelationID)
	w.WriteHeader(status)
	_, _ = io.WriteString(w, body)
}

func wrap(data string) string {
	return `{"status":"success","msg":"success","data":` + data + `}`
}

func listEnvelope(data string, total int) string {
	return `{"status":"success","msg":"","total_count":` + itoa(total) + `,"data":` + data + `}`
}

func itoa(i int) string {
	b, _ := json.Marshal(i)
	return string(b)
}

func ctx() context.Context { return context.Background() }

// assertRequest checks the method and path of the most recent request.
func assertRequest(t *testing.T, g *gateway, method, path string) recorded {
	t.Helper()
	r := g.last()
	if r.Method != method || r.Path != path {
		t.Fatalf("request = %s %s, want %s %s", r.Method, r.Path, method, path)
	}
	if got := r.Header.Get("Authorization"); got != testAPIKey {
		t.Errorf("Authorization = %q, want %q", got, testAPIKey)
	}
	return r
}

func mustNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func mustError(t *testing.T, err error, contains string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected an error containing %q, got nil", contains)
	}
	if !strings.Contains(err.Error(), contains) {
		t.Fatalf("error %q does not contain %q", err.Error(), contains)
	}
}

// equal compares two comparable values; dynamic types must match too.
func equal(t *testing.T, name string, got, want any) {
	t.Helper()
	if got != want {
		t.Errorf("%s = %#v, want %#v", name, got, want)
	}
}

// keys returns the sorted top-level keys of a JSON object.
func keys(m map[string]any) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sortStrings(out)
	return out
}

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}

func joinKeys(m map[string]any) string { return strings.Join(keys(m), ",") }
