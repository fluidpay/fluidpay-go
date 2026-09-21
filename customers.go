package fluidpay

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

// CustomersService is the Customer Vault: the canonical way to store
// customers, their addresses and their payment methods for later charging.
// Access it through Client.Customers.
//
// Card verification data (CVC/CVV) is never stored, to keep the vault PCI
// compliant. Some operations run a small verification transaction when
// VerificationOptions ask for it.
type CustomersService struct{ client *Client }

// Customer flags.
const (
	CustomerFlagSurchargeExempt = "surcharge_exempt"
)

// Payment method types stored in the vault.
const (
	PaymentMethodTypeCard = "card"
	PaymentMethodTypeACH  = "ach"
)

// VerificationOptions control whether the gateway verifies a payment method
// while storing it. They are sent as query parameters, not in the body.
type VerificationOptions struct {
	// Validate runs a $0 verification on a card.
	Validate bool
	// Authorize runs a $1.00 authorization on a card.
	Authorize bool
	// BypassRuleEngine skips rule engine checks for the verification.
	BypassRuleEngine bool
}

func (o VerificationOptions) query() url.Values {
	q := url.Values{}
	if o.Validate {
		q.Set("validate", "true")
	}
	if o.Authorize {
		q.Set("authorize", "true")
	}
	if o.BypassRuleEngine {
		q.Set("bypass_rule_engine", "true")
	}
	if len(q) == 0 {
		return nil
	}
	return q
}

