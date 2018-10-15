package fluidpay

// GeneralResponse is the general Response containing the status, message and data.
type GeneralResponse struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
	Data   string `json:"data,omitempty"`
}

// ------------------ Authentication --------------------

// JWTTokenRequest is the request body for requesting a new JWT token.
type JWTTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// JWTTokenResponse is the Response from requesting a new JWT token.
type JWTTokenResponse struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
	Token  string `json:"token"`
	Sid    string `json:"sid"`
}

// ForgottenUsernameRequest is the request body for requesting a username reminder e-mail.
type ForgottenUsernameRequest struct {
	Email string `json:"email"`
}

// ForgottenPasswordRequest is the request body for requesting a password reminder a-mail.
type ForgottenPasswordRequest struct {
	Username string `json:"username"`
}

// PasswordResetRequest is the request body for reseting the password.
// The password field is the new password.
type PasswordResetRequest struct {
	Username  string `json:"username"`
	ResetCode string `json:"reset_code"`
	Password  string `json:"password"`
}

// -------------------------------- Users -----------------------------------

// UserResponse is the Response for a user.
type UserResponse struct {
	Status string           `json:"status"`
	Msg    string           `json:"msg"`
	Data   UserResponseData `json:"data,omitempty"`
}

// UserResponseData is the data of UserResponse.
type UserResponseData struct {
	ID               string          `json:"id"`
	Username         string          `json:"username"`
	Name             string          `json:"name"`
	Phone            string          `json:"phone"`
	Email            string          `json:"email"`
	Timezone         string          `json:"timezone"`
	Status           string          `json:"status"`
	Role             string          `json:"role"`
	AccountType      string          `json:"account_type"`
	AccountTypeID    string          `json:"account_type_id"`
	Permissions      UserPermissions `json:"permissions"`
	TwoFactorEnabled bool            `json:"two_factor_enabled"`
	CreatedAt        string          `json:"created_at"`
	UpdatedAt        string          `json:"updated_at"`
}

// UserPermissions is a list of all the permissions a user can have.
type UserPermissions struct {
	ManageUsers           bool `json:"manage_users"`
	ManageAPIKeys         bool `json:"manage_api_keys"`
	ViewBillingReports    bool `json:"view_billing_reports"`
	ManageTerminals       bool `json:"manage_terminals"`
	ManageRuleEngine      bool `json:"manage_rule_engine"`
	ViewSettlementBatches bool `json:"view_settlement_batches"`
	ProcessAuthorization  bool `json:"process_authorization"`
	ProcessCapture        bool `json:"process_capture"`
	ProcessSale           bool `json:"process_sale"`
	ProcessVoid           bool `json:"process_void"`
	ProcessCredit         bool `json:"process_credit"`
	ProcessRefund         bool `json:"process_refund"`
	ProcessVerification   bool `json:"process_verification"`
}

// ChangePasswordRequest is the request body for changing the password.
type ChangePasswordRequest struct {
	Username        string `json:"username"`
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// CreateUserRequest is the request body for creating a new user.
// Status can be active or disabled. Role can be admin or standard.
type CreateUserRequest struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Timezone string `json:"timezone"`
	Password string `json:"password"`
	Status   string `json:"status"`
	Role     string `json:"role"`
}

// UpdateUserRequest is the request body for updating a user.
type UpdateUserRequest struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Timezone string `json:"timezone"`
	Status   string `json:"status"`
	Role     string `json:"role"`
}

//UsersResponse is the Response for all the users.
type UsersResponse struct {
	Status     string             `json:"status"`
	Msg        string             `json:"msg"`
	TotalCount uint               `json:"total_count"`
	Data       []UserResponseData `json:"data,omitempty"`
}

// -------------------------------- API keys -----------------------------------

// KeyRequest is the request to create a new API key.
type KeyRequest struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

// KeyResponse is the Response for creating a new API key.
type KeyResponse struct {
	Status string          `json:"status"`
	Msg    string          `json:"msg"`
	Data   KeyResponseData `json:"data,omitempty"`
}

