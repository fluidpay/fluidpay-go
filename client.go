package fluidpay

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// Version is the SDK version reported in the User-Agent header.
const Version = "1.0.0"

// Base URLs for the two FluidPay environments.
const (
	ProductionURL = "https://app.fluidpay.com/api"
	SandboxURL    = "https://sandbox.fluidpay.com/api"
)

// DefaultTimeout is the HTTP timeout applied when no custom http.Client is
// supplied. FluidPay recommends at least three minutes because some
// authorizations legitimately take longer than a minute to complete.
const DefaultTimeout = 3 * time.Minute

// Environment variables read by NewClientFromEnv.
const (
	EnvAPIKey      = "FLUIDPAY_API_KEY"
	EnvEnvironment = "FLUIDPAY_ENVIRONMENT"
	EnvBaseURL     = "FLUIDPAY_BASE_URL"
)

// Environment selects which gateway a Client talks to.
type Environment string

// Supported environments.
const (
	// Production processes live transactions.
	Production Environment = "production"
	// Sandbox simulates processing. Nothing leaves the platform and no
	// billing occurs. All development should happen here.
	Sandbox Environment = "sandbox"
)

// baseURL returns the API base URL for the environment.
func (e Environment) baseURL() (string, error) {
	switch strings.ToLower(strings.TrimSpace(string(e))) {
	case "", string(Production), "prod", "live":
		return ProductionURL, nil
	case string(Sandbox), "test":
		return SandboxURL, nil
	default:
		return "", fmt.Errorf("fluidpay: unknown environment %q (use %q or %q)", e, Production, Sandbox)
	}
}

// maxResponseBytes bounds how much of a response body the client will read.
const maxResponseBytes = 16 << 20

// Client talks to the FluidPay API. Create one with NewClient or
// NewClientFromEnv and use the service fields to make requests. A Client is
// safe for concurrent use by multiple goroutines.
type Client struct {
	baseURL       *url.URL
	httpClient    *http.Client
	authorization string
	userAgent     string

	// Transactions processes payments: sale, authorize, capture, void,
	// refund, verification and credit, plus lookup and search.
	Transactions *TransactionsService
	// Customers is the Customer Vault: stored customers, addresses and
	// payment methods.
	Customers *CustomersService
	// LegacyCustomers exposes the deprecated /customer endpoints. Prefer
	// Customers for new integrations.
	LegacyCustomers *LegacyCustomersService
	// AddOns manages recurring add-ons.
	AddOns *AddOnsService
	// Discounts manages recurring discounts.
	Discounts *DiscountsService
	// Plans manages recurring plans.
	Plans *PlansService
	// Subscriptions manages recurring subscriptions and their lifecycle.
	Subscriptions *SubscriptionsService
	// Terminals lists physical terminals and settles their batches.
	Terminals *TerminalsService
	// Settlements searches settlement batches.
	Settlements *SettlementsService
	// Lookup performs BIN lookups.
	Lookup *LookupService
	// Users administers user accounts.
	Users *UsersService
	// APIKeys creates, lists and deletes API keys for the current user.
	APIKeys *APIKeysService
	// Auth obtains JWT tokens and handles password recovery.
	Auth *AuthService
}

// Option configures a Client.
type Option func(*Client) error

// WithEnvironment points the client at the given environment. The default
// is Production.
func WithEnvironment(env Environment) Option {
	return func(c *Client) error {
		raw, err := env.baseURL()
		if err != nil {
			return err
		}
		return c.setBaseURL(raw)
	}
}

// WithSandbox points the client at the sandbox environment. It is shorthand
// for WithEnvironment(Sandbox).
func WithSandbox() Option { return WithEnvironment(Sandbox) }

// WithBaseURL overrides the API base URL entirely, for example to target a
// local gateway or a proxy. The URL should include the "/api" path prefix.
func WithBaseURL(raw string) Option {
	return func(c *Client) error { return c.setBaseURL(raw) }
}

// WithHTTPClient uses a custom *http.Client for all requests. Configure its
// Timeout with the gateway's guidance in mind (see DefaultTimeout).
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) error {
		if hc == nil {
			return errors.New("fluidpay: http client must not be nil")
		}
		c.httpClient = hc
		return nil
	}
}

