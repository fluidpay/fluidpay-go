package fluidpay

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

// Billing frequencies for plans and subscriptions.
const (
	BillingMonthly      = "monthly"
	BillingTwiceMonthly = "twice_monthly"
	BillingDaily        = "daily"
)

// Subscription statuses. Active, Failing, Failed, Error and Completed are
// set by the billing engine; Paused, PastDue and Cancelled are set through
// the lifecycle methods on SubscriptionsService.
const (
	SubscriptionActive    = "active"
	SubscriptionFailing   = "failing"
	SubscriptionFailed    = "failed"
	SubscriptionError     = "error"
	SubscriptionCompleted = "completed"
	SubscriptionPaused    = "paused"
	SubscriptionPastDue   = "past_due"
	SubscriptionCancelled = "cancelled"
	SubscriptionStopped   = "stopped"
)

// AddOnsService manages recurring add-ons: fixed or percentage amounts added
// to a subscription charge. Access it through Client.AddOns.
type AddOnsService struct{ client *Client }

// DiscountsService manages recurring discounts: fixed or percentage amounts
// subtracted from a subscription charge. Access it through Client.Discounts.
type DiscountsService struct{ client *Client }

// PlansService manages recurring plans. Access it through Client.Plans.
type PlansService struct{ client *Client }

// SubscriptionsService manages subscriptions and their lifecycle. Access it
// through Client.Subscriptions.
type SubscriptionsService struct{ client *Client }

// AdjustmentRequest creates or updates an add-on or discount. Provide either
// Amount (cents) or Percentage (thousandths of a percent: 43440 is
// 43.440%), never both. Duration is how many times it applies; 0 means until
// cancelled.
type AdjustmentRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Amount      int    `json:"amount,omitempty"`
	Percentage  int    `json:"percentage,omitempty"`
	Duration    int    `json:"duration"`
}

// AddOn is a recurring add-on.
type AddOn struct {
	APIResource

	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Amount      int       `json:"amount"`
	Percentage  int       `json:"percentage"`
	Duration    int       `json:"duration"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Discount is a recurring discount.
type Discount struct {
	APIResource

	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Amount      int       `json:"amount"`
	Percentage  int       `json:"percentage"`
	Duration    int       `json:"duration"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// AddOnList is a list of add-ons.
type AddOnList struct {
	APIResource
	Data       []*AddOn
	TotalCount int
}

// DiscountList is a list of discounts.
type DiscountList struct {
	APIResource
	Data       []*Discount
	TotalCount int
}

// PlanAdjustment references an add-on or discount from a plan or
// subscription. Only ID is required; the other fields override the
// referenced record's values for this plan or subscription.
type PlanAdjustment struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Amount      int    `json:"amount,omitempty"`
	Percentage  int    `json:"percentage,omitempty"`
	Duration    int    `json:"duration,omitempty"`
}

// PlanRequest creates or updates a plan.
type PlanRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	// Amount is the recurring charge in cents.
	Amount int `json:"amount"`
	// BillingCycleInterval runs the cycle every N months.
	BillingCycleInterval int `json:"billing_cycle_interval"`
	// BillingFrequency is one of the Billing* constants.
	BillingFrequency string `json:"billing_frequency"`
	// BillingDays is the day of month to bill ("1"), two comma separated
	// days for BillingTwiceMonthly ("1,15"), or "0" for the last day.
	BillingDays string `json:"billing_days"`
	// ChargeOnDay bills on the current day instead of BillingDays.
	ChargeOnDay bool `json:"charge_on_day,omitempty"`
	// Duration is the number of billings; 0 runs until cancelled.
	Duration  int              `json:"duration"`
	AddOns    []PlanAdjustment `json:"add_ons,omitempty"`
	Discounts []PlanAdjustment `json:"discounts,omitempty"`
}

// Plan is a recurring plan.
type Plan struct {
	APIResource

	ID                   string      `json:"id"`
	Name                 string      `json:"name"`
	Description          string      `json:"description"`
	Amount               int         `json:"amount"`
	BillingCycleInterval int         `json:"billing_cycle_interval"`
	BillingFrequency     string      `json:"billing_frequency"`
	BillingDays          string      `json:"billing_days"`
	ChargeOnDay          bool        `json:"charge_on_day"`
	TotalAddOns          int         `json:"total_add_ons"`
	TotalDiscounts       int         `json:"total_discounts"`
	Duration             int         `json:"duration"`
	AddOns               []*AddOn    `json:"add_ons"`
	Discounts            []*Discount `json:"discounts"`
	CreatedAt            time.Time   `json:"created_at"`
	UpdatedAt            time.Time   `json:"updated_at"`
}

// PlanList is a list of plans.
type PlanList struct {
	APIResource
	Data       []*Plan
	TotalCount int
}