// KeyResponseData is the data for KeyResponse.
type KeyResponseData struct {
	ID            string `json:"id"`
	UserID        string `json:"user_id"`
	Type          string `json:"type"`
	Name          string `json:"name"`
	APIKey        string `json:"api_key"`
	AccountType   string `json:"account_type"`
	AccountTypeID string `json:"account_type_id"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// KeysResponse is the Response for getting all the keys for a user.
type KeysResponse struct {
	Status     string            `json:"status"`
	Msg        string            `json:"msg"`
	TotalCount uint              `json:"total_count"`
	Data       []KeyResponseData `json:"data,omitempty"`
}

// -------------------------------- Transactions -----------------------------------

// TransactionRequest is the request to process a transaction.
type TransactionRequest struct {
	Type            string                `json:"type"`
	Amount          uint                  `json:"amount"`
	TaxExempt       bool                  `json:"tax_exempt"`
	TaxAmount       uint                  `json:"tax_amount"`
	ShippingAmount  uint                  `json:"shipping_amount"`
	Currency        string                `json:"currency"`
	Description     string                `json:"description"`
	OrderID         string                `json:"order_id"`
	PoNumber        string                `json:"po_number"`
	IPAddress       string                `json:"ip_address"`
	EmailReciept    bool                  `json:"email_reciept"`
	EmailAddress    string                `json:"email_address"`
	PaymentMethod   PaymentMethodsRequest `json:"payment_method,omitempty"`
	BillingAddress  Address               `json:"billing_address,omitempty"`
	ShippingAddress Address               `json:"shipping_address,omitempty"`
}

// PaymentMethodsRequest is the method for the transaction.
type PaymentMethodsRequest struct {
	Card     *CardRequest     `json:"card,omitempty"`
	Customer *CustomerRequest `json:"customer,omitempty"`
	Terminal *TerminalRequest `json:"terminal,omitempty"`
	Ach      *AchRequest      `json:"ach,omitempty"`
}

// CardRequest is the request for a card transaction.
type CardRequest struct {
	EntryType                string                   `json:"entry_type,omitempty"`
	Number                   string                   `json:"number,omitempty"`
	ExpirationDate           string                   `json:"expiration_date,omitempty"`
	Cvc                      string                   `json:"cvc,omitempty"`
	CardholderAuthentication CardholderAuthentication `json:"cardholder_authentication,omitempty"`
}

// CardholderAuthentication is the authentication of the cardholder.
type CardholderAuthentication struct {
	Condition string `json:"condition,omitempty"`
	Eci       string `json:"eci,omitempty"`
	Cavv      string `json:"cavv,omitempty"`
	Xid       string `json:"xid,omitempty"`
}

// CustomerRequest is the request for a customer transaction.
type CustomerRequest struct {
	ID                string `json:"id,omitempty"`
	PaymentMethodType string `json:"payment_method_type,omitempty"`
	PaymentMethodID   string `json:"payment_merhod_id,omitempty"`
	BillingAddressID  string `json:"billing_address_id,omitempty"`
	ShippingAddressID string `json:"shipping_address_id,omitempty"`
}

// TerminalRequest is the request for a terminal transaction.
type TerminalRequest struct {
	ID                string `json:"id,omitempty"`
	ExpirationDate    string `json:"expiration_date,omitempty"`
	Cvc               string `json:"cvc,omitempty"`
	PrintReceipt      string `json:"print_receipt,omitempty"`
	SignatureRequired bool   `json:"signature_required,omitempty"`
}

// AchRequest is the request for a ACH transaction.
type AchRequest struct {
	RoutingNumber               string                      `json:"routing_number,omitempty"`
	AccountNumber               string                      `json:"account_number,omitempty"`
	SecCode                     string                      `json:"sec_code,omitempty"`
	AccountType                 string                      `json:"account_type,omitempty"`
	CheckNumber                 string                      `json:"check_number,omitempty"`
	AccountholderAuthentication AccountholderAuthentication `json:"accountholder_authentication,omitempty"`
}

// AccountholderAuthentication is the authentication of the accountholder.
type AccountholderAuthentication struct {
	DlState  string `json:"dl_state,omitempty"`
	DlNumber string `json:"dl_number,omitempty"`
}

// Address is the general struct for address.
type Address struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Company      string `json:"company"`
	AddressLine1 string `json:"address_line_1"`
	AddressLine2 string `json:"address_line_2,omitempty"`
	City         string `json:"city"`
	State        string `json:"state"`
	PostalCode   string `json:"postal_code"`
	Country      string `json:"country"`
	Phone        string `json:"phone"`
	Fax          string `json:"fax"`
	Email        string `json:"email"`
}

// TransactionResponse is the Response after a transaction.
type TransactionResponse struct {
	Status string                  `json:"status"`
	Msg    string                  `json:"msg"`
	Data   TransactionResponseData `json:"data,omitempty"`
}

// TransactionResponseData is the data of TransactionResponse.
type TransactionResponseData struct {
	ID              string                  `json:"id"`
	Type            string                  `json:"type"`
	Amount          uint                    `json:"amount"`
	TipAmount       uint                    `json:"tip_amount,omitempty"`
	TaxAmount       uint                    `json:"tax_amount"`
	TaxExempt       bool                    `json:"tax_exempt"`
	ShippingAmount  uint                    `json:"shipping_amount"`
	Currency        string                  `json:"currency"`
	Description     string                  `json:"description"`
	OrderID         string                  `json:"order_id"`
	PoNumber        string                  `json:"po_number"`
	IPAddress       string                  `json:"ip_address"`
	EmailReceipt    bool                    `json:"email_receipt"`
	EmailAddress    string                  `json:"email_address"`
	PaymentMethod   string                  `json:"payment_method"`
	Response        TransactionDataResponse `json:"response"`
	Status          string                  `json:"status"`
	BillingAddress  Address                 `json:"billing_address,omitempty"`
	ShippingAddress Address                 `json:"shipping_address,omitempty"`
	CreatedAt       string                  `json:"created_at"`
	UpdatedAt       string                  `json:"updated_at"`
}

// TransactionDataResponse is the response from a transaction.
type TransactionDataResponse struct {
	Card CardResponseBody `json:"card,omitempty"`
}

// CardResponseBody is the card response from a transaction.
type CardResponseBody struct {
	ID                    string            `json:"id"`
	CardType              string            `json:"card_type"`
	FirstSix              string            `json:"first_six"`
	LastFour              string            `json:"last_four"`
	MaskedCard            string            `json:"masked_card"`
	ExpirationDate        string            `json:"expiration_date"`
	Status                string            `json:"status"`
	AuthCode              string            `json:"auth_code"`
	ProcessorResponseCode string            `json:"processor_response_code"`
	ProcessorResponseText string            `json:"processor_response_text"`
	ProcessorType         string            `json:"processor_type"`
	ProcessorID           string            `json:"processor_id"`
	AvsResponseCode       string            `json:"avs_response_code"`
	CvvResponseCode       string            `json:"cvv_response_code"`
	ProcessorSpecific     ProcessorSpecific `json:"processor_specific,omitempty"`
	CreatedAt             string            `json:"created_at"`
	UpdatedAt             string            `json:"updated_at"`
}

// ProcessorSpecific is the processor response of a transaction.
type ProcessorSpecific struct {
	BatchNum             string `json:"BatchNum"`
	CashBack             string `json:"CashBack"`
	ClerkID              string `json:"ClerkID"`
	Disc                 string `json:"DISC"`
	EbtCashAvailBalance  string `json:"EBTCashAvailBalance"`
	EbtCashBeginBalance  string `json:"EBTCashBeginBalance"`
	EbtCashLedgerBalance string `json:"EBTCashLedgerBalance"`
	EbtfsAvailBalance    string `json:"EBTFSAvailBalance"`
	EbtfsBeginBalance    string `json:"EBTFSBeginBalance"`
	EbtfsLedgerBalance   string `json:"EBTFSLedgerBalance"`
	Fee                  string `json:"Fee"`
	InvNum               string `json:"InvNum"`
	Language             string `json:"Language"`
	ProcessData          string `json:"ProcessData"`
	RefNo                string `json:"RefNo"`
	RewardCode           string `json:"RewardCode"`
	RewardQR             string `json:"RewardQR"`
	RwdBalance           string `json:"RwdBalance"`
	RwdIssued            string `json:"RwdIssued"`
	RwdPoints            string `json:"RwdPoints"`
	ShFee                string `json:"SHFee"`
	Svc                  string `json:"SVC"`
	TableNum             string `json:"TableNum"`
	TaxCity              string `json:"tax_city"`
	TaxState             string `json:"tay_state"`
	TicketNum            string `json:"TicketNum"`
	Tip                  string `json:"Tip"`
	TotalAmt             string `json:"TotalAmt"`
}

// TerminalResponse is the Response from the terminal after a transaction.
type TerminalResponse struct {
	Status string               `json:"status"`
	Msg    string               `json:"msg"`
	Data   TerminalResponseData `json:"data,omitempty"`
}

// TerminalResponseData is the data of a TerminalResponse.
type TerminalResponseData struct {
	ID                      string       `json:"id"`
	UserID                  string       `json:"user_id"`
	UserName                string       `json:"user_name"`
	IdempotencyKey          string       `json:"idempotency_key"`
	IdempotencyTime         uint         `json:"idempotency_time"`
	Type                    string       `json:"type"`
	Amount                  uint         `json:"amount"`
	AmountAuthorized        uint         `json:"amount_authorized"`
	AmountCaptured          uint         `json:"amount_captured"`
	AmountSettled           uint         `json:"amount_settled"`
	TipAmount               uint         `json:"tip_amount,omitempty"`
	PaymentAdjustment       uint         `json:"payment_adjustment"`
	PaymentAdjustmentAdd    bool         `json:"payment_adjustment_add"`
	ProcessorID             string       `json:"processor_id"`
	ProcessorType           string       `json:"processor_type"`
	ProcessorName           string       `json:"processor_name"`
	PaymentMethod           string       `json:"payment_method"`
	PaymentType             string       `json:"payment_type"`
	TaxAmount               uint         `json:"tax_amount"`
	TaxExempt               bool         `json:"tax_exempt"`
	ShippingAmount          uint         `json:"shipping_amount"`
	Currency                string       `json:"currency"`
	Description             string       `json:"description"`
	OrderID                 string       `json:"order_id"`
	PoNumber                string       `json:"po_number"`
	IPAddress               string       `json:"ip_address"`
	TransactionSource       string       `json:"trasaction_source"`
	EmailReceipt            bool         `json:"email_receipt"`
	EmailAddress            string       `json:"email_address"`
	CustomerID              string       `json:"customer_id"`
	ReferencedTransactionID string       `json:"referenced_transaction_is"`
	ResponseBody            ResponseBody `json:"response_body"`
	Status                  string       `json:"status"`
	Response                string       `json:"Response"`
	ResponseCode            uint         `json:"response_code"`
	BillingAddress          Address      `json:"billing_address"`
	ShippingAddress         Address      `json:"shipping_address"`
	CreatedAt               string       `json:"created_at"`
	UpdatedAt               string       `json:"updated_at"`
	CapturedAt              string       `json:"captured_at"`
	SettledAt               string       `json:"settled_at,omitempty"`
}

// ResponseBody is the detailed response from a transaction.
type ResponseBody struct {
	Card     CardResponseBody     `json:"card,omitempty"`
	Terminal TerminalResponseBody `json:"terminal,omitempty"`
	Ach      AchResponseBody      `json:"ach,omitempty"`
}

// TerminalResponseBody is the terminal response from a transaction.
type TerminalResponseBody struct {
	ID                    string            `json:"id"`
	TerminalID            string            `json:"terminal_id"`
	TerminalDescription   string            `json:"terminal_description"`
	CardType              string            `json:"card_type"`
	PaymentType           string            `json:"payment_type"`
	EntryType             string            `json:"entry_type"`
	FirstFour             string            `json:"first_four"`
	LastFour              string            `json:"last_four"`
	MaskedCard            string            `json:"masked_card"`
	CardholderName        string            `json:"cardholder_name"`
	AuthCode              string            `json:"auth_code"`
	ResponseCode          uint              `json:"response_code"`
	ProcessorResponseText string            `json:"processor_response_text"`
	ProcessorSpecific     ProcessorSpecific `json:"processor_specific,omitempty"`
	EmvAid                string            `json:"emv_aid"`
	EmvAppName            string            `json:"emv_app_name"`
	EmvTvr                string            `json:"emv_tvr"`
	EmvTsi                string            `json:"emv_tsi"`
	SignatureData         string            `json:"signature_data"`
	CreatedAt             string            `json:"created_at"`
	UpdatedAt             string            `json:"updated_at"`
}

// AchResponseBody is the ACH response from a transaction.
type AchResponseBody struct {
	ID                    string      `json:"id"`
	SecCode               string      `json:"sec_code"`
	AccountType           string      `json:"account_type"`
	MaskedAccountNumber   string      `json:"masked_account_number"`
	RoutingNumber         string      `json:"routing_number"`
	AuthCode              string      `json:"auth_code"`
	Response              uint        `json:"Response"`
	ResponseCode          string      `json:"response_code"`
	ProcessorResponseCode string      `json:"processor_response_code"`
	ProcessorResponseText string      `json:"processor_response_text"`
	ProcessorType         string      `json:"processor_type"`
	ProcessorID           string      `json:"processor_id"`
	ProcessorSpecific     AchSpecific `json:"processor_specific"`
	CreatedAt             string      `json:"created_at"`
	UpdatedAt             string      `json:"updated_at"`
}

// AchSpecific is the ACH specific processor response.
type AchSpecific struct {
	ResultCodes []string `json:"result_codes"`
	TypeCodes   []string `json:"type_codes"`
}

// TransactionSearchResponse is the Response of a specific or queried transaction.
type TransactionSearchResponse struct {
	Status     string                  `json:"status"`
	Msg        string                  `json:"msg"`
	TotalCount uint                    `json:"total_count"`
	Data       TransactionResponseData `json:"data,omitempty"`
}

// TransactionQueryRequest is the request for querying a transaction.
// All fields are optional.
type TransactionQueryRequest struct {
	TransactionID    StringQuery  `json:"transaction_id,omitempty"`
	UserID           StringQuery  `json:"user_id,omitempty"`
	Type             StringQuery  `json:"type,omitempty"`
	IPAddress        StringQuery  `json:"ip_address,omitempty"`
	Amount           IntQuery     `json:"amount,omitempty"`
	AmountAuthorized IntQuery     `json:"amount_authorized,omitempty"`
	AmountCaptured   IntQuery     `json:"amount_captured,omitempty"`
	AmountSettled    IntQuery     `json:"amount_settled,omitempty"`
	TaxAmount        IntQuery     `json:"tax_amount,omitempty"`
	PoNumber         StringQuery  `json:"po_number,omitempty"`
	OrderID          StringQuery  `json:"order_id,omitempty"`
	PaymentMethod    StringQuery  `json:"payment_method,omitempty"`
	PaymentType      StringQuery  `json:"payment_type,omitempty"`
	Status           StringQuery  `json:"status,omitempty"`
	ProcessorID      StringQuery  `json:"processor_id,omitempty"`
	CustomerID       StringQuery  `json:"customer_id,omitempty"`
	CreatedAt        DateQuery    `json:"created_at,omitempty"`
	CapturedAt       DateQuery    `json:"captured_at,omitempty"`
	SettledAt        DateQuery    `json:"settled_at,omitempty"`
	BillingAddress   AddressQuery `json:"billing_address,omitempty"`
	ShippingAddress  AddressQuery `json:"shipping_address,omitempty"`
	Limit            uint         `json:"limit,omitempty"`
	Offset           uint         `json:"offset,omitempty"`
}

// AddressQuery is the query for an address.
type AddressQuery struct {
	AddressID    StringQuery `json:"address_id,omitempty"`
	FirstName    StringQuery `json:"first_name,omitempty"`
	LastName     StringQuery `json:"last_name,omitempty"`
	Company      StringQuery `json:"company,omitempty"`
	AddressLine1 StringQuery `json:"address_line_1,omitempty"`
	AddressLine2 StringQuery `json:"address_line_2,omitempty"`
	City         StringQuery `json:"city,omitempty"`
	State        StringQuery `json:"state,omitempty"`
	PostalCode   StringQuery `json:"postal_code,omitempty"`
	Country      StringQuery `json:"country,omitempty"`
	Email        StringQuery `json:"email,omitempty"`
	Phone        StringQuery `json:"phone,omitempty"`
	Fax          StringQuery `json:"fax,omitempty"`
}

// StringQuery is the query for a string.
type StringQuery struct {
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

// IntQuery is the query for an integer.
type IntQuery struct {
	Operator string `json:"operator"`
	Value    uint   `json:"value"`
}

// DateQuery is the query for a date.
type DateQuery struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// TransactionCaptureRequest is the request to capture a transaction.
type TransactionCaptureRequest struct {
	Amount         uint   `json:"amount"`
	TaxAmount      uint   `json:"tax_amount"`
	TaxExempt      bool   `json:"tax_exempt"`
	ShippingAmount uint   `json:"shipping_amount"`
	OrderID        string `json:"order_id"`
	PoNumber       string `json:"po_number"`
	IPAddress      string `json:"ip_address"`
}

// TransactionRefundRequest is the request to refund a transaction.
type TransactionRefundRequest struct {
	Amount uint `json:"amount"`
}

// -------------------------------- Customers -----------------------------------

// CreateCustomerRequest is the request to create a new customer token.
type CreateCustomerRequest struct {
	Description     string                       `json:"description"`
	PaymentMethod   CustomerPaymentMethodRequest `json:"payment_method"`
	BillingAddress  Address                      `json:"billing_address"`
	ShippingAddress Address                      `json:"shipping_address"`
}

// CustomerPaymentMethodRequest is the payment method for a customer.
type CustomerPaymentMethodRequest struct {
	Card CustomerPaymentRequest `json:"card"`
}

// CustomerResponse is the Response for a customer token.
type CustomerResponse struct {
	Status string               `json:"status"`
	Msg    string               `json:"msg"`
	Data   CustomerResponseData `json:"data,omitempty"`
}

// CustomerResponseData is the data of CustomerResponse.
type CustomerResponseData struct {
	ID                string                        `json:"id"`
	Description       string                        `json:"description"`
	PaymentMethod     CustomerPaymentMethodResponse `json:"payment_method"`
	PaymentMethodType string                        `json:"payment_method_type"`
	BillingAddress    Address                       `json:"billing_address"`
	ShippingAddress   Address                       `json:"shipping_address"`
	CreatedAt         string                        `json:"created_at"`
	UpdatedAt         string                        `json:"updated_at"`
}

// CustomerPaymentMethodResponse is the payment method response for a customer.
type CustomerPaymentMethodResponse struct {
	Card CustomerCardResponse `json:"card"`
	Ach  string               `json:"ach"`
}

// CustomerCardResponse is the card response for a customer.
type CustomerCardResponse struct {
	ID             string `json:"id"`
	CardType       string `json:"card_type"`
	FirstSix       string `json:"first_six"`
	LastFour       string `json:"last_four"`
	MaskedCard     string `json:"masked_card"`
	ExpirationDate string `json:"expiration_date"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

// UpdateCustomerRequest is the request to update a customer token.
type UpdateCustomerRequest struct {
	Description       string `json:"description"`
	PaymentMethod     string `json:"payment_method"`
	PaymentMethodID   string `json:"payment_method_id"`
	BillingAddressID  string `json:"billing_address_id"`
	ShippingAddressID string `json:"shipping_address_id"`
}

// CustomerAddressRequest is the request to create or update a customer Address token.
type CustomerAddressRequest struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Company      string `json:"company"`
	AddressLine1 string `json:"address_line_1"`
	AddressLine2 string `json:"address_line_2,omitempty"`
	City         string `json:"city"`
	State        string `json:"state"`
	PostalCode   string `json:"postal_code"`
	Country      string `json:"country"`
	Phone        string `json:"phone"`
	Fax          string `json:"fax"`
	Email        string `json:"email"`
}