// WithTimeout sets the timeout of the client's *http.Client. It applies to
// the default client and to one supplied via WithHTTPClient.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) error {
		if d <= 0 {
			return errors.New("fluidpay: timeout must be positive")
		}
		c.httpClient.Timeout = d
		return nil
	}
}

// WithUserAgent appends an identifier for your application to the
// User-Agent header, for example "my-shop/2.3".
func WithUserAgent(ua string) Option {
	return func(c *Client) error {
		ua = strings.TrimSpace(ua)
		if ua != "" {
			c.userAgent = ua + " " + defaultUserAgent()
		}
		return nil
	}
}

// WithBearerToken authenticates with a JWT obtained from Auth.ObtainJWT
// instead of an API key. The token is sent as "Authorization: Bearer <token>".
func WithBearerToken(token string) Option {
	return func(c *Client) error {
		token = strings.TrimSpace(token)
		if token == "" {
			return errors.New("fluidpay: bearer token must not be empty")
		}
		c.authorization = "Bearer " + token
		return nil
	}
}

func defaultUserAgent() string { return "fluidpay-go/" + Version }

// NewClient returns a Client authenticated with a private API key. The key
// must start with "api_"; public keys ("pub_") are for browser-side use and
// are rejected by the gateway for server-side calls, so they are rejected
// here too. Pass WithBearerToken to authenticate with a JWT instead, in
// which case apiKey may be empty.
//
// By default the client targets production with a three minute timeout.
func NewClient(apiKey string, opts ...Option) (*Client, error) {
	c := &Client{
		httpClient: &http.Client{Timeout: DefaultTimeout},
		userAgent:  defaultUserAgent(),
	}
	if err := c.setBaseURL(ProductionURL); err != nil {
		return nil, err
	}

	apiKey = strings.TrimSpace(apiKey)
	if strings.HasPrefix(apiKey, "pub_") {
		return nil, errors.New("fluidpay: public API keys (pub_...) are client-side only; use your private api_... key")
	}
	c.authorization = apiKey

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	if c.authorization == "" {
		return nil, errors.New("fluidpay: an API key or bearer token is required")
	}

	c.Transactions = &TransactionsService{client: c}
	c.Customers = &CustomersService{client: c}
	c.LegacyCustomers = &LegacyCustomersService{client: c}
	c.AddOns = &AddOnsService{client: c}
	c.Discounts = &DiscountsService{client: c}
	c.Plans = &PlansService{client: c}
	c.Subscriptions = &SubscriptionsService{client: c}
	c.Terminals = &TerminalsService{client: c}
	c.Settlements = &SettlementsService{client: c}
	c.Lookup = &LookupService{client: c}
	c.Users = &UsersService{client: c}
	c.APIKeys = &APIKeysService{client: c}
	c.Auth = &AuthService{client: c}
	return c, nil
}

// NewClientFromEnv builds a Client from environment variables so secrets stay
// out of source code:
//
//   - FLUIDPAY_API_KEY (required): your private api_... key
//   - FLUIDPAY_ENVIRONMENT (optional): "sandbox" or "production" (default)
//   - FLUIDPAY_BASE_URL (optional): full base URL override
//
// Options passed explicitly are applied after the environment and take
// precedence.
func NewClientFromEnv(opts ...Option) (*Client, error) {
	apiKey := os.Getenv(EnvAPIKey)
	if strings.TrimSpace(apiKey) == "" {
		return nil, fmt.Errorf("fluidpay: %s is not set", EnvAPIKey)
	}
	var envOpts []Option
	if env := os.Getenv(EnvEnvironment); env != "" {
		envOpts = append(envOpts, WithEnvironment(Environment(env)))
	}
	if base := os.Getenv(EnvBaseURL); base != "" {
		envOpts = append(envOpts, WithBaseURL(base))
	}
	return NewClient(apiKey, append(envOpts, opts...)...)
}

// BaseURL returns the API base URL the client sends requests to.
func (c *Client) BaseURL() string { return c.baseURL.String() }

func (c *Client) setBaseURL(raw string) error {
	raw = strings.TrimRight(strings.TrimSpace(raw), "/")
	u, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("fluidpay: invalid base URL %q: %w", raw, err)
	}
	if u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("fluidpay: invalid base URL %q: scheme and host are required", raw)
	}
	c.baseURL = u
	return nil
}