// VaultAddress is an address to store on a customer. Note that the vault
// uses line_1 and line_2 whereas transactions use address_line_1 and
// address_line_2.
type VaultAddress struct {
	FirstName  string `json:"first_name,omitempty"`
	LastName   string `json:"last_name,omitempty"`
	Company    string `json:"company,omitempty"`
	Line1      string `json:"line_1,omitempty"`
	Line2      string `json:"line_2,omitempty"`
	City       string `json:"city,omitempty"`
	State      string `json:"state,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
	Country    string `json:"country,omitempty"`
	Email      string `json:"email,omitempty"`
	Phone      string `json:"phone,omitempty"`
	Fax        string `json:"fax,omitempty"`
}

// VaultCard stores a card. ID is optional (2-40 characters: letters,
// digits, ".", "-", "_" and "^"); the gateway generates one otherwise.
type VaultCard struct {
	VerificationOptions `json:"-"`

	ID             string `json:"id,omitempty"`
	Number         string `json:"number"`
	ExpirationDate string `json:"expiration_date"`
	// Flags is honoured on update only.
	Flags []string `json:"flags,omitempty"`
}

// VaultACH stores a bank account.
type VaultACH struct {
	VerificationOptions `json:"-"`

	ID            string `json:"id,omitempty"`
	AccountNumber string `json:"account_number"`
	RoutingNumber string `json:"routing_number"`
	// AccountType is AccountTypeChecking or AccountTypeSavings.
	AccountType string `json:"account_type"`
	// SecCode is SecCodeWeb, SecCodeCCD or SecCodePPD.
	SecCode string `json:"sec_code"`
}

// VaultToken stores a card or bank account from a Tokenizer token.
type VaultToken struct {
	VerificationOptions `json:"-"`

	Token string `json:"token"`
}

// VaultApplePay stores a card from an Apple Pay token.
type VaultApplePay struct {
	VerificationOptions `json:"-"`

	KeyID          string          `json:"key_id,omitempty"`
	ProcessorID    string          `json:"processor_id,omitempty"`
	TemporaryToken string          `json:"temporary_token,omitempty"`
	PKPaymentToken json.RawMessage `json:"pkpaymenttoken,omitempty"`
}

// VaultPaymentMethod holds exactly one payment method to store as a
// customer's default.
type VaultPaymentMethod struct {
	Card           *VaultCard      `json:"card,omitempty"`
	ACH            *VaultACH       `json:"ach,omitempty"`
	Token          string          `json:"token,omitempty"`
	ApplePay       *VaultApplePay  `json:"apple_pay,omitempty"`
	GooglePayToken json.RawMessage `json:"google_pay_token,omitempty"`
}

// CustomerCreateRequest creates a vault customer. Everything is optional; a
// customer can be created empty and populated later.
type CustomerCreateRequest struct {
	VerificationOptions `json:"-"`

	// ID is an optional custom identifier (2-40 characters: letters,
	// digits, ".", "-", "_" and "^").
	ID string `json:"id,omitempty"`
	// IDFormat asks the gateway to generate an ID in a given format, for
	// example "xid_type_last4". It cannot be combined with ID.
	IDFormat    string   `json:"id_format,omitempty"`
	Description string   `json:"description,omitempty"`
	Flags       []string `json:"flags,omitempty"`

	DefaultPayment         *VaultPaymentMethod `json:"default_payment,omitempty"`
	DefaultBillingAddress  *VaultAddress       `json:"default_billing_address,omitempty"`
	DefaultShippingAddress *VaultAddress       `json:"default_shipping_address,omitempty"`
}

// CustomerUpdateRequest edits a vault customer.
type CustomerUpdateRequest struct {
	Description string   `json:"description,omitempty"`
	Notes       string   `json:"notes,omitempty"`
	Flags       []string `json:"flags,omitempty"`
	// Defaults selects the stored address and payment method used when a
	// transaction does not specify them.
	Defaults *CustomerDefaults `json:"defaults,omitempty"`
}

// CustomerDefaults identifies a customer's default address and payment
// method. PaymentMethodID and PaymentMethodType must be set together.
type CustomerDefaults struct {
	BillingAddressID  string `json:"billing_address_id,omitempty"`
	ShippingAddressID string `json:"shipping_address_id,omitempty"`
	PaymentMethodID   string `json:"payment_method_id,omitempty"`
	PaymentMethodType string `json:"payment_method_type,omitempty"`
}

// CustomerAddress is an address stored on a customer.
type CustomerAddress struct {
	VaultAddress
	ID string `json:"id"`
	// Hash identifies the address content; the gateway rejects duplicates.
	Hash string `json:"hash"`
}

// CustomerCard is a card stored on a customer. The full number is never
// returned.
type CustomerCard struct {
	ID             string   `json:"id"`
	CardType       string   `json:"card_type"`
	ExpirationDate string   `json:"expiration_date"`
	MaskedNumber   string   `json:"masked_number"`
	Flags          []string `json:"flags"`
	ProcessorID    string   `json:"processor_id"`
}

// CustomerACH is a bank account stored on a customer.
type CustomerACH struct {
	ID                  string   `json:"id"`
	AccountType         string   `json:"account_type"`
	SecCode             string   `json:"sec_code"`
	MaskedAccountNumber string   `json:"masked_account_number"`
	RoutingNumber       string   `json:"routing_number"`
	Flags               []string `json:"flags"`
	ProcessorID         string   `json:"processor_id"`
}

// Customer is a vault customer record. The gateway nests the interesting
// fields under data.customer; the SDK flattens them for convenience.
type Customer struct {
	APIResource

	ID          string   `json:"id"`
	OwnerID     string   `json:"owner_id"`
	Description string   `json:"description"`
	Notes       string   `json:"notes"`
	Flags       []string `json:"flags"`

	Addresses []*CustomerAddress `json:"addresses"`
	Defaults  CustomerDefaults   `json:"defaults"`
	Cards     []*CustomerCard    `json:"cards"`
	ACH       []*CustomerACH     `json:"ach"`

	// CreatedAddressID is set on the result of CreateAddress.
	CreatedAddressID string `json:"created_address_id,omitempty"`
	// CreatedPaymentMethodID is set on the results of CreateCard,
	// CreateACH, CreateToken, CreateApplePay and CreateGooglePay.
	CreatedPaymentMethodID string `json:"created_payment_method_id,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// customerWire mirrors the gateway's nested customer shape.
type customerWire struct {
	ID                     string `json:"id"`
	OwnerID                string `json:"owner_id"`
	CreatedAddressID       string `json:"created_address_id"`
	CreatedPaymentMethodID string `json:"created_payment_method_id"`
	Data                   struct {
		Customer struct {
			Addresses   []*CustomerAddress `json:"addresses"`
			Defaults    CustomerDefaults   `json:"defaults"`
			Description string             `json:"description"`
			Notes       string             `json:"notes"`
			Flags       []string           `json:"flags"`
			Payments    struct {
				ACH   []*CustomerACH  `json:"ach"`
				Cards []*CustomerCard `json:"cards"`
			} `json:"payments"`
		} `json:"customer"`
	} `json:"data"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UnmarshalJSON flattens the gateway's nested representation.
func (c *Customer) UnmarshalJSON(b []byte) error {
	var w customerWire
	if err := json.Unmarshal(b, &w); err != nil {
		return err
	}
	cust := w.Data.Customer
	*c = Customer{
		ID:                     w.ID,
		OwnerID:                w.OwnerID,
		Description:            cust.Description,
		Notes:                  cust.Notes,
		Flags:                  cust.Flags,
		Addresses:              cust.Addresses,
		Defaults:               cust.Defaults,
		Cards:                  cust.Payments.Cards,
		ACH:                    cust.Payments.ACH,
		CreatedAddressID:       w.CreatedAddressID,
		CreatedPaymentMethodID: w.CreatedPaymentMethodID,
		CreatedAt:              w.CreatedAt,
		UpdatedAt:              w.UpdatedAt,
	}
	return nil
}

// Address returns the stored address with the given ID, or nil.
func (c *Customer) Address(id string) *CustomerAddress {
	for _, a := range c.Addresses {
		if a != nil && a.ID == id {
			return a
		}
	}
	return nil
}

// Card returns the stored card with the given ID, or nil.
func (c *Customer) Card(id string) *CustomerCard {
	for _, card := range c.Cards {
		if card != nil && card.ID == id {
			return card
		}
	}
	return nil
}

// DefaultCard returns the customer's default card, or nil when the default
// payment method is not a card.
func (c *Customer) DefaultCard() *CustomerCard {
	if c.Defaults.PaymentMethodType != PaymentMethodTypeCard {
		return nil
	}
	return c.Card(c.Defaults.PaymentMethodID)
}

// CustomerList is a page of vault customers.
type CustomerList struct {
	APIResource
	Data       []*Customer
	TotalCount int
}

// CustomerSearchRequest filters vault customers. Every field is optional.
type CustomerSearchRequest struct {
	ID                *StringQuery    `json:"id,omitempty"`
	AddressID         *StringQuery    `json:"address_id,omitempty"`
	FirstName         *StringQuery    `json:"first_name,omitempty"`
	LastName          *StringQuery    `json:"last_name,omitempty"`
	Company           *StringQuery    `json:"company,omitempty"`
	AddressLine1      *StringQuery    `json:"address_line_1,omitempty"`
	AddressLine2      *StringQuery    `json:"address_line_2,omitempty"`
	City              *StringQuery    `json:"city,omitempty"`
	State             *StringQuery    `json:"state,omitempty"`
	PostalCode        *StringQuery    `json:"postal_code,omitempty"`
	Country           *StringQuery    `json:"country,omitempty"`
	Email             *StringQuery    `json:"email,omitempty"`
	Phone             *StringQuery    `json:"phone,omitempty"`
	Fax               *StringQuery    `json:"fax,omitempty"`
	PaymentMethodType *StringQuery    `json:"payment_method_type,omitempty"`
	BillingAddressID  *StringQuery    `json:"billing_address_id,omitempty"`
	ShippingAddressID *StringQuery    `json:"shipping_address_id,omitempty"`
	CreatedAt         *DateRangeQuery `json:"created_at,omitempty"`
	UpdatedAt         *DateRangeQuery `json:"updated_at,omitempty"`
	// Limit is the page size (1-100, default 10). Offset skips records.
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

// Create stores a new customer.
func (s *CustomersService) Create(ctx context.Context, req *CustomerCreateRequest) (*Customer, error) {
	if req == nil {
		req = &CustomerCreateRequest{}
	}
	return doOne[Customer](ctx, s.client, http.MethodPost, "vault/customer", req.query(), req)
}

// Get retrieves a customer with all stored addresses and payment methods.
func (s *CustomersService) Get(ctx context.Context, customerID string) (*Customer, error) {
	if err := requireID("customer id", customerID); err != nil {
		return nil, err
	}
	return doOne[Customer](ctx, s.client, http.MethodGet, joinPath("vault", customerID), nil, nil)
}

// Search returns customers matching req.
func (s *CustomersService) Search(ctx context.Context, req *CustomerSearchRequest) (*CustomerList, error) {
	if req == nil {
		req = &CustomerSearchRequest{}
	}
	items, resp, err := doList[Customer](ctx, s.client, http.MethodPost, "vault/customer/search", nil, req)
	if err != nil {
		return nil, err
	}
	return &CustomerList{APIResource: APIResource{LastResponse: resp}, Data: items, TotalCount: resp.TotalCount}, nil
}

// Update edits a customer's description, notes, flags or defaults.
func (s *CustomersService) Update(ctx context.Context, customerID string, req *CustomerUpdateRequest) (*Customer, error) {
	if err := requireID("customer id", customerID); err != nil {
		return nil, err
	}
	if req == nil {
		return nil, errNilRequest("customer update")
	}
	return doOne[Customer](ctx, s.client, http.MethodPost, joinPath("vault", "customer", customerID), nil, req)
}

// Delete removes a customer and everything stored on it.
func (s *CustomersService) Delete(ctx context.Context, customerID string) (*APIResponse, error) {
	if err := requireID("customer id", customerID); err != nil {
		return nil, err
	}
	return doEmpty(ctx, s.client, http.MethodDelete, joinPath("vault", customerID), nil, nil)
}

// CreateAddress stores an address on a customer. The new address ID is in
// Customer.CreatedAddressID. The gateway rejects an address identical to one
// already stored; check Customer.Addresses first and reuse it instead.
func (s *CustomersService) CreateAddress(ctx context.Context, customerID string, addr *VaultAddress) (*Customer, error) {
	if err := requireID("customer id", customerID); err != nil {
		return nil, err
	}
	if addr == nil {
		return nil, errNilRequest("address")
	}
	return doOne[Customer](ctx, s.client, http.MethodPost, joinPath("vault", "customer", customerID, "address"), nil, addr)
}

// UpdateAddress replaces a stored address.
func (s *CustomersService) UpdateAddress(ctx context.Context, customerID, addressID string, addr *VaultAddress) (*Customer, error) {
	if err := requireIDs("customer id", customerID, "address id", addressID); err != nil {
		return nil, err
	}
	if addr == nil {
		return nil, errNilRequest("address")
	}
	return doOne[Customer](ctx, s.client, http.MethodPost, joinPath("vault", "customer", customerID, "address", addressID), nil, addr)
}

// DeleteAddress removes a stored address.
func (s *CustomersService) DeleteAddress(ctx context.Context, customerID, addressID string) (*APIResponse, error) {
	if err := requireIDs("customer id", customerID, "address id", addressID); err != nil {
		return nil, err
	}
	return doEmpty(ctx, s.client, http.MethodDelete, joinPath("vault", "customer", customerID, "address", addressID), nil, nil)
}

// CreateCard stores a card. The new payment method ID is in
// Customer.CreatedPaymentMethodID.
func (s *CustomersService) CreateCard(ctx context.Context, customerID string, card *VaultCard) (*Customer, error) {
	if err := requireID("customer id", customerID); err != nil {
		return nil, err
	}
	if card == nil {
		return nil, errNilRequest("card")
	}
	return doOne[Customer](ctx, s.client, http.MethodPost, joinPath("vault", "customer", customerID, "card"), card.query(), card)
}

// CreateACH stores a bank account.
func (s *CustomersService) CreateACH(ctx context.Context, customerID string, ach *VaultACH) (*Customer, error) {
	if err := requireID("customer id", customerID); err != nil {
		return nil, err
	}
	if ach == nil {
		return nil, errNilRequest("ach")
	}
	return doOne[Customer](ctx, s.client, http.MethodPost, joinPath("vault", "customer", customerID, "ach"), ach.query(), ach)
}

// CreateToken stores the card or bank account behind a Tokenizer token.
func (s *CustomersService) CreateToken(ctx context.Context, customerID string, token *VaultToken) (*Customer, error) {
	if err := requireID("customer id", customerID); err != nil {
		return nil, err
	}
	if token == nil {
		return nil, errNilRequest("token")
	}
	return doOne[Customer](ctx, s.client, http.MethodPost, joinPath("vault", "customer", customerID, "token"), token.query(), token)
}

// CreateApplePay stores a card from an Apple Pay token.
func (s *CustomersService) CreateApplePay(ctx context.Context, customerID string, applePay *VaultApplePay) (*Customer, error) {
	if err := requireID("customer id", customerID); err != nil {
		return nil, err
	}
	if applePay == nil {
		return nil, errNilRequest("apple pay")
	}
	return doOne[Customer](ctx, s.client, http.MethodPost, joinPath("vault", "customer", customerID, "applepay"), applePay.query(), applePay)
}

// CreateGooglePay stores a card from a Google Pay token. Pass the token
// object exactly as the Google Pay client library produced it.
func (s *CustomersService) CreateGooglePay(ctx context.Context, customerID string, token json.RawMessage, opts VerificationOptions) (*Customer, error) {
	if err := requireID("customer id", customerID); err != nil {
		return nil, err
	}
	if len(token) == 0 {
		return nil, errNilRequest("google pay token")
	}
	return doOne[Customer](ctx, s.client, http.MethodPost, joinPath("vault", "customer", customerID, "googlepay"), opts.query(), token)
}

// UpdateCard replaces a stored card's number, expiration date or flags.
func (s *CustomersService) UpdateCard(ctx context.Context, customerID, paymentMethodID string, card *VaultCard) (*Customer, error) {
	if err := requireIDs("customer id", customerID, "payment method id", paymentMethodID); err != nil {
		return nil, err
	}
	if card == nil {
		return nil, errNilRequest("card")
	}
	return doOne[Customer](ctx, s.client, http.MethodPost, joinPath("vault", "customer", customerID, "card", paymentMethodID), card.query(), card)
}

// UpdateACH replaces a stored bank account.
func (s *CustomersService) UpdateACH(ctx context.Context, customerID, paymentMethodID string, ach *VaultACH) (*Customer, error) {
	if err := requireIDs("customer id", customerID, "payment method id", paymentMethodID); err != nil {
		return nil, err
	}
	if ach == nil {
		return nil, errNilRequest("ach")
	}
	return doOne[Customer](ctx, s.client, http.MethodPost, joinPath("vault", "customer", customerID, "ach", paymentMethodID), ach.query(), ach)
}

// UpdateToken replaces a stored payment method with the one behind a
// Tokenizer token.
func (s *CustomersService) UpdateToken(ctx context.Context, customerID, paymentMethodID string, token *VaultToken) (*Customer, error) {
	if err := requireIDs("customer id", customerID, "payment method id", paymentMethodID); err != nil {
		return nil, err
	}
	if token == nil {
		return nil, errNilRequest("token")
	}
	return doOne[Customer](ctx, s.client, http.MethodPost, joinPath("vault", "customer", customerID, "token", paymentMethodID), token.query(), token)
}

// DeleteCard removes a stored card.
func (s *CustomersService) DeleteCard(ctx context.Context, customerID, cardID string) (*APIResponse, error) {
	if err := requireIDs("customer id", customerID, "card id", cardID); err != nil {
		return nil, err
	}
	return doEmpty(ctx, s.client, http.MethodDelete, joinPath("vault", "customer", customerID, "card", cardID), nil, nil)
}

// DeleteACH removes a stored bank account.
func (s *CustomersService) DeleteACH(ctx context.Context, customerID, achID string) (*APIResponse, error) {
	if err := requireIDs("customer id", customerID, "ach id", achID); err != nil {
		return nil, err
	}
	return doEmpty(ctx, s.client, http.MethodDelete, joinPath("vault", "customer", customerID, "ach", achID), nil, nil)
}