// CustomerAddressResponse is the Response from a customer Address token.
type CustomerAddressResponse struct {
	Status string                      `json:"status"`
	Msg    string                      `json:"msg"`
	Data   CustomerAddressResponseData `json:"data,omitempty"`
}

// CustomerAddressesResponse is the Response from getting all of the customer's Address tokens.
type CustomerAddressesResponse struct {
	Status     string                        `json:"status"`
	Msg        string                        `json:"msg"`
	Data       []CustomerAddressResponseData `json:"data,omitempty"`
	TotalCount uint                          `json:"total_count"`
}

// CustomerAddressResponseData is the data of CustomerAddressResponse.
type CustomerAddressResponseData struct {
	ID           string `json:"id"`
	CustomerID   string `ȷson:"customer_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Company      string `json:"company"`
	AddressLine1 string `json:"address_line_1"`
	AddressLine2 string `json:"address_line_2,omitempty"`
	City         string `json:"city"`
	State        string `json:"state"`
	PostalCode   string `json:"postal_code"`
	Country      string `json:"country"`
	Phone        string `json:"phone"`
	Fax          string `json:"fax"`
	Email        string `json:"email"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// CustomerPaymentRequest is the request to create or update a customer payment token.
type CustomerPaymentRequest struct {
	CardNumber     string `json:"card_number"`
	ExpirationDate string `json:"expiration_date"`
	EntryType      string `json:"entry_type,omitempty"`
}

// CustomerPaymentResponse is the Response from a customer payment token.
type CustomerPaymentResponse struct {
	Status string                        `json:"status"`
	Msg    string                        `json:"msg"`
	Data   CustomerPaymentMethodResponse `json:"data,omitempty"`
}

// CustomerPaymentsResponse is the Response from getting all of the customer's payment tokens.
type CustomerPaymentsResponse struct {
	Status     string                          `json:"status"`
	Msg        string                          `json:"msg"`
	Data       []CustomerPaymentMethodResponse `json:"data"`
	TotalCount uint                            `json:"total_count"`
}

// -------------------------------- Recurring -----------------------------------

// RecurrenceRequest is the request to create or update an add on or discount.
type RecurrenceRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Amount      uint   `json:"amount,omitempty"`
	Percentage  uint   `json:"percentage,omitempty"`
	Duration    uint   `json:"duration"`
}