// APIResponse describes the HTTP response that produced a value. It is
// attached to every result through the embedded APIResource so callers can
// reach the correlation ID and raw headers when they need them.
type APIResponse struct {
	// StatusCode is the HTTP status code.
	StatusCode int
	// Header holds the response headers.
	Header http.Header
	// CorrelationID is the x-correlation-id header. Quote it in support
	// requests; FluidPay uses it to find the request in their logs.
	CorrelationID string
	// Status is the envelope status, normally "success".
	Status string
	// Msg is the envelope message.
	Msg string
	// TotalCount is the envelope total_count for list responses.
	TotalCount int

	// raw is the response body, kept for the few endpoints whose payload
	// sits beside status and msg instead of under data.
	raw []byte
}

// decodeTopLevel unmarshals the whole response body into out.
func decodeTopLevel(resp *APIResponse, out any) error {
	if resp == nil || len(bytes.TrimSpace(resp.raw)) == 0 {
		return errors.New("fluidpay: empty response body")
	}
	return json.Unmarshal(resp.raw, out)
}

// APIResource is embedded in every value returned by the SDK and records the
// HTTP response that produced it. It is ignored when marshaling.
type APIResource struct {
	// LastResponse is metadata about the HTTP response, or nil for values
	// that were not produced by an API call.
	LastResponse *APIResponse `json:"-"`
}

func (r *APIResource) setLastResponse(resp *APIResponse) { r.LastResponse = resp }

type lastResponseSetter interface{ setLastResponse(*APIResponse) }

// envelope is the standard FluidPay response wrapper.
type envelope struct {
	Status     string          `json:"status"`
	Msg        string          `json:"msg"`
	Data       json.RawMessage `json:"data"`
	TotalCount int             `json:"total_count"`
	// Count is used instead of total_count by the deprecated customer
	// endpoints.
	Count int `json:"count"`
}

func (e envelope) total() int {
	if e.TotalCount != 0 {
		return e.TotalCount
	}
	return e.Count
}

// call performs an HTTP request and decodes the envelope. If out is non-nil,
// the envelope's data member is decoded into it. The returned APIResponse is
// non-nil whenever an HTTP response was received, including on error.
func (c *Client) call(ctx context.Context, method, path string, query url.Values, body any, out any) (*APIResponse, error) {
	if ctx == nil {
		return nil, errors.New("fluidpay: context must not be nil")
	}
	req, err := c.newRequest(ctx, method, path, query, body)
	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fluidpay: %s %s: %w", method, path, err)
	}
	defer res.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(res.Body, maxResponseBytes))
	if err != nil {
		return nil, fmt.Errorf("fluidpay: %s %s: reading response: %w", method, path, err)
	}

	apiResp := &APIResponse{
		StatusCode:    res.StatusCode,
		Header:        res.Header,
		CorrelationID: res.Header.Get("X-Correlation-Id"),
		raw:           raw,
	}

	var env envelope
	if len(bytes.TrimSpace(raw)) > 0 {
		// A malformed body on an error status is still an error; on a
		// success status it is a decoding failure reported below.
		if jsonErr := json.Unmarshal(raw, &env); jsonErr != nil && res.StatusCode < 400 {
			return apiResp, fmt.Errorf("fluidpay: %s %s: decoding response: %w", method, path, jsonErr)
		}
	}
	apiResp.Status = env.Status
	apiResp.Msg = env.Msg
	apiResp.TotalCount = env.total()

	if res.StatusCode >= 400 || isFailureStatus(env.Status) {
		return apiResp, &Error{
			StatusCode:    res.StatusCode,
			Status:        env.Status,
			Msg:           env.Msg,
			CorrelationID: apiResp.CorrelationID,
			Method:        method,
			Path:          path,
			Body:          raw,
		}
	}

	if out == nil {
		return apiResp, nil
	}
	data := env.Data
	if len(bytes.TrimSpace(data)) == 0 && env.Status == "" {
		// Some endpoints (for example BIN lookup) return the payload
		// without the standard envelope.
		data = raw
	}
	if err := decodeData(data, out); err != nil {
		return apiResp, fmt.Errorf("fluidpay: %s %s: decoding data: %w", method, path, err)
	}
	return apiResp, nil
}

