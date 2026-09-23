package fluidpay

import (
	"context"
	"errors"
	"net/http"
)

// LookupService answers questions about payment instruments. Access it
// through Client.Lookup.
type LookupService struct{ client *Client }

// BINLookupRequest looks up card metadata by BIN. Provide exactly one of
// BIN, TempToken or CustomerID. Country and State together determine
// surchargeability.
type BINLookupRequest struct {
	// BIN is the first six or more digits of the card number.
	BIN string `json:"bin,omitempty"`
	// TempToken is a Tokenizer token; the BIN is derived from it.
	TempToken string `json:"temp_token,omitempty"`
	// CustomerID looks up the customer's default card, or the card named
	// by PaymentMethodID.
	CustomerID      string `json:"customer_id,omitempty"`
	PaymentMethodID string `json:"payment_method_id,omitempty"`
	Country         string `json:"country,omitempty"`
	State           string `json:"state,omitempty"`
}

// BINLookup is the metadata for a BIN. Unknown BINs return an empty result
// rather than an error.
type BINLookup struct {
	APIResource

	BIN              string `json:"bin"`
	CardBrand        string `json:"card_brand"`
	IssuingBank      string `json:"issuing_bank"`
	CardType         string `json:"card_type"`
	CardLevelGeneric string `json:"card_level_generic"`
	Country          string `json:"country"`
	// IsSurchargeable is only meaningful when Country and State were
	// supplied in the request.
	IsSurchargeable   bool   `json:"is_surchargeable"`
	PaymentMethodType string `json:"payment_method_type"`
}

// BIN looks up card metadata and surchargeability.
func (s *LookupService) BIN(ctx context.Context, req *BINLookupRequest) (*BINLookup, error) {
	if req == nil {
		return nil, errNilRequest("bin lookup")
	}
	if req.BIN == "" && req.TempToken == "" && req.CustomerID == "" {
		return nil, errors.New("fluidpay: bin lookup requires one of BIN, TempToken or CustomerID")
	}
	return doOne[BINLookup](ctx, s.client, http.MethodPost, "lookup/bin/protected", nil, req)
}