// RecurrenceResponse is the Response from a recurring.
type RecurrenceResponse struct {
	Status string                 `json:"status"`
	Msg    string                 `json:"msg"`
	Data   RecurrenceResponseData `json:"data,omitempty"`
}

// RecurrenceResponseData is the data of RecurrenceResponse.
type RecurrenceResponseData struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Amount      uint   `json:"amount,omitempty"`
	Percentage  uint   `json:"percentage,omitempty"`
	Duration    uint   `json:"duration"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// RecurrencesResponse is the Response from getting all the recurrences.
type RecurrencesResponse struct {
	Status     string                   `json:"status"`
	Msg        string                   `json:"msg"`
	Data       []RecurrenceResponseData `json:"data,omitempty"`
	TotalCount uint                     `json:"total_count"`
}

// PlanRequest is the request to create or update a plan.
type PlanRequest struct {
	Name                 string                     `json:"name"`
	Description          string                     `json:"description"`
	Amount               uint                       `json:"amount"`
	BillingCycleInterval uint                       `json:"billing_cycle_interval"`
	BillingFrequency     string                     `json:"billing_frequency"`
	BillingDays          string                     `json:"billing_days"`
	Duration             uint                       `json:"duration,omitempty"`
	AddOns               []OptionalRecurringRequest `json:"add_ons,omitempty"`
	Discounts            []OptionalRecurringRequest `json:"discounts,omitempty"`
}

// OptionalRecurringRequest is for optional overwriting a recurring.
type OptionalRecurringRequest struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Amount      uint   `json:"amount,omitempty"`
	Percentage  uint   `json:"percentage,omitempty"`
	Duration    uint   `json:"duration,omitempty"`
}

// PlanResponse is the Response from a plan.
type PlanResponse struct {
	Status string           `json:"status"`
	Msg    string           `json:"msg"`
	Data   PlanResponseData `json:"data,omitempty"`
}

// PlanResponseData is the data of PlanResponse.
type PlanResponseData struct {
	ID                   string                   `json:"id"`
	Name                 string                   `json:"name"`
	Description          string                   `json:"description"`
	Amount               uint                     `json:"amount"`
	Percentage           uint                     `json:"percentage"`
	BillingCycleInterval uint                     `json:"billing_cycle_interval"`
	BillingFrequency     string                   `json:"billing_frequency"`
	BillingDays          string                   `json:"billing_days"`
	TotalAddOns          uint                     `json:"total_add_ons"`
	TotalDiscounts       uint                     `json:"total_discounts"`
	Duration             uint                     `json:"duration"`
	AddOns               []RecurrenceResponseData `json:"add_ons"`
	Discounts            []RecurrenceResponseData `json:"discounts"`
	CreatedAt            string                   `json:"created_at"`
	UpdatedAt            string                   `json:"updated_at"`
}

// PlansResponse is the Response from getting all the plans.
type PlansResponse struct {
	Status     string             `json:"status"`
	Msg        string             `json:"msg"`
	Data       []PlanResponseData `json:"data,omitempty"`
	TotalCount uint               `json:"total_count"`
}

// SubscriptionRequest is the request for creating or updating a subscription.
type SubscriptionRequest struct {
	PlanID               string                     `json:"plan_id"`
	Description          string                     `json:"description"`
	Customer             IDData                     `json:"customer"`
	Amount               uint                       `json:"amount"`
	BillingCycleInterval uint                       `json:"billing_cycle_interval"`
	BillingFrequency     string                     `json:"billing_frequency"`
	BillingDays          string                     `json:"billing_days"`
	Duration             uint                       `json:"duration,omitempty"`
	NextBillDate         string                     `json:"next_bill_date,omitempty"`
	AddOns               []OptionalRecurringRequest `json:"add_ons,omitempty"`
	Discounts            []OptionalRecurringRequest `json:"discounts,omitempty"`
}

// IDData is the id of the customer for a subscription.
type IDData struct {
	ID string `json:"id"`
}

// SubscriptionResponse is the Response from a subscription.
type SubscriptionResponse struct {
	Status string                   `json:"status"`
	Msg    string                   `json:"msg"`
	Data   SubscriptionResponseData `json:"data,omitempty"`
}

// SubscriptionsResponse is the Response from getting a subscription.
type SubscriptionsResponse struct {
	Status     string                     `json:"status"`
	Msg        string                     `json:"msg"`
	Data       []SubscriptionResponseData `json:"data,omitempty"`
	TotalCount uint                       `json:"total_count"`
}

// SubscriptionResponseData is the data of SubscriptionResponse.
type SubscriptionResponseData struct {
	ID                   string                   `json:"id"`
	PlanID               string                   `json:"plan_id"`
	Status               string                   `json:"status"`
	Description          string                   `json:"description"`
	Customer             IDData                   `json:"customer"`
	Amount               uint                     `json:"amount"`
	TotalAdds            uint                     `json:"total_adds"`
	TotalDiscounts       uint                     `json:"total_discounts"`
	BillingCycleInterval uint                     `json:"billing_cycle_interval"`
	BillingFrequency     string                   `json:"billing_frequency"`
	BillingDays          string                   `json:"billing_days"`
	Duration             uint                     `json:"duration"`
	NextBillDate         string                   `json:"next_bill_date"`
	AddOns               []RecurrenceResponseData `json:"add_ons"`
	Discounts            []RecurrenceResponseData `json:"discounts"`
	CreatedAt            string                   `json:"created_at"`
	UpdatedAt            string                   `json:"updated_at"`
}

// -------------------------------- Terminals -----------------------------------

// TerminalsResponse is the Response from getting all the terminals.
type TerminalsResponse struct {
	Status     string                  `json:"status"`
	Msg        string                  `json:"msg"`
	TotalCount uint                    `json:"total_count"`
	Data       []TerminalsDataResponse `json:"data"`
}

// TerminalsDataResponse is the data of TerminalsResponse.
type TerminalsDataResponse struct {
	ID           string `json:"id"`
	MerchantID   string `json:"merchant_id"`
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"model"`
	SerialNumber string `json:"serial_number"`
	Tpn          string `json:"tpn"`
	Description  string `json:"description"`
	Status       string `json:"status"`
	AuthKey      string `json:"auth_key"`
	RegisterID   string `json:"register_id"`
	AutoSettle   bool   `json:"auto_settle"`
	SettleAt     string `json:"settle_at"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}
