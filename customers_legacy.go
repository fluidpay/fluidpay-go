package fluidpay

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

// LegacyCustomersService wraps the deprecated /customer endpoints. FluidPay
// no longer develops them and warns that mixing them with the Customer
// Vault in one application can produce unintended behaviour.
//
// Deprecated: use CustomersService (Client.Customers) for new integrations.
type LegacyCustomersService struct{ client *Client }

// LegacyCustomerRequest creates a customer through the deprecated API.
//
// Deprecated: use CustomerCreateRequest with Client.Customers.
type LegacyCustomerRequest struct {
	// Validate runs a $0 verification on a card while storing it.
	Validate bool `json:"-"`

	Description     string               `json:"description,omitempty"`
	PaymentMethod   *LegacyPaymentMethod `json:"payment_method,omitempty"`
	BillingAddress  *Address             `json:"billing_address,omitempty"`
	ShippingAddress *Address             `json:"shipping_address,omitempty"`
}

// LegacyPaymentMethod is the payment method of a LegacyCustomerRequest.
//
// Deprecated: part of the deprecated customer API.
type LegacyPaymentMethod struct {
	ProcessorID   string           `json:"processor_id,omitempty"`
	Card          *LegacyCardInput `json:"card,omitempty"`
	ACH           *LegacyACHInput  `json:"ach,omitempty"`
	ApplePayToken string           `json:"apple_pay_token,omitempty"`
}

// LegacyCardInput is a card in a LegacyCustomerRequest. Note the field is
// card_number here but number when creating a payment method token.
//
// Deprecated: part of the deprecated customer API.
type LegacyCardInput struct {
	CardNumber     string `json:"card_number"`
	ExpirationDate string `json:"expiration_date"`
}

// LegacyACHInput is a bank account in the deprecated customer API.
//
// Deprecated: part of the deprecated customer API.
type LegacyACHInput struct {
	AccountNumber string `json:"account_number"`
	RoutingNumber string `json:"routing_number"`
	AccountType   string `json:"account_type"`
	SecCode       string `json:"sec_code"`
}

// LegacyCustomerUpdateRequest sets a customer's description and default
// tokens. Only gateway-generated 20 character IDs are accepted.
//
// Deprecated: part of the deprecated customer API.
type LegacyCustomerUpdateRequest struct {
	Description string `json:"description,omitempty"`
	// PaymentMethod is PaymentMethodTypeCard or PaymentMethodTypeACH.
	PaymentMethod     string `json:"payment_method,omitempty"`
	PaymentMethodID   string `json:"payment_method_id,omitempty"`
	BillingAddressID  string `json:"billing_address_id,omitempty"`
	ShippingAddressID string `json:"shipping_address_id,omitempty"`
}

// LegacyCustomerSearchRequest filters deprecated customers.
//
// Deprecated: part of the deprecated customer API.
type LegacyCustomerSearchRequest struct {
	ID                *StringQuery    `json:"id,omitempty"`
	PaymentMethodID   *StringQuery    `json:"payment_method_id,omitempty"`
	BillingAddressID  *StringQuery    `json:"billing_address_id,omitempty"`
	ShippingAddressID *StringQuery    `json:"shipping_address_id,omitempty"`
	CreatedAt         *DateRangeQuery `json:"created_at,omitempty"`
	Limit             int             `json:"limit,omitempty"`
	Offset            int             `json:"offset,omitempty"`
}

