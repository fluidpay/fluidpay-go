package fluidpay

import (
	"context"
	"net/http"
	"time"
)

// SettlementsService searches settlement batches. Access it through
// Client.Settlements.
type SettlementsService struct{ client *Client }

// SettlementBatchSearchRequest filters settlement batches.
type SettlementBatchSearchRequest struct {
	BatchDate         *DateRangeQuery `json:"batch_date,omitempty"`
	SettlementBatchID *StringQuery    `json:"settlement_batch_id,omitempty"`
	// Limit is the page size (1-100, default 10). Offset skips records.
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

// SettlementBatchSearchResult is the result of a settlement batch search.
type SettlementBatchSearchResult struct {
	APIResource

	// Summary aggregates the batches per merchant, day and processor.
	Summary []*SettlementSummary `json:"summary"`
	// Results lists the matching batches.
	Results []*SettlementBatch `json:"results"`
	// TotalCount is the number of matching batches across all pages.
	TotalCount int `json:"-"`
}

// SettlementSummary aggregates settled activity for one day and processor.
type SettlementSummary struct {
	MerchantID      string `json:"merchant_id"`
	BatchDate       string `json:"batch_date"`
	ProcessorID     string `json:"processor_id"`
	ProcessorName   string `json:"processor_name"`
	NumTransactions int    `json:"num_transactions"`
	Captured        int    `json:"captured"`
	Credit          int    `json:"credit"`
}

// SettlementBatch is one settled batch.
type SettlementBatch struct {
	ID              string    `json:"id"`
	MerchantID      string    `json:"merchant_id"`
	BatchDate       time.Time `json:"batch_date"`
	ProcessorID     string    `json:"processor_id"`
	ProcessorName   string    `json:"processor_name"`
	ProcessorType   string    `json:"processor_type"`
	BatchNumber     int       `json:"batch_number"`
	NumTransactions int       `json:"num_transactions"`
	AmountCaptured  int       `json:"amount_captured"`
	AmountCredit    int       `json:"amount_credit"`
	NetDeposit      int       `json:"net_deposit"`
	ResponseCode    int       `json:"response_code"`
	ResponseMessage string    `json:"response_message"`
}

// SearchBatches returns settlement batches matching req.
func (s *SettlementsService) SearchBatches(ctx context.Context, req *SettlementBatchSearchRequest) (*SettlementBatchSearchResult, error) {
	if req == nil {
		req = &SettlementBatchSearchRequest{}
	}
	out, err := doOne[SettlementBatchSearchResult](ctx, s.client, http.MethodPost, "settlement/batch/search", nil, req)
	if err != nil {
		return nil, err
	}
	if out.LastResponse != nil {
		out.TotalCount = out.LastResponse.TotalCount
	}
	return out, nil
}