func isFailureStatus(status string) bool {
	switch strings.ToLower(status) {
	case "", "success":
		return false
	default:
		return true
	}
}

// decodeData unmarshals an envelope data member into out. When out is not a
// slice but the gateway wrapped a single record in an array, the first
// element is used.
func decodeData(data json.RawMessage, out any) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		return nil
	}
	if _, wantsSlice := out.(sliceTarget); data[0] == '[' && !wantsSlice {
		var items []json.RawMessage
		if err := json.Unmarshal(data, &items); err != nil {
			return err
		}
		if len(items) == 0 {
			return ErrNotFound
		}
		data = items[0]
	}
	return json.Unmarshal(data, out)
}

// sliceTarget is implemented by list decoders so decodeData knows an array
// is expected.
type sliceTarget interface{ isSliceTarget() }

// newRequest builds an *http.Request with authentication and JSON headers.
func (c *Client) newRequest(ctx context.Context, method, path string, query url.Values, body any) (*http.Request, error) {
	// path arrives with each segment already escaped (see joinPath), so it
	// becomes the raw path and is unescaped for URL.Path.
	rawPath := strings.TrimRight(c.baseURL.EscapedPath(), "/") + "/" + strings.TrimLeft(path, "/")
	unescaped, err := url.PathUnescape(rawPath)
	if err != nil {
		return nil, fmt.Errorf("fluidpay: %s %s: invalid path: %w", method, path, err)
	}
	u := *c.baseURL
	u.Path = unescaped
	u.RawPath = rawPath
	if len(query) > 0 {
		u.RawQuery = query.Encode()
	}

	var payload io.Reader
	if body != nil {
		buf, err := encodeJSON(body)
		if err != nil {
			return nil, fmt.Errorf("fluidpay: %s %s: encoding request: %w", method, path, err)
		}
		payload = bytes.NewReader(buf)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), payload)
	if err != nil {
		return nil, fmt.Errorf("fluidpay: %s %s: %w", method, path, err)
	}
	req.Header.Set("Authorization", c.authorization)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

// encodeJSON marshals v without HTML escaping so search operators such as
// "<" are sent literally.
func encodeJSON(v any) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return bytes.TrimRight(buf.Bytes(), "\n"), nil
}

// doOne performs a request whose data member is a single record of type T.
func doOne[T any](ctx context.Context, c *Client, method, path string, query url.Values, body any) (*T, error) {
	var out T
	resp, err := c.call(ctx, method, path, query, body, &out)
	if err != nil {
		return nil, err
	}
	if s, ok := any(&out).(lastResponseSetter); ok {
		s.setLastResponse(resp)
	}
	return &out, nil
}

// listDecoder lets decodeData recognise a slice destination.
type listDecoder[T any] struct{ items []*T }

func (*listDecoder[T]) isSliceTarget() {}

func (l *listDecoder[T]) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &l.items)
}

// doList performs a request whose data member is an array of T.
func doList[T any](ctx context.Context, c *Client, method, path string, query url.Values, body any) ([]*T, *APIResponse, error) {
	var out listDecoder[T]
	resp, err := c.call(ctx, method, path, query, body, &out)
	if err != nil {
		return nil, resp, err
	}
	for _, item := range out.items {
		if s, ok := any(item).(lastResponseSetter); ok {
			s.setLastResponse(resp)
		}
	}
	return out.items, resp, nil
}

// doEmpty performs a request whose response carries no data of interest.
func doEmpty(ctx context.Context, c *Client, method, path string, query url.Values, body any) error {
	_, err := c.call(ctx, method, path, query, body, nil)
	return err
}

// joinPath escapes each segment and joins them with "/".
func joinPath(segments ...string) string {
	parts := make([]string, 0, len(segments))
	for _, s := range segments {
		parts = append(parts, url.PathEscape(s))
	}
	return strings.Join(parts, "/")
}

// requireID returns an error when a path identifier is blank so a typo never
// turns "get one" into "list all".
func requireID(name, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("fluidpay: %s is required", name)
	}
	return nil
}

// requireIDs validates several path identifiers given as name/value pairs.
func requireIDs(pairs ...string) error {
	for i := 0; i+1 < len(pairs); i += 2 {
		if err := requireID(pairs[i], pairs[i+1]); err != nil {
			return err
		}
	}
	return nil
}