// LegacyCustomer is a customer from the deprecated API.
//
// Deprecated: part of the deprecated customer API.
type LegacyCustomer struct {
	APIResource

	ID              string               `json:"id"`
	Description     string               `json:"description"`
	PaymentMethod   LegacyPaymentMethods `json:"payment_method"`
	BillingAddress  *LegacyAddress       `json:"billing_address"`
	ShippingAddress *LegacyAddress       `json:"shipping_address"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
}

// LegacyPaymentMethods holds the stored payment method(s) of a
// LegacyCustomer, or the payment method just created.
//
// Deprecated: part of the deprecated customer API.
type LegacyPaymentMethods struct {
	APIResource

	Card *LegacyCard `json:"card,omitempty"`
	ACH  *LegacyACH  `json:"ach,omitempty"`
}

// LegacyCard is a stored card token.
//
// Deprecated: part of the deprecated customer API.
type LegacyCard struct {
	APIResource

	ID             string    `json:"id"`
	CardType       string    `json:"card_type"`
	FirstSix       string    `json:"first_six"`
	LastFour       string    `json:"last_four"`
	MaskedCard     string    `json:"masked_card"`
	ExpirationDate string    `json:"expiration_date"`
	ProcessorID    string    `json:"processor_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// LegacyACH is a stored bank account token.
//
// Deprecated: part of the deprecated customer API.
type LegacyACH struct {
	APIResource

	ID                  string    `json:"id"`
	SecCode             string    `json:"sec_code"`
	AccountType         string    `json:"account_type"`
	MaskedAccountNumber string    `json:"masked_account_number"`
	RoutingNumber       string    `json:"routing_number"`
	ProcessorID         string    `json:"processor_id"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// LegacyAddress is a stored address token.
//
// Deprecated: part of the deprecated customer API.
type LegacyAddress struct {
	APIResource
	Address

	ID         string    `json:"id"`
	CustomerID string    `json:"customer_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// LegacyCustomerList is a page of deprecated customers.
//
// Deprecated: part of the deprecated customer API.
type LegacyCustomerList struct {
	APIResource
	Data       []*LegacyCustomer
	TotalCount int
}

// LegacyAddressList is a list of stored address tokens.
//
// Deprecated: part of the deprecated customer API.
type LegacyAddressList struct {
	APIResource
	Data       []*LegacyAddress
	TotalCount int
}

// LegacyCardList is a list of stored card tokens.
//
// Deprecated: part of the deprecated customer API.
type LegacyCardList struct {
	APIResource
	Data       []*LegacyCard
	TotalCount int
}

// LegacyACHList is a list of stored bank account tokens.
//
// Deprecated: part of the deprecated customer API.
type LegacyACHList struct {
	APIResource
	Data       []*LegacyACH
	TotalCount int
}

// LegacyCardTokenRequest creates or updates a stored card token.
//
// Deprecated: part of the deprecated customer API.
type LegacyCardTokenRequest struct {
	Validate bool `json:"-"`

	Number         string `json:"number"`
	ExpirationDate string `json:"expiration_date"`
}

// LegacyACHTokenRequest creates or updates a stored bank account token.
//
// Deprecated: part of the deprecated customer API.
type LegacyACHTokenRequest struct {
	Validate bool `json:"-"`

	AccountNumber string `json:"account_number"`
	RoutingNumber string `json:"routing_number"`
	AccountType   string `json:"account_type"`
	SecCode       string `json:"sec_code"`
}

// LegacyTokenRequest creates a stored payment method from a Tokenizer token.
//
// Deprecated: part of the deprecated customer API.
type LegacyTokenRequest struct {
	Validate bool `json:"-"`

	Token string `json:"token"`
}

func validateQuery(validate bool) url.Values {
	if !validate {
		return nil
	}
	return url.Values{"validate": []string{"true"}}
}

// Create creates a customer.
//
// Deprecated: use Client.Customers.Create.
func (s *LegacyCustomersService) Create(ctx context.Context, req *LegacyCustomerRequest) (*LegacyCustomer, error) {
	if req == nil {
		return nil, errNilRequest("customer")
	}
	return doOne[LegacyCustomer](ctx, s.client, http.MethodPost, "customer", validateQuery(req.Validate), req)
}

// Get retrieves a customer.
//
// Deprecated: use Client.Customers.Get.
func (s *LegacyCustomersService) Get(ctx context.Context, customerID string) (*LegacyCustomer, error) {
	if err := requireID("customer id", customerID); err != nil {
		return nil, err
	}
	return doOne[LegacyCustomer](ctx, s.client, http.MethodGet, joinPath("customer", customerID), nil, nil)
}

// Search returns customers matching req.
//
// Deprecated: use Client.Customers.Search.
func (s *LegacyCustomersService) Search(ctx context.Context, req *LegacyCustomerSearchRequest) (*LegacyCustomerList, error) {
	if req == nil {
		req = &LegacyCustomerSearchRequest{}
	}
	items, resp, err := doList[LegacyCustomer](ctx, s.client, http.MethodPost, "customer/search", nil, req)
	if err != nil {
		return nil, err
	}
	return &LegacyCustomerList{APIResource: APIResource{LastResponse: resp}, Data: items, TotalCount: resp.TotalCount}, nil
}

// Update sets the description and default tokens of a customer.
//
// Deprecated: use Client.Customers.Update.
func (s *LegacyCustomersService) Update(ctx context.Context, customerID string, req *LegacyCustomerUpdateRequest) error {
	if err := requireID("customer id", customerID); err != nil {
		return err
	}
	if req == nil {
		return errNilRequest("customer update")
	}
	return doEmpty(ctx, s.client, http.MethodPost, joinPath("customer", customerID), nil, req)
}

// Delete removes a customer.
//
// Deprecated: use Client.Customers.Delete.
func (s *LegacyCustomersService) Delete(ctx context.Context, customerID string) error {
	if err := requireID("customer id", customerID); err != nil {
		return err
	}
	return doEmpty(ctx, s.client, http.MethodDelete, joinPath("customer", customerID), nil, nil)
}

// CreateAddress stores an address token.
//
// Deprecated: use Client.Customers.CreateAddress.
func (s *LegacyCustomersService) CreateAddress(ctx context.Context, customerID string, addr *Address) (*LegacyAddress, error) {
	if err := requireID("customer id", customerID); err != nil {
		return nil, err
	}
	if addr == nil {
		return nil, errNilRequest("address")
	}
	return doOne[LegacyAddress](ctx, s.client, http.MethodPost, joinPath("customer", customerID, "address"), nil, addr)
}

// GetAddress retrieves an address token.
//
// Deprecated: use Client.Customers.Get and Customer.Address.
func (s *LegacyCustomersService) GetAddress(ctx context.Context, customerID, addressID string) (*LegacyAddress, error) {
	if err := requireIDs("customer id", customerID, "address id", addressID); err != nil {
		return nil, err
	}
	return doOne[LegacyAddress](ctx, s.client, http.MethodGet, joinPath("customer", customerID, "address", addressID), nil, nil)
}

// ListAddresses retrieves all address tokens of a customer.
//
// Deprecated: use Client.Customers.Get and Customer.Addresses.
func (s *LegacyCustomersService) ListAddresses(ctx context.Context, customerID string) (*LegacyAddressList, error) {
	if err := requireID("customer id", customerID); err != nil {
		return nil, err
	}
	items, resp, err := doList[LegacyAddress](ctx, s.client, http.MethodGet, joinPath("customer", customerID, "addresses"), nil, nil)
	if err != nil {
		return nil, err
	}
	return &LegacyAddressList{APIResource: APIResource{LastResponse: resp}, Data: items, TotalCount: resp.TotalCount}, nil
}

// UpdateAddress replaces an address token.
//
// Deprecated: use Client.Customers.UpdateAddress.
func (s *LegacyCustomersService) UpdateAddress(ctx context.Context, customerID, addressID string, addr *Address) (*LegacyAddress, error) {
	if err := requireIDs("customer id", customerID, "address id", addressID); err != nil {
		return nil, err
	}
	if addr == nil {
		return nil, errNilRequest("address")
	}
	return doOne[LegacyAddress](ctx, s.client, http.MethodPost, joinPath("customer", customerID, "address", addressID), nil, addr)
}

// DeleteAddress removes an address token.
//
// Deprecated: use Client.Customers.DeleteAddress.
func (s *LegacyCustomersService) DeleteAddress(ctx context.Context, customerID, addressID string) error {
	if err := requireIDs("customer id", customerID, "address id", addressID); err != nil {
		return err
	}
	return doEmpty(ctx, s.client, http.MethodDelete, joinPath("customer", customerID, "address", addressID), nil, nil)
}

// CreateCard stores a card token.
//
// Deprecated: use Client.Customers.CreateCard.
func (s *LegacyCustomersService) CreateCard(ctx context.Context, customerID string, req *LegacyCardTokenRequest) (*LegacyPaymentMethods, error) {
	if err := requireID("customer id", customerID); err != nil {
		return nil, err
	}
	if req == nil {
		return nil, errNilRequest("card")
	}
	return doOne[LegacyPaymentMethods](ctx, s.client, http.MethodPost, joinPath("customer", customerID, "paymentmethod", "card"), validateQuery(req.Validate), req)
}

// CreateACH stores a bank account token.
//
// Deprecated: use Client.Customers.CreateACH.
func (s *LegacyCustomersService) CreateACH(ctx context.Context, customerID string, req *LegacyACHTokenRequest) (*LegacyPaymentMethods, error) {
	if err := requireID("customer id", customerID); err != nil {
		return nil, err
	}
	if req == nil {
		return nil, errNilRequest("ach")
	}
	return doOne[LegacyPaymentMethods](ctx, s.client, http.MethodPost, joinPath("customer", customerID, "paymentmethod", "ach"), validateQuery(req.Validate), req)
}

// CreateToken stores the payment method behind a Tokenizer token.
//
// Deprecated: use Client.Customers.CreateToken.
func (s *LegacyCustomersService) CreateToken(ctx context.Context, customerID string, req *LegacyTokenRequest) (*LegacyPaymentMethods, error) {
	if err := requireID("customer id", customerID); err != nil {
		return nil, err
	}
	if req == nil {
		return nil, errNilRequest("token")
	}
	return doOne[LegacyPaymentMethods](ctx, s.client, http.MethodPost, joinPath("customer", customerID, "paymentmethod", "token"), validateQuery(req.Validate), req)
}

// GetCard retrieves a stored card token.
//
// Deprecated: use Client.Customers.Get and Customer.Card.
func (s *LegacyCustomersService) GetCard(ctx context.Context, customerID, cardID string) (*LegacyCard, error) {
	if err := requireIDs("customer id", customerID, "card id", cardID); err != nil {
		return nil, err
	}
	return doOne[LegacyCard](ctx, s.client, http.MethodGet, joinPath("customer", customerID, "paymentmethod", "card", cardID), nil, nil)
}

// GetACH retrieves a stored bank account token.
//
// Deprecated: use Client.Customers.Get and Customer.ACH.
func (s *LegacyCustomersService) GetACH(ctx context.Context, customerID, achID string) (*LegacyACH, error) {
	if err := requireIDs("customer id", customerID, "ach id", achID); err != nil {
		return nil, err
	}
	return doOne[LegacyACH](ctx, s.client, http.MethodGet, joinPath("customer", customerID, "paymentmethod", "ach", achID), nil, nil)
}

// ListCards retrieves all stored card tokens of a customer.
//
// Deprecated: use Client.Customers.Get and Customer.Cards.
func (s *LegacyCustomersService) ListCards(ctx context.Context, customerID string) (*LegacyCardList, error) {
	if err := requireID("customer id", customerID); err != nil {
		return nil, err
	}
	items, resp, err := doList[LegacyCard](ctx, s.client, http.MethodGet, joinPath("customer", customerID, "paymentmethod", "card"), nil, nil)
	if err != nil {
		return nil, err
	}
	return &LegacyCardList{APIResource: APIResource{LastResponse: resp}, Data: items, TotalCount: resp.TotalCount}, nil
}

// ListACH retrieves all stored bank account tokens of a customer.
//
// Deprecated: use Client.Customers.Get and Customer.ACH.
func (s *LegacyCustomersService) ListACH(ctx context.Context, customerID string) (*LegacyACHList, error) {
	if err := requireID("customer id", customerID); err != nil {
		return nil, err
	}
	items, resp, err := doList[LegacyACH](ctx, s.client, http.MethodGet, joinPath("customer", customerID, "paymentmethod", "ach"), nil, nil)
	if err != nil {
		return nil, err
	}
	return &LegacyACHList{APIResource: APIResource{LastResponse: resp}, Data: items, TotalCount: resp.TotalCount}, nil
}

// UpdateCard replaces a stored card token. The gateway expects the card
// fields nested under "card" and returns no data.
//
// Deprecated: use Client.Customers.UpdateCard.
func (s *LegacyCustomersService) UpdateCard(ctx context.Context, customerID, cardID string, card *LegacyCardInput) error {
	if err := requireIDs("customer id", customerID, "card id", cardID); err != nil {
		return err
	}
	if card == nil {
		return errNilRequest("card")
	}
	body := struct {
		Card *LegacyCardInput `json:"card"`
	}{card}
	return doEmpty(ctx, s.client, http.MethodPost, joinPath("customer", customerID, "paymentmethod", "card", cardID), nil, body)
}

// UpdateACH replaces a stored bank account token.
//
// Deprecated: use Client.Customers.UpdateACH.
func (s *LegacyCustomersService) UpdateACH(ctx context.Context, customerID, achID string, ach *LegacyACHInput) error {
	if err := requireIDs("customer id", customerID, "ach id", achID); err != nil {
		return err
	}
	if ach == nil {
		return errNilRequest("ach")
	}
	body := struct {
		ACH *LegacyACHInput `json:"ach"`
	}{ach}
	return doEmpty(ctx, s.client, http.MethodPost, joinPath("customer", customerID, "paymentmethod", "ach", achID), nil, body)
}

// DeleteCard removes a stored card token.
//
// Deprecated: use Client.Customers.DeleteCard.
func (s *LegacyCustomersService) DeleteCard(ctx context.Context, customerID, cardID string) error {
	if err := requireIDs("customer id", customerID, "card id", cardID); err != nil {
		return err
	}
	return doEmpty(ctx, s.client, http.MethodDelete, joinPath("customer", customerID, "paymentmethod", "card", cardID), nil, nil)
}

// DeleteACH removes a stored bank account token.
//
// Deprecated: use Client.Customers.DeleteACH.
func (s *LegacyCustomersService) DeleteACH(ctx context.Context, customerID, achID string) error {
	if err := requireIDs("customer id", customerID, "ach id", achID); err != nil {
		return err
	}
	return doEmpty(ctx, s.client, http.MethodDelete, joinPath("customer", customerID, "paymentmethod", "ach", achID), nil, nil)
}