// SubscriptionCustomer selects the vault customer and, optionally, which of
// their stored payment methods and addresses a subscription charges.
type SubscriptionCustomer struct {
	ID                string `json:"id"`
	PaymentMethodType string `json:"payment_method_type,omitempty"`
	PaymentMethodID   string `json:"payment_method_id,omitempty"`
	BillingAddressID  string `json:"billing_address_id,omitempty"`
	ShippingAddressID string `json:"shipping_address_id,omitempty"`
}

// SubscriptionRequest creates or updates a subscription. Billing fields
// mirror PlanRequest and override the plan's values.
type SubscriptionRequest struct {
	PlanID      string               `json:"plan_id,omitempty"`
	Description string               `json:"description,omitempty"`
	Customer    SubscriptionCustomer `json:"customer"`
	// Amount in cents. Optional when DeriveAmountFromLineItems is true.
	Amount               int    `json:"amount,omitempty"`
	BillingCycleInterval int    `json:"billing_cycle_interval,omitempty"`
	BillingFrequency     string `json:"billing_frequency,omitempty"`
	BillingDays          string `json:"billing_days,omitempty"`
	ChargeOnDay          bool   `json:"charge_on_day,omitempty"`
	Duration             int    `json:"duration,omitempty"`
	// NextBillDate is "YYYY-MM-DD".
	NextBillDate string           `json:"next_bill_date,omitempty"`
	AddOns       []PlanAdjustment `json:"add_ons,omitempty"`
	Discounts    []PlanAdjustment `json:"discounts,omitempty"`

	// Level 3 data forwarded on every recurring charge.
	LineItems                     []LineItem `json:"line_items,omitempty"`
	DeriveAmountFromLineItems     bool       `json:"derive_amount_from_line_items,omitempty"`
	SummaryCommodityCode          string     `json:"summary_commodity_code,omitempty"`
	ShipFromPostalCode            string     `json:"ship_from_postal_code,omitempty"`
	NationalTaxAmount             int        `json:"national_tax_amount,omitempty"`
	DutyAmount                    int        `json:"duty_amount,omitempty"`
	MerchantVATRegistrationNumber string     `json:"merchant_vat_registration_number,omitempty"`
	CustomerVATRegistrationNumber string     `json:"customer_vat_registration_number,omitempty"`
	PoNumber                      string     `json:"po_number,omitempty"`
	TaxAmount                     int        `json:"tax_amount,omitempty"`
	TaxExempt                     bool       `json:"tax_exempt,omitempty"`
}

