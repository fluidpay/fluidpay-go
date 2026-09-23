package fluidpay

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Values from the webhook documentation's worked example.
const (
	docWebhookSecret    = "12345678-1234-1234-1234-123456789012"
	docWebhookBody      = `{"data":"this is test data"}`
	docWebhookSignature = "JacUiw_ztpEZJWvOhhKoHTLBf4b-aZv9n_0YmJJxltc"
)

func TestComputeWebhookSignature_MatchesDocs(t *testing.T) {
	equal(t, "signature", ComputeWebhookSignature(docWebhookSecret, []byte(docWebhookBody)), docWebhookSignature)
}

func TestVerifyWebhookSignature(t *testing.T) {
	body := []byte(docWebhookBody)
	equal(t, "valid", VerifyWebhookSignature(docWebhookSecret, body, docWebhookSignature), true)
	equal(t, "wrong secret", VerifyWebhookSignature("other", body, docWebhookSignature), false)
	equal(t, "tampered body", VerifyWebhookSignature(docWebhookSecret, []byte(docWebhookBody+"\n"), docWebhookSignature), false)
	equal(t, "tampered signature", VerifyWebhookSignature(docWebhookSecret, body, "x"+docWebhookSignature[1:]), false)
	equal(t, "empty signature", VerifyWebhookSignature(docWebhookSecret, body, ""), false)
	equal(t, "empty secret", VerifyWebhookSignature("", body, docWebhookSignature), false)
}

func TestParseWebhook(t *testing.T) {
	body := []byte(webhookTransactionFixture)
	secret := "live-secret"
	sig := ComputeWebhookSignature(secret, body)

	event, err := ParseWebhook(body, sig, secret)
	mustNoError(t, err)
	equal(t, "type", event.Type, WebhookTypeTransactionCreate)
	equal(t, "account", event.AccountTypeID, "testmerchant12345678")
	equal(t, "transaction id", event.TransactionID, "bm5s8gm9ku6ejcu15t9g")
	equal(t, "action_at", event.ActionAt.Year(), 2019)
	equal(t, "IsTest", event.IsTest(), false)

	tx, err := event.Transaction()
	mustNoError(t, err)
	equal(t, "tx id", tx.ID, "bm5s8gm9ku6ejcu15t9g")
	equal(t, "approved", tx.Approved(), true)
	equal(t, "auth code", tx.ResponseBody.Card.AuthCode, "TAS000")

	// Wrong secret.
	_, err = ParseWebhook(body, sig, "other")
	if !errors.Is(err, ErrInvalidWebhookSignature) {
		t.Fatalf("expected ErrInvalidWebhookSignature, got %v", err)
	}

	// Previous secret still accepted through extraSecrets.
	_, err = ParseWebhook(body, sig, "rotated", secret)
	mustNoError(t, err)

	// Malformed JSON.
	_, err = ParseWebhook([]byte("{"), sig, secret)
	mustError(t, err, "not valid JSON")

	// Test deliveries are signed with the shared test key.
	testBody := []byte(`{"status":"success","msg":"success","type":"test","account_type":"merchant","account_type_id":"m","action_at":"2026-08-03T15:04:05Z"}` + "\n")
	event, err = ParseWebhook(testBody, ComputeWebhookSignature(WebhookTestSecret, testBody), secret)
	mustNoError(t, err)
	equal(t, "IsTest", event.IsTest(), true)
	_, err = event.Transaction()
	mustError(t, err, "has no data")

	// A test-typed body signed with neither key is rejected.
	_, err = ParseWebhook(testBody, ComputeWebhookSignature("bogus", testBody), secret)
	if !errors.Is(err, ErrInvalidWebhookSignature) {
		t.Fatalf("expected ErrInvalidWebhookSignature, got %v", err)
	}
}

func TestParseWebhookRequest(t *testing.T) {
	secret := "live-secret"
	body := webhookTransactionFixture
	newReq := func(sig, prev string) *http.Request {
		r := httptest.NewRequest(http.MethodPost, "/hooks", strings.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		if sig != "" {
			r.Header.Set(WebhookSignatureHeader, sig)
		}
		if prev != "" {
			r.Header.Set(WebhookSignaturePreviousHeader, prev)
		}
		r.Header.Set(WebhookEventIDHeader, "evt_1")
		r.Header.Set(WebhookReplayOfHeader, "evt_0")
		return r
	}

	event, err := ParseWebhookRequest(newReq(ComputeWebhookSignature(secret, []byte(body)), ""), secret)
	mustNoError(t, err)
	equal(t, "EventID", event.EventID, "evt_1")
	equal(t, "ReplayOf", event.ReplayOf, "evt_0")
	equal(t, "type", event.Type, WebhookTypeTransactionCreate)

	// During rotation the Signature header may be made with the new secret
	// we do not have yet, while Signature-Previous matches ours.
	event, err = ParseWebhookRequest(newReq("bad", ComputeWebhookSignature(secret, []byte(body))), secret)
	mustNoError(t, err)
	equal(t, "EventID after fallback", event.EventID, "evt_1")

	_, err = ParseWebhookRequest(newReq("bad", "also-bad"), secret)
	if !errors.Is(err, ErrInvalidWebhookSignature) {
		t.Fatalf("expected ErrInvalidWebhookSignature, got %v", err)
	}
	_, err = ParseWebhookRequest(newReq("", ""), secret)
	if !errors.Is(err, ErrInvalidWebhookSignature) {
		t.Fatalf("expected ErrInvalidWebhookSignature for missing header, got %v", err)
	}
}

func TestParseWebhookRequest_HandlerIntegration(t *testing.T) {
	secret := "live-secret"
	var seen *WebhookEvent
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		event, err := ParseWebhookRequest(r, secret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		seen = event
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/hooks", strings.NewReader(webhookTransactionFixture))
	req.Header.Set(WebhookSignatureHeader, ComputeWebhookSignature(secret, []byte(webhookTransactionFixture)))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	equal(t, "status", rec.Code, http.StatusOK)
	if seen == nil || seen.TransactionID != "bm5s8gm9ku6ejcu15t9g" {
		t.Fatalf("event not captured: %+v", seen)
	}

	req = httptest.NewRequest(http.MethodPost, "/hooks", strings.NewReader(webhookTransactionFixture))
	req.Header.Set(WebhookSignatureHeader, "forged")
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	equal(t, "forged status", rec.Code, http.StatusUnauthorized)
}
