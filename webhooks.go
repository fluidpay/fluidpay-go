package fluidpay

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

// Webhook header names.
const (
	// WebhookSignatureHeader carries the HMAC-SHA256 of the raw body.
	WebhookSignatureHeader = "Signature"
	// WebhookSignaturePreviousHeader carries a signature made with the
	// previous secret during rotation.
	WebhookSignaturePreviousHeader = "Signature-Previous"
	// WebhookEventIDHeader is a stable ID for idempotent processing.
	WebhookEventIDHeader = "X-Webhook-Event-Id"
	// WebhookReplayOfHeader names the original event when a delivery is a
	// replay.
	WebhookReplayOfHeader = "X-Webhook-Replay-Of"
)

// WebhookTestSecret signs deliveries of type "test", which the gateway sends
// when a webhook is created, updated or tested. Verify them with this key
// rather than skipping verification.
const WebhookTestSecret = "00000000-0000-4000-8000-000000000001"

// Webhook event types.
const (
	WebhookTypeTest              = "test"
	WebhookTypeTransactionCreate = "transaction_create"
)

// maxWebhookBody bounds the request body read by ParseWebhookRequest.
const maxWebhookBody = 4 << 20

// ErrInvalidWebhookSignature is returned when a webhook's signature does not
// match any of the supplied secrets.
var ErrInvalidWebhookSignature = errors.New("fluidpay: invalid webhook signature")

// ComputeWebhookSignature returns the signature the gateway would send for
// body under secret: HMAC-SHA256 encoded with unpadded base64url.
func ComputeWebhookSignature(secret string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

// VerifyWebhookSignature reports whether signature matches body under
// secret using a constant-time comparison. Pass the raw request bytes
// exactly as received; the gateway includes a trailing newline in the body
// it signs.
func VerifyWebhookSignature(secret string, body []byte, signature string) bool {
	if secret == "" || signature == "" {
		return false
	}
	expected := ComputeWebhookSignature(secret, body)
	return hmac.Equal([]byte(expected), []byte(signature))
}

// WebhookEvent is the envelope of a webhook delivery.
type WebhookEvent struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
	// Type is the event subtype, for example WebhookTypeTransactionCreate.
	Type          string `json:"type"`
	AccountType   string `json:"account_type"`
	AccountTypeID string `json:"account_type_id"`
	// TransactionID is present for transaction related events.
	TransactionID string    `json:"transaction_id,omitempty"`
	ActionAt      time.Time `json:"action_at"`
	// Data is the event specific payload. Use Transaction to decode it
	// for transaction events.
	Data json.RawMessage `json:"data,omitempty"`

	// EventID and ReplayOf come from the delivery headers when the event
	// was parsed with ParseWebhookRequest.
	EventID  string `json:"-"`
	ReplayOf string `json:"-"`
}

// IsTest reports whether this is a test delivery.
func (e *WebhookEvent) IsTest() bool { return e.Type == WebhookTypeTest }

// Transaction decodes Data as a transaction. It returns an error when the
// event carries no data.
func (e *WebhookEvent) Transaction() (*Transaction, error) {
	if len(e.Data) == 0 {
		return nil, errors.New("fluidpay: webhook event has no data")
	}
	var tx Transaction
	if err := json.Unmarshal(e.Data, &tx); err != nil {
		return nil, err
	}
	return &tx, nil
}

// ParseWebhook verifies signature against body and decodes the event.
// Supply your live signing secret; test deliveries are verified against
// WebhookTestSecret automatically. Extra secrets (for example the previous
// secret during rotation) may be passed and are tried in order.
func ParseWebhook(body []byte, signature string, secret string, extraSecrets ...string) (*WebhookEvent, error) {
	var event WebhookEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return nil, errors.New("fluidpay: webhook body is not valid JSON: " + err.Error())
	}

	secrets := make([]string, 0, 2+len(extraSecrets))
	if event.IsTest() {
		secrets = append(secrets, WebhookTestSecret)
	}
	secrets = append(secrets, secret)
	secrets = append(secrets, extraSecrets...)

	for _, s := range secrets {
		if VerifyWebhookSignature(s, body, signature) {
			return &event, nil
		}
	}
	return nil, ErrInvalidWebhookSignature
}

// ParseWebhookRequest reads and verifies an incoming webhook HTTP request.
// It checks the Signature header, then Signature-Previous, against secret
// and extraSecrets. The event's EventID and ReplayOf are filled from the
// headers. The request body is consumed.
func ParseWebhookRequest(r *http.Request, secret string, extraSecrets ...string) (*WebhookEvent, error) {
	body, err := io.ReadAll(io.LimitReader(r.Body, maxWebhookBody))
	if err != nil {
		return nil, err
	}
	event, err := ParseWebhook(body, r.Header.Get(WebhookSignatureHeader), secret, extraSecrets...)
	if err != nil && errors.Is(err, ErrInvalidWebhookSignature) {
		if prev := r.Header.Get(WebhookSignaturePreviousHeader); prev != "" {
			event, err = ParseWebhook(body, prev, secret, extraSecrets...)
		}
	}
	if err != nil {
		return nil, err
	}
	event.EventID = r.Header.Get(WebhookEventIDHeader)
	event.ReplayOf = r.Header.Get(WebhookReplayOfHeader)
	return event, nil
}
