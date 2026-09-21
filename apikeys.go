package fluidpay

import (
	"context"
	"net/http"
	"time"
)

// APIKeysService creates, lists and deletes API keys belonging to the
// authenticated user. Access it through Client.APIKeys.
//
// Private keys ("api_") authenticate server-side calls such as this SDK's;
// public keys ("pub_") are for browser-side tooling such as the Tokenizer.
// A private key is only returned once, at creation: store it securely.
type APIKeysService struct{ client *Client }

// API key types.
const (
	APIKeyTypePrivate = "api"
	APIKeyTypePublic  = "pub"
)

// APIKeyCreateRequest creates an API key.
type APIKeyCreateRequest struct {
	// Type is APIKeyTypePrivate or APIKeyTypePublic.
	Type string `json:"type"`
	// Name labels the key in the control panel.
	Name string `json:"name"`
}

// APIKey is an API key. APIKey.Key is only populated on creation.
type APIKey struct {
	APIResource

	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Type   string `json:"type"`
	Name   string `json:"name"`
	// Key is the secret value, returned once when the key is created.
	Key           string    `json:"api_key"`
	AccountType   string    `json:"account_type"`
	AccountTypeID string    `json:"account_type_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// APIKeyList is a list of API keys.
type APIKeyList struct {
	APIResource
	Data       []*APIKey
	TotalCount int
}

// Create creates an API key. The secret is in APIKey.Key and is never
// returned again.
func (s *APIKeysService) Create(ctx context.Context, req *APIKeyCreateRequest) (*APIKey, error) {
	if req == nil {
		return nil, errNilRequest("api key")
	}
	return doOne[APIKey](ctx, s.client, http.MethodPost, "user/apikey", nil, req)
}

// List retrieves the current user's API keys.
func (s *APIKeysService) List(ctx context.Context) (*APIKeyList, error) {
	items, resp, err := doList[APIKey](ctx, s.client, http.MethodGet, "user/apikeys", nil, nil)
	if err != nil {
		return nil, err
	}
	return &APIKeyList{APIResource: APIResource{LastResponse: resp}, Data: items, TotalCount: resp.TotalCount}, nil
}

// Delete revokes an API key. Requests made with a revoked key fail with an
// unauthorized error.
func (s *APIKeysService) Delete(ctx context.Context, apiKeyID string) error {
	if err := requireID("api key id", apiKeyID); err != nil {
		return err
	}
	return doEmpty(ctx, s.client, http.MethodDelete, joinPath("user", "apikey", apiKeyID), nil, nil)
}
