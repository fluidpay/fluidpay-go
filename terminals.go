package fluidpay

import (
	"context"
	"net/http"
	"time"
)

// TerminalsService lists physical terminals and settles their batches.
// Transactions are sent to a terminal through Transactions.Sale with
// PaymentMethod.Terminal set. Access it through Client.Terminals.
type TerminalsService struct{ client *Client }

// Terminal is a physical payment terminal registered on the account.
type Terminal struct {
	APIResource

	ID           string    `json:"id"`
	MerchantID   string    `json:"merchant_id"`
	Manufacturer string    `json:"manufacturer"`
	Model        string    `json:"model"`
	SerialNumber string    `json:"serial_number"`
	TPN          string    `json:"tpn"`
	Description  string    `json:"description"`
	Status       string    `json:"status"`
	AuthKey      string    `json:"auth_key"`
	RegisterID   string    `json:"register_id"`
	AutoSettle   bool      `json:"auto_settle"`
	SettleAt     string    `json:"settle_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Active reports whether the terminal is enabled.
func (t *Terminal) Active() bool { return t.Status == "active" }

// TerminalList is a list of terminals.
type TerminalList struct {
	APIResource
	Data       []*Terminal
	TotalCount int
}

// List retrieves every terminal on the account, including disabled ones.
func (s *TerminalsService) List(ctx context.Context) (*TerminalList, error) {
	items, resp, err := doList[Terminal](ctx, s.client, http.MethodGet, "terminals", nil, nil)
	if err != nil {
		return nil, err
	}
	return &TerminalList{APIResource: APIResource{LastResponse: resp}, Data: items, TotalCount: resp.TotalCount}, nil
}

// Settle closes the open batch on a terminal.
func (s *TerminalsService) Settle(ctx context.Context, terminalID string) (*APIResponse, error) {
	if err := requireID("terminal id", terminalID); err != nil {
		return nil, err
	}
	return doEmpty(ctx, s.client, http.MethodPost, joinPath("terminal", terminalID, "settle"), nil, nil)
}