// Subscription applies a plan to a customer.
type Subscription struct {
	APIResource

	ID          string `json:"id"`
	PlanID      string `json:"plan_id"`
	Status      string `json:"status"`
	Description string `json:"description"`

	Customer SubscriptionCustomer `json:"customer"`

	Amount               int    `json:"amount"`
	TotalAdds            int    `json:"total_adds"`
	TotalDiscounts       int    `json:"total_discounts"`
	BillingCycleInterval int    `json:"billing_cycle_interval"`
	BillingFrequency     string `json:"billing_frequency"`
	BillingDays          string `json:"billing_days"`
	ChargeOnDay          bool   `json:"charge_on_day"`
	Duration             int    `json:"duration"`
	// NextBillDate is "YYYY-MM-DD".
	NextBillDate string `json:"next_bill_date"`

	AddOns    []*AddOn    `json:"add_ons"`
	Discounts []*Discount `json:"discounts"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NeedsAttention reports whether billing has failed or errored. Monitor for
// all three statuses, not just SubscriptionError.
func (s *Subscription) NeedsAttention() bool {
	switch s.Status {
	case SubscriptionFailing, SubscriptionFailed, SubscriptionError:
		return true
	}
	return false
}

// SubscriptionList is a page of subscriptions.
type SubscriptionList struct {
	APIResource
	Data       []*Subscription
	TotalCount int
}

// SubscriptionSearchRequest filters subscriptions.
type SubscriptionSearchRequest struct {
	PlanID *StringQuery `json:"plan_id,omitempty"`
	// Customer filters on the customer ID. The gateway documents the field
	// as customer.id; it is sent nested to mirror the response shape.
	Customer       *SubscriptionCustomerQuery `json:"customer,omitempty"`
	NextBillDate   *DateRangeQuery            `json:"next_bill_date,omitempty"`
	ExpirationDate *DateRangeQuery            `json:"expiration_date,omitempty"`
	Status         *StringQuery               `json:"status,omitempty"`
	Limit          int                        `json:"limit,omitempty"`
	Offset         int                        `json:"offset,omitempty"`
}

// SubscriptionCustomerQuery filters subscriptions by customer.
type SubscriptionCustomerQuery struct {
	ID *StringQuery `json:"id,omitempty"`
}

// --- Add-ons ---

// Create creates an add-on.
func (s *AddOnsService) Create(ctx context.Context, req *AdjustmentRequest) (*AddOn, error) {
	if req == nil {
		return nil, errNilRequest("add-on")
	}
	return doOne[AddOn](ctx, s.client, http.MethodPost, "recurring/addon", nil, req)
}

// Get retrieves an add-on.
func (s *AddOnsService) Get(ctx context.Context, addOnID string) (*AddOn, error) {
	if err := requireID("add-on id", addOnID); err != nil {
		return nil, err
	}
	return doOne[AddOn](ctx, s.client, http.MethodGet, joinPath("recurring", "addon", addOnID), nil, nil)
}

// List retrieves all add-ons.
func (s *AddOnsService) List(ctx context.Context) (*AddOnList, error) {
	items, resp, err := doList[AddOn](ctx, s.client, http.MethodGet, "recurring/addons", nil, nil)
	if err != nil {
		return nil, err
	}
	return &AddOnList{APIResource: APIResource{LastResponse: resp}, Data: items, TotalCount: resp.TotalCount}, nil
}

// Update edits an add-on.
func (s *AddOnsService) Update(ctx context.Context, addOnID string, req *AdjustmentRequest) (*AddOn, error) {
	if err := requireID("add-on id", addOnID); err != nil {
		return nil, err
	}
	if req == nil {
		return nil, errNilRequest("add-on")
	}
	return doOne[AddOn](ctx, s.client, http.MethodPost, joinPath("recurring", "addon", addOnID), nil, req)
}

// Delete removes an add-on.
func (s *AddOnsService) Delete(ctx context.Context, addOnID string) error {
	if err := requireID("add-on id", addOnID); err != nil {
		return err
	}
	return doEmpty(ctx, s.client, http.MethodDelete, joinPath("recurring", "addon", addOnID), nil, nil)
}

// --- Discounts ---

// Create creates a discount.
func (s *DiscountsService) Create(ctx context.Context, req *AdjustmentRequest) (*Discount, error) {
	if req == nil {
		return nil, errNilRequest("discount")
	}
	return doOne[Discount](ctx, s.client, http.MethodPost, "recurring/discount", nil, req)
}

// Get retrieves a discount.
func (s *DiscountsService) Get(ctx context.Context, discountID string) (*Discount, error) {
	if err := requireID("discount id", discountID); err != nil {
		return nil, err
	}
	return doOne[Discount](ctx, s.client, http.MethodGet, joinPath("recurring", "discount", discountID), nil, nil)
}

// List retrieves all discounts.
func (s *DiscountsService) List(ctx context.Context) (*DiscountList, error) {
	items, resp, err := doList[Discount](ctx, s.client, http.MethodGet, "recurring/discounts", nil, nil)
	if err != nil {
		return nil, err
	}
	return &DiscountList{APIResource: APIResource{LastResponse: resp}, Data: items, TotalCount: resp.TotalCount}, nil
}

// Update edits a discount.
func (s *DiscountsService) Update(ctx context.Context, discountID string, req *AdjustmentRequest) (*Discount, error) {
	if err := requireID("discount id", discountID); err != nil {
		return nil, err
	}
	if req == nil {
		return nil, errNilRequest("discount")
	}
	return doOne[Discount](ctx, s.client, http.MethodPost, joinPath("recurring", "discount", discountID), nil, req)
}

// Delete removes a discount.
func (s *DiscountsService) Delete(ctx context.Context, discountID string) error {
	if err := requireID("discount id", discountID); err != nil {
		return err
	}
	return doEmpty(ctx, s.client, http.MethodDelete, joinPath("recurring", "discount", discountID), nil, nil)
}

// --- Plans ---

// Create creates a plan.
func (s *PlansService) Create(ctx context.Context, req *PlanRequest) (*Plan, error) {
	if req == nil {
		return nil, errNilRequest("plan")
	}
	return doOne[Plan](ctx, s.client, http.MethodPost, "recurring/plan", nil, req)
}

// Get retrieves a plan.
func (s *PlansService) Get(ctx context.Context, planID string) (*Plan, error) {
	if err := requireID("plan id", planID); err != nil {
		return nil, err
	}
	return doOne[Plan](ctx, s.client, http.MethodGet, joinPath("recurring", "plan", planID), nil, nil)
}

// List retrieves all plans.
func (s *PlansService) List(ctx context.Context) (*PlanList, error) {
	items, resp, err := doList[Plan](ctx, s.client, http.MethodGet, "recurring/plans", nil, nil)
	if err != nil {
		return nil, err
	}
	return &PlanList{APIResource: APIResource{LastResponse: resp}, Data: items, TotalCount: resp.TotalCount}, nil
}

// Update edits a plan.
func (s *PlansService) Update(ctx context.Context, planID string, req *PlanRequest) (*Plan, error) {
	if err := requireID("plan id", planID); err != nil {
		return nil, err
	}
	if req == nil {
		return nil, errNilRequest("plan")
	}
	return doOne[Plan](ctx, s.client, http.MethodPost, joinPath("recurring", "plan", planID), nil, req)
}

// Delete removes a plan.
func (s *PlansService) Delete(ctx context.Context, planID string) error {
	if err := requireID("plan id", planID); err != nil {
		return err
	}
	return doEmpty(ctx, s.client, http.MethodDelete, joinPath("recurring", "plan", planID), nil, nil)
}

// --- Subscriptions ---

// Create subscribes a customer to a plan.
func (s *SubscriptionsService) Create(ctx context.Context, req *SubscriptionRequest) (*Subscription, error) {
	if req == nil {
		return nil, errNilRequest("subscription")
	}
	return doOne[Subscription](ctx, s.client, http.MethodPost, "recurring/subscription", nil, req)
}

// Get retrieves a subscription.
func (s *SubscriptionsService) Get(ctx context.Context, subscriptionID string) (*Subscription, error) {
	if err := requireID("subscription id", subscriptionID); err != nil {
		return nil, err
	}
	return doOne[Subscription](ctx, s.client, http.MethodGet, joinPath("recurring", "subscription", subscriptionID), nil, nil)
}

// Search returns subscriptions matching req.
func (s *SubscriptionsService) Search(ctx context.Context, req *SubscriptionSearchRequest) (*SubscriptionList, error) {
	if req == nil {
		req = &SubscriptionSearchRequest{}
	}
	items, resp, err := doList[Subscription](ctx, s.client, http.MethodPost, "recurring/subscription/search", nil, req)
	if err != nil {
		return nil, err
	}
	return &SubscriptionList{APIResource: APIResource{LastResponse: resp}, Data: items, TotalCount: resp.TotalCount}, nil
}

// Update edits a subscription.
func (s *SubscriptionsService) Update(ctx context.Context, subscriptionID string, req *SubscriptionRequest) (*Subscription, error) {
	if err := requireID("subscription id", subscriptionID); err != nil {
		return nil, err
	}
	if req == nil {
		return nil, errNilRequest("subscription")
	}
	return doOne[Subscription](ctx, s.client, http.MethodPost, joinPath("recurring", "subscription", subscriptionID), nil, req)
}

// Delete removes a subscription.
func (s *SubscriptionsService) Delete(ctx context.Context, subscriptionID string) error {
	if err := requireID("subscription id", subscriptionID); err != nil {
		return err
	}
	return doEmpty(ctx, s.client, http.MethodDelete, joinPath("recurring", "subscription", subscriptionID), nil, nil)
}

// Pause suspends billing until Activate is called.
func (s *SubscriptionsService) Pause(ctx context.Context, subscriptionID string) (*Subscription, error) {
	return s.setStatus(ctx, subscriptionID, SubscriptionPaused, nil)
}

// MarkPastDue flags the subscription as past due. The gateway never sets
// this automatically.
func (s *SubscriptionsService) MarkPastDue(ctx context.Context, subscriptionID string) (*Subscription, error) {
	return s.setStatus(ctx, subscriptionID, SubscriptionPastDue, nil)
}

// Cancel permanently stops billing.
func (s *SubscriptionsService) Cancel(ctx context.Context, subscriptionID string) (*Subscription, error) {
	return s.setStatus(ctx, subscriptionID, SubscriptionCancelled, nil)
}

// Complete marks the subscription as finished.
func (s *SubscriptionsService) Complete(ctx context.Context, subscriptionID string) (*Subscription, error) {
	return s.setStatus(ctx, subscriptionID, SubscriptionCompleted, nil)
}

// Activate resumes billing. Pass a non-zero nextBillDate to choose when the
// next charge runs; the zero time keeps the gateway default.
func (s *SubscriptionsService) Activate(ctx context.Context, subscriptionID string, nextBillDate time.Time) (*Subscription, error) {
	var q url.Values
	if !nextBillDate.IsZero() {
		q = url.Values{"next_bill_date": []string{nextBillDate.Format("2006-01-02")}}
	}
	return s.setStatus(ctx, subscriptionID, SubscriptionActive, q)
}

func (s *SubscriptionsService) setStatus(ctx context.Context, subscriptionID, status string, q url.Values) (*Subscription, error) {
	if err := requireID("subscription id", subscriptionID); err != nil {
		return nil, err
	}
	return doOne[Subscription](ctx, s.client, http.MethodGet, joinPath("recurring", "subscription", subscriptionID, "status", status), q, nil)
}
