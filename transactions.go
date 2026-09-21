package fluidpay

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// TransactionsService processes payments and looks up transaction history.
// Access it through Client.Transactions.
type TransactionsService struct{ client *Client }

// Transaction types accepted by TransactionRequest.Type.
const (
	TransactionTypeSale         = "sale"
	TransactionTypeAuthorize    = "authorize"
	TransactionTypeVerification = "verification"
	TransactionTypeCredit       = "credit"
)

// Transaction statuses reported in Transaction.Status.
const (
	StatusUnknown           = "unknown"
	StatusPending           = "pending"
	StatusDeclined          = "declined"
	StatusAuthorized        = "authorized"
	StatusPendingSettlement = "pending_settlement"
	StatusSettled           = "settled"
	StatusVoided            = "voided"
	StatusReversed          = "reversed"
	StatusRefunded          = "refunded"
	StatusPartiallyRefunded = "partially_refunded"
	StatusReturned          = "returned"
	StatusLateReturn        = "late_return"
	StatusFlagged           = "flagged"
	StatusFlaggedPartner    = "flagged_partner"
)

// Well known response codes. Codes 100-199 are approvals, 200-299 issuer
// declines, 300-399 gateway declines and 400-499 processor errors.
const (
	ResponseCodeUnknown          = 0
	ResponseCodePendingPayment   = 99
	ResponseCodeApproved         = 100
	ResponseCodeApprovedPending  = 101
	ResponseCodePartialApproval  = 110
	ResponseCodeDeclined         = 200
	ResponseCodeDoNotHonor       = 201
	ResponseCodeInsufficient     = 202
	ResponseCodeExpiredCard      = 223
	ResponseCodeInvalidCVC       = 225
	ResponseCodeGatewayDecline   = 300
	ResponseCodeDuplicate        = 301
	ResponseCodeRuleEngine       = 310
	ResponseCodeProcessorError   = 400
	ResponseCodeMerchantConfig   = 410
	ResponseCodeProcessorOffline = 421
)

// Card entry types.
const (
	EntryTypeKeyed  = "keyed"
	EntryTypeSwiped = "swiped"
)

// ACH SEC codes.
const (
	SecCodeWeb = "web"
	SecCodeCCD = "ccd"
	SecCodePPD = "ppd"
	SecCodeTel = "tel"
)

// ACH account types.
const (
	AccountTypeChecking = "checking"
	AccountTypeSavings  = "savings"
)

// Terminal receipt options for TerminalPayment.PrintReceipt.
const (
	ReceiptNone     = "no"
	ReceiptCustomer = "customer"
	ReceiptMerchant = "merchant"
	ReceiptBoth     = "both"
)

// Payment adjustment types.
const (
	AdjustmentFlat       = "flat"
	AdjustmentPercentage = "percentage"
)

// Transaction sources reported in Transaction.TransactionSource.
const (
	SourceAPI          = "api"
	SourceControlPanel = "cp"
	SourceCart         = "cart"
	SourceInvoice      = "invoice"
	SourceRecurring    = "recurring"
	SourceBatch        = "batch"
)

// TransactionRequest describes a sale, authorization, verification or credit.
// Amounts are integers in cents: 1299 is $12.99. Provide exactly one payment
// method. Only Type, Amount (or BaseAmount) and PaymentMethod are required;
// everything else is optional and omitted from the request when zero.
type TransactionRequest struct {
	// Type is one of the TransactionType* constants. The convenience
	// methods Sale, Authorize, Verify and Credit set it for you.
	Type string `json:"type,omitempty"`
	// Amount is the final amount to charge in cents, including all fees
	// and taxes. Use BaseAmount instead to have the gateway add surcharges.
	Amount int `json:"amount,omitempty"`
	// BaseAmount is the amount before surcharges and fees, which the
	// gateway calculates and adds.
	BaseAmount int `json:"base_amount,omitempty"`
	// Currency is an ISO 4217 code. Defaults to "USD".
	Currency string `json:"currency,omitempty"`
	// ProcessorID selects a processor. Required when the account has no
	// default processor for the payment method.
	ProcessorID string `json:"processor_id,omitempty"`
	// PaymentMethod carries exactly one payment method.
	PaymentMethod PaymentMethod `json:"payment_method"`

	// IdempotencyKey (UUID) makes retries safe. See NewIdempotencyKey.
	IdempotencyKey string `json:"idempotency_key,omitempty"`
	// IdempotencyTime is the key's time to live in seconds (default 300).
	IdempotencyTime int `json:"idempotency_time,omitempty"`

	// OrderID supports up to 17 alphanumeric characters. Required for
	// Level 3 processing.
	OrderID string `json:"order_id,omitempty"`
	// PoNumber supports up to 17 alphanumeric characters.
	PoNumber string `json:"po_number,omitempty"`
	// Description is free text up to 255 characters.
	Description string `json:"description,omitempty"`
	// IPAddress is the end user's IPv4 or IPv6 address without a port.
	IPAddress string `json:"ip_address,omitempty"`
	// VendorID is a processor specific field; only use it when instructed
	// by support.
	VendorID string `json:"vendor_id,omitempty"`

	TaxAmount      int  `json:"tax_amount,omitempty"`
	TaxExempt      bool `json:"tax_exempt,omitempty"`
	ShippingAmount int  `json:"shipping_amount,omitempty"`
	DiscountAmount int  `json:"discount_amount,omitempty"`
	TipAmount      int  `json:"tip_amount,omitempty"`
	// PaymentAdjustment applies a convenience fee, service fee or
	// surcharge.
	PaymentAdjustment *PaymentAdjustment `json:"payment_adjustment,omitempty"`

	BillingAddress  *Address `json:"billing_address,omitempty"`
	ShippingAddress *Address `json:"shipping_address,omitempty"`

	// LineItems itemises the purchase and is required for Level 3.
	LineItems []LineItem `json:"line_items,omitempty"`
	// SummaryCommodityCode and ShipFromPostalCode are required for Level 3.
	SummaryCommodityCode string `json:"summary_commodity_code,omitempty"`
	ShipFromPostalCode   string `json:"ship_from_postal_code,omitempty"`

	// CustomFields maps custom field IDs to their values. Create fields
	// with the Custom Fields API first.
	CustomFields map[string][]string `json:"custom_fields,omitempty"`
	// GroupName is required when the custom fields belong to a non-default
	// group.
	GroupName string `json:"group_name,omitempty"`

	// CreateVaultRecord stores the payment method as a new customer after
	// a successful transaction.
	CreateVaultRecord bool `json:"create_vault_record,omitempty"`
	// CreateVaultRecordFor adds the payment method to an existing
	// customer after a successful transaction.
	CreateVaultRecordFor string `json:"create_vault_record_for,omitempty"`

	// EmailReceipt sends a receipt to EmailAddress.
	EmailReceipt bool   `json:"email_receipt,omitempty"`
	EmailAddress string `json:"email_address,omitempty"`

	// Descriptor customises the text on the cardholder's statement.
	Descriptor *Descriptor `json:"descriptor,omitempty"`
	// AllowPartialPayment lets the processor approve a lower amount.
	AllowPartialPayment bool `json:"allow_partial_payment,omitempty"`
	// SplitTransactionAmount runs a secondary transaction on a second
	// processor (requires the Split Transaction feature).
	SplitTransactionAmount int `json:"split_transaction_amount,omitempty"`

	// Card-on-file (CIT/MIT) indicators.
	CardOnFileIndicator       string `json:"card_on_file_indicator,omitempty"`
	InitiatedBy               string `json:"initiated_by,omitempty"`
	InitialTransactionID      string `json:"initial_transaction_id,omitempty"`
	StoredCredentialIndicator string `json:"stored_credential_indicator,omitempty"`
	BillingMethod             string `json:"billing_method,omitempty"`

	// HSA/FSA processing.
	IIASStatus        string             `json:"iias_status,omitempty"`
	AdditionalAmounts *AdditionalAmounts `json:"additional_amounts,omitempty"`

	// ProcessorSpecific carries processor specific options such as the
	// PaySafe Direct subscription fields.
	ProcessorSpecific map[string]any `json:"processor_specific,omitempty"`
}

// PaymentMethod holds exactly one way to pay.
type PaymentMethod struct {
	Card *CardPayment `json:"card,omitempty"`
	// EMV carries chip data captured alongside a card.
	EMV      *EMVData         `json:"emv,omitempty"`
	ACH      *ACHPayment      `json:"ach,omitempty"`
	Customer *CustomerPayment `json:"customer,omitempty"`
	Terminal *TerminalPayment `json:"terminal,omitempty"`
	// Token is a temporary token from the Tokenizer.
	Token string `json:"token,omitempty"`
	// APM selects an alternative payment method such as Klarna or OXXO.
	APM           *APMPayment    `json:"apm,omitempty"`
	ApplePayToken *ApplePayToken `json:"apple_pay_token,omitempty"`
	// GooglePayToken is the token from
	// paymentData.paymentMethodData.tokenizationData.token, passed through
	// unchanged.
	GooglePayToken json.RawMessage `json:"google_pay_token,omitempty"`
}

// CardPayment is a card entered by number.
type CardPayment struct {
	// EntryType is EntryTypeKeyed or EntryTypeSwiped.
	EntryType string `json:"entry_type,omitempty"`
	// Number is the card number, digits only.
	Number string `json:"number"`
	// ExpirationDate is "MM/YY".
	ExpirationDate string `json:"expiration_date"`
	// CVC is the card verification code. Required if the gateway rule is
	// enabled.
	CVC             string `json:"cvc,omitempty"`
	Track1          string `json:"track_1,omitempty"`
	Track2          string `json:"track_2,omitempty"`
	EncryptedTrack1 string `json:"encrypted_track_1,omitempty"`
	EncryptedTrack2 string `json:"encrypted_track_2,omitempty"`
	KSN             string `json:"ksn,omitempty"`
	// CardholderAuthentication carries 3-D Secure results.
	CardholderAuthentication *CardholderAuthentication `json:"cardholder_authentication,omitempty"`
}

// CardholderAuthentication carries 3-D Secure data collected elsewhere.
type CardholderAuthentication struct {
	ECI              string `json:"eci,omitempty"`
	CAVV             string `json:"cavv,omitempty"`
	XID              string `json:"xid,omitempty"`
	Cryptogram       string `json:"cryptogram,omitempty"`
	Version          string `json:"version,omitempty"`
	DSTransactionID  string `json:"ds_transaction_id,omitempty"`
	ACSTransactionID string `json:"acs_transaction_id,omitempty"`
}

// EMVData carries chip card data.
type EMVData struct {
	// TLVData maps EMV tags to hex encoded values.
	TLVData            map[string]string `json:"tlv_data,omitempty"`
	DeviceSerialNumber string            `json:"device_serial_number,omitempty"`
}

// ACHPayment debits a bank account.
type ACHPayment struct {
	RoutingNumber string `json:"routing_number"`
	AccountNumber string `json:"account_number"`
	// SecCode is one of the SecCode* constants.
	SecCode string `json:"sec_code"`
	// AccountType is AccountTypeChecking or AccountTypeSavings.
	AccountType string `json:"account_type"`
	// CheckNumber is required when SecCode is SecCodeTel.
	CheckNumber string `json:"check_number,omitempty"`
	// FundingSpeed is "standard" (default) or "sameday".
	FundingSpeed                string                       `json:"funding_speed,omitempty"`
	AccountholderAuthentication *AccountholderAuthentication `json:"accountholder_authentication,omitempty"`
}

// AccountholderAuthentication carries driver's licence data some ACH
// processors require.
type AccountholderAuthentication struct {
	DLState  string `json:"dl_state"`
	DLNumber string `json:"dl_number"`
}

// CustomerPayment charges a payment method stored in the Customer Vault.
type CustomerPayment struct {
	// ID is the vault customer ID.
	ID string `json:"id"`
	// PaymentMethodID and PaymentMethodType select a stored payment
	// method. They default to the customer's defaults.
	PaymentMethodID   string `json:"payment_method_id,omitempty"`
	PaymentMethodType string `json:"payment_method_type,omitempty"`
	BillingAddressID  string `json:"billing_address_id,omitempty"`
	ShippingAddressID string `json:"shipping_address_id,omitempty"`
}

// TerminalPayment routes the transaction to a physical terminal.
type TerminalPayment struct {
	ID             string `json:"id"`
	ExpirationDate string `json:"expiration_date,omitempty"`
	CVC            string `json:"cvc,omitempty"`
	// PrintReceipt is one of the Receipt* constants.
	PrintReceipt      string `json:"print_receipt"`
	SignatureRequired bool   `json:"signature_required"`
}

// APMPayment selects an alternative payment method (Klarna, OXXO, Alipay,
// WeChat Pay, DragonPay, SEPA). Fields beyond Type, MerchantRedirectURL,
// Locale and MobileView are method specific.
type APMPayment struct {
	Type                string `json:"type"`
	MerchantRedirectURL string `json:"merchant_redirect_url"`
	Locale              string `json:"locale,omitempty"`
	MobileView          bool   `json:"mobile_view"`
	NationalID          string `json:"national_id,omitempty"`
	ConsumerRef         string `json:"consumer_ref,omitempty"`
	// SEPA fields.
	ConsumerID           string `json:"consumer_id,omitempty"`
	IBAN                 string `json:"iban,omitempty"`
	MandateReference     string `json:"mandate_reference,omitempty"`
	MandateURL           string `json:"mandate_url,omitempty"`
	MandateSignatureDate string `json:"mandate_signature_date,omitempty"`
	// OXXO fields.
	DueDate string `json:"due_date,omitempty"`
	// Klarna fields.
	PaymentMethodCategory string `json:"payment_method_category,omitempty"`
	PurchaseType          string `json:"purchase_type,omitempty"`
	HPPTitle              string `json:"hpp_title,omitempty"`
	LogoURL               string `json:"logo_url,omitempty"`
}

// ApplePayToken pays with Apple Pay, either through a Wallet.js temporary
// token or a native PKPaymentToken. Apple Pay is production only.
type ApplePayToken struct {
	// TemporaryToken is the one-time token issued by Wallet.js.
	TemporaryToken string `json:"temporary_token,omitempty"`
	// KeyID is the registered Apple Pay credential ID (native flows).
	KeyID string `json:"key_id,omitempty"`
	// PKPaymentToken is Apple's encrypted payload, passed through unchanged.
	PKPaymentToken json.RawMessage `json:"pkpaymenttoken,omitempty"`
}

// PaymentAdjustment adds a flat or percentage fee. Value is cents for
// AdjustmentFlat or thousandths of a percent for AdjustmentPercentage
// (1000 is 1.000%).
type PaymentAdjustment struct {
	Type  string `json:"type"`
	Value int    `json:"value"`
}

// Descriptor customises the merchant descriptor on card statements.
// Processor support varies.
type Descriptor struct {
	Name       string `json:"name,omitempty"`
	Address    string `json:"address,omitempty"`
	City       string `json:"city,omitempty"`
	State      string `json:"state,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
}

// AdditionalAmounts carries HSA/FSA breakdowns.
type AdditionalAmounts struct {
	HSA *HSAAmounts `json:"hsa,omitempty"`
}

// HSAAmounts breaks down an HSA/FSA purchase in cents.
type HSAAmounts struct {
	Total        int `json:"total"`
	RxAmount     int `json:"rx_amount,omitempty"`
	VisionAmount int `json:"vision_amount,omitempty"`
	ClinicAmount int `json:"clinic_amount,omitempty"`
	DentalAmount int `json:"dental_amount,omitempty"`
}

// Address is a billing or shipping address on a transaction.
type Address struct {
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
	Company      string `json:"company,omitempty"`
	AddressLine1 string `json:"address_line_1,omitempty"`
	AddressLine2 string `json:"address_line_2,omitempty"`
	City         string `json:"city,omitempty"`
	State        string `json:"state,omitempty"`
	PostalCode   string `json:"postal_code,omitempty"`
	Country      string `json:"country,omitempty"`
	Email        string `json:"email,omitempty"`
	Phone        string `json:"phone,omitempty"`
	Fax          string `json:"fax,omitempty"`
}

// LineItem itemises a purchase. Amounts are in cents and rates are strings
// with three decimals ("10.000" is 10%).
type LineItem struct {
	Name              string  `json:"name,omitempty"`
	Description       string  `json:"description,omitempty"`
	ProductCode       string  `json:"product_code,omitempty"`
	CommodityCode     string  `json:"commodity_code,omitempty"`
	Quantity          float64 `json:"quantity,omitempty"`
	UnitPrice         int     `json:"unit_price,omitempty"`
	Amount            int     `json:"amount,omitempty"`
	TaxAmount         int     `json:"tax_amount,omitempty"`
	TaxRate           string  `json:"tax_rate,omitempty"`
	NationalTaxAmount int     `json:"national_tax_amount,omitempty"`
	NationalTaxRate   string  `json:"national_tax_rate,omitempty"`
	DiscountAmount    int     `json:"discount_amount,omitempty"`
	FreightAmount     int     `json:"freight_amount,omitempty"`
	UnitOfMeasure     string  `json:"unit_of_measure,omitempty"`
	QuantityShipped   float64 `json:"quantity_shipped,omitempty"`
	LocalTax          int     `json:"local_tax,omitempty"`
}

// Transaction is a processed transaction as returned by the gateway. Fields
// that the gateway omits when empty are pointers or slices; everything else
// is always present with a zero default.
type Transaction struct {
	APIResource

	ID                string `json:"id"`
	Type              string `json:"type"`
	Status            string `json:"status"`
	TransactionSource string `json:"transaction_source"`
	// Response is the gateway's verdict, for example "approved" or
	// "declined". Use Approved for a code based check.
	Response     string `json:"response"`
	ResponseCode int    `json:"response_code"`

	Amount            int  `json:"amount"`
	BaseAmount        int  `json:"base_amount"`
	AmountAuthorized  int  `json:"amount_authorized"`
	AmountCaptured    int  `json:"amount_captured"`
	AmountSettled     int  `json:"amount_settled"`
	AmountRefunded    int  `json:"amount_refunded"`
	TaxAmount         int  `json:"tax_amount"`
	TaxExempt         bool `json:"tax_exempt"`
	ShippingAmount    int  `json:"shipping_amount"`
	DiscountAmount    int  `json:"discount_amount"`
	Surcharge         int  `json:"surcharge"`
	ServiceFee        int  `json:"service_fee"`
	TipAmount         int  `json:"tip_amount"`
	PaymentAdjustment int  `json:"payment_adjustment"`
	NationalTaxAmount int  `json:"national_tax_amount"`
	DutyAmount        int  `json:"duty_amount"`

	Currency      string       `json:"currency"`
	PaymentMethod string       `json:"payment_method"`
	PaymentType   string       `json:"payment_type"`
	ResponseBody  ResponseBody `json:"response_body"`

	ProcessorID   string `json:"processor_id"`
	ProcessorType string `json:"processor_type"`
	ProcessorName string `json:"processor_name"`

	MerchantID   string `json:"merchant_id"`
	MerchantName string `json:"merchant_name"`
	UserID       string `json:"user_id"`
	UserName     string `json:"user_name"`

	CustomerID              string `json:"customer_id"`
	CustomerPaymentType     string `json:"customer_payment_type"`
	CustomerPaymentID       string `json:"customer_payment_id"`
	SubscriptionID          string `json:"subscription_id"`
	ReferencedTransactionID string `json:"referenced_transaction_id"`
	SettlementBatchID       string `json:"settlement_batch_id"`

	OrderID      string `json:"order_id"`
	PoNumber     string `json:"po_number"`
	Description  string `json:"description"`
	IPAddress    string `json:"ip_address"`
	EmailReceipt bool   `json:"email_receipt"`
	EmailAddress string `json:"email_address"`

	IdempotencyKey  string `json:"idempotency_key"`
	IdempotencyTime int    `json:"idempotency_time"`

	Features     []string            `json:"features"`
	CustomFields map[string][]string `json:"custom_fields"`
	LineItems    []LineItem          `json:"line_items"`

	BillingAddress  Address `json:"billing_address"`
	ShippingAddress Address `json:"shipping_address"`

	SummaryCommodityCode          string `json:"summary_commodity_code"`
	ShipFromPostalCode            string `json:"ship_from_postal_code"`
	MerchantVATRegistrationNumber string `json:"merchant_vat_registration_number"`
	CustomerVATRegistrationNumber string `json:"customer_vat_registration_number"`
	ReceiptData                   string `json:"receipt_data"`

	// SplitTransactionResponse describes the secondary transaction of a
	// split transaction, when one was processed.
	SplitTransactionResponse json.RawMessage `json:"split_transaction_response,omitempty"`

	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	CapturedAt *time.Time `json:"captured_at"`
	SettledAt  *time.Time `json:"settled_at"`
}

// Approved reports whether the transaction was approved, including partial
// approvals (response codes 100-199).
func (t *Transaction) Approved() bool {
	return t.ResponseCode >= 100 && t.ResponseCode < 200
}

// PartiallyApproved reports whether the processor approved less than the
// requested amount. Compare AmountAuthorized with Amount to see how much.
func (t *Transaction) PartiallyApproved() bool {
	return t.ResponseCode == ResponseCodePartialApproval ||
		(t.Approved() && t.AmountAuthorized > 0 && t.AmountAuthorized < t.Amount)
}

// Declined reports whether the issuer declined (codes 200-299).
func (t *Transaction) Declined() bool {
	return t.ResponseCode >= 200 && t.ResponseCode < 300
}

// GatewayDeclined reports whether the gateway declined for configuration or
// fraud reasons (codes 300-399).
func (t *Transaction) GatewayDeclined() bool {
	return t.ResponseCode >= 300 && t.ResponseCode < 400
}

// ProcessorError reports whether the processor returned an error (codes
// 400-499).
func (t *Transaction) ProcessorError() bool {
	return t.ResponseCode >= 400 && t.ResponseCode < 500
}

// ResponseBody holds the processor response for whichever payment method
// was used. Only one member is populated.
type ResponseBody struct {
	Card     *CardResponse     `json:"card,omitempty"`
	ACH      *ACHResponse      `json:"ach,omitempty"`
	Terminal *TerminalResponse `json:"terminal,omitempty"`
	APM      json.RawMessage   `json:"apm,omitempty"`
	Cash     json.RawMessage   `json:"cash,omitempty"`
}

// CardResponse is the processor response for a card payment.
type CardResponse struct {
	ID                     string `json:"id"`
	CardType               string `json:"card_type"`
	FirstSix               string `json:"first_six"`
	LastFour               string `json:"last_four"`
	MaskedCard             string `json:"masked_card"`
	ExpirationDate         string `json:"expiration_date"`
	Response               string `json:"response"`
	ResponseCode           int    `json:"response_code"`
	AuthCode               string `json:"auth_code"`
	ProcessorResponseCode  string `json:"processor_response_code"`
	ProcessorResponseText  string `json:"processor_response_text"`
	ProcessorTransactionID string `json:"processor_transaction_id"`
	ProcessorType          string `json:"processor_type"`
	ProcessorID            string `json:"processor_id"`
	BINType                string `json:"bin_type"`
	// Type is "credit" or "debit".
	Type            string `json:"type"`
	AVSResponseCode string `json:"avs_response_code"`
	CVVResponseCode string `json:"cvv_response_code"`
	// ProcessorSpecific varies by processor; it may be an object or an
	// empty string.
	ProcessorSpecific json.RawMessage `json:"processor_specific,omitempty"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

// ACHResponse is the processor response for an ACH payment.
type ACHResponse struct {
	ID                    string          `json:"id"`
	SecCode               string          `json:"sec_code"`
	AccountType           string          `json:"account_type"`
	MaskedAccountNumber   string          `json:"masked_account_number"`
	RoutingNumber         string          `json:"routing_number"`
	Response              string          `json:"response"`
	ResponseCode          int             `json:"response_code"`
	AuthCode              string          `json:"auth_code"`
	ProcessorResponseCode string          `json:"processor_response_code"`
	ProcessorResponseText string          `json:"processor_response_text"`
	ProcessorType         string          `json:"processor_type"`
	ProcessorID           string          `json:"processor_id"`
	ProcessorSpecific     json.RawMessage `json:"processor_specific,omitempty"`
	CreatedAt             time.Time       `json:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at"`
}

// TerminalResponse is the processor response for a terminal payment.
type TerminalResponse struct {
	ID                    string          `json:"id"`
	TerminalID            string          `json:"terminal_id"`
	TerminalDescription   string          `json:"terminal_description"`
	CardType              string          `json:"card_type"`
	PaymentType           string          `json:"payment_type"`
	EntryType             string          `json:"entry_type"`
	FirstFour             string          `json:"first_four"`
	LastFour              string          `json:"last_four"`
	MaskedCard            string          `json:"masked_card"`
	CardholderName        string          `json:"cardholder_name"`
	AuthCode              string          `json:"auth_code"`
	ResponseCode          int             `json:"response_code"`
	ProcessorResponseText string          `json:"processor_response_text"`
	ProcessorSpecific     json.RawMessage `json:"processor_specific,omitempty"`
	EMVAID                string          `json:"emv_aid"`
	EMVAppName            string          `json:"emv_app_name"`
	EMVTVR                string          `json:"emv_tvr"`
	EMVTSI                string          `json:"emv_tsi"`
	SignatureData         string          `json:"signature_data"`
	CreatedAt             time.Time       `json:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at"`
}

// TransactionList is a page of transactions.
type TransactionList struct {
	APIResource
	Data []*Transaction
	// TotalCount is the number of matching records across all pages.
	TotalCount int
}

// TransactionSearchRequest filters transaction history. Every field is
// optional; use the query constructors (Equal, IntEqual, DateRange, ...) to
// populate them. Without a CreatedAt range the gateway searches the prior
// four months.
type TransactionSearchRequest struct {
	TransactionID     *StringQuery    `json:"transaction_id,omitempty"`
	UserID            *StringQuery    `json:"user_id,omitempty"`
	Type              *StringQuery    `json:"type,omitempty"`
	IPAddress         *StringQuery    `json:"ip_address,omitempty"`
	Amount            *IntQuery       `json:"amount,omitempty"`
	AmountAuthorized  *IntQuery       `json:"amount_authorized,omitempty"`
	AmountCaptured    *IntQuery       `json:"amount_captured,omitempty"`
	AmountSettled     *IntQuery       `json:"amount_settled,omitempty"`
	TaxAmount         *IntQuery       `json:"tax_amount,omitempty"`
	PoNumber          *StringQuery    `json:"po_number,omitempty"`
	OrderID           *StringQuery    `json:"order_id,omitempty"`
	PaymentMethod     *StringQuery    `json:"payment_method,omitempty"`
	PaymentType       *StringQuery    `json:"payment_type,omitempty"`
	Status            *StringQuery    `json:"status,omitempty"`
	ProcessorID       *StringQuery    `json:"processor_id,omitempty"`
	CustomerID        *StringQuery    `json:"customer_id,omitempty"`
	SettlementBatchID *StringQuery    `json:"settlement_batch_id,omitempty"`
	CreatedAt         *DateRangeQuery `json:"created_at,omitempty"`
	CapturedAt        *DateRangeQuery `json:"captured_at,omitempty"`
	SettledAt         *DateRangeQuery `json:"settled_at,omitempty"`
	BillingAddress    *AddressQuery   `json:"billing_address,omitempty"`
	ShippingAddress   *AddressQuery   `json:"shipping_address,omitempty"`
	// Limit caps the page size (0-100). Offset skips records.
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

// AddressQuery filters on address fields.
type AddressQuery struct {
	AddressID    *StringQuery `json:"address_id,omitempty"`
	FirstName    *StringQuery `json:"first_name,omitempty"`
	LastName     *StringQuery `json:"last_name,omitempty"`
	Company      *StringQuery `json:"company,omitempty"`
	AddressLine1 *StringQuery `json:"address_line_1,omitempty"`
	AddressLine2 *StringQuery `json:"address_line_2,omitempty"`
	City         *StringQuery `json:"city,omitempty"`
	State        *StringQuery `json:"state,omitempty"`
	PostalCode   *StringQuery `json:"postal_code,omitempty"`
	Country      *StringQuery `json:"country,omitempty"`
	Email        *StringQuery `json:"email,omitempty"`
	Phone        *StringQuery `json:"phone,omitempty"`
	Fax          *StringQuery `json:"fax,omitempty"`
}

// CaptureRequest captures a previously authorized transaction. Every field
// is optional and defaults to the value from the authorization.
type CaptureRequest struct {
	Amount         int    `json:"amount,omitempty"`
	TaxAmount      int    `json:"tax_amount,omitempty"`
	TaxExempt      bool   `json:"tax_exempt,omitempty"`
	ShippingAmount int    `json:"shipping_amount,omitempty"`
	OrderID        string `json:"order_id,omitempty"`
	PoNumber       string `json:"po_number,omitempty"`
	IPAddress      string `json:"ip_address,omitempty"`
}

// RefundRequest refunds a settled transaction. Amount is required for a
// partial refund and defaults to the full settled amount when zero.
type RefundRequest struct {
	Amount    int `json:"amount,omitempty"`
	Surcharge int `json:"surcharge,omitempty"`
}

// Create processes a transaction of the type set in req.Type. The
// convenience wrappers Sale, Authorize, Verify and Credit set the type for
// you.
//
// A declined card is not an error: check Transaction.Approved on the result.
func (s *TransactionsService) Create(ctx context.Context, req *TransactionRequest) (*Transaction, error) {
	if req == nil {
		return nil, errNilRequest("transaction")
	}
	return doOne[Transaction](ctx, s.client, http.MethodPost, "transaction", nil, req)
}

// Sale authorizes and captures in a single call.
func (s *TransactionsService) Sale(ctx context.Context, req *TransactionRequest) (*Transaction, error) {
	return s.createTyped(ctx, TransactionTypeSale, req)
}

// Authorize places a hold that Capture completes later.
func (s *TransactionsService) Authorize(ctx context.Context, req *TransactionRequest) (*Transaction, error) {
	return s.createTyped(ctx, TransactionTypeAuthorize, req)
}

// Verify runs a zero-dollar verification of the payment method.
func (s *TransactionsService) Verify(ctx context.Context, req *TransactionRequest) (*Transaction, error) {
	return s.createTyped(ctx, TransactionTypeVerification, req)
}

// Credit sends funds to the payment method without a prior sale.
func (s *TransactionsService) Credit(ctx context.Context, req *TransactionRequest) (*Transaction, error) {
	return s.createTyped(ctx, TransactionTypeCredit, req)
}

func (s *TransactionsService) createTyped(ctx context.Context, typ string, req *TransactionRequest) (*Transaction, error) {
	if req == nil {
		return nil, errNilRequest("transaction")
	}
	copied := *req
	copied.Type = typ
	return s.Create(ctx, &copied)
}

// Get retrieves a transaction by ID.
func (s *TransactionsService) Get(ctx context.Context, transactionID string) (*Transaction, error) {
	if err := requireID("transaction id", transactionID); err != nil {
		return nil, err
	}
	return doOne[Transaction](ctx, s.client, http.MethodGet, joinPath("transaction", transactionID), nil, nil)
}

// Search returns transactions matching req. Pass an empty request to list
// the prior four months.
func (s *TransactionsService) Search(ctx context.Context, req *TransactionSearchRequest) (*TransactionList, error) {
	if req == nil {
		req = &TransactionSearchRequest{}
	}
	items, resp, err := doList[Transaction](ctx, s.client, http.MethodPost, "transaction/search", nil, req)
	if err != nil {
		return nil, err
	}
	return &TransactionList{APIResource: APIResource{LastResponse: resp}, Data: items, TotalCount: resp.TotalCount}, nil
}

// Capture captures funds for an authorized transaction. A nil req captures
// the full authorized amount.
func (s *TransactionsService) Capture(ctx context.Context, transactionID string, req *CaptureRequest) (*Transaction, error) {
	if err := requireID("transaction id", transactionID); err != nil {
		return nil, err
	}
	if req == nil {
		req = &CaptureRequest{}
	}
	return doOne[Transaction](ctx, s.client, http.MethodPost, joinPath("transaction", transactionID, "capture"), nil, req)
}

// Void cancels a transaction that is pending settlement. Where the
// processor supports it, the void is sent as an authorization reversal.
func (s *TransactionsService) Void(ctx context.Context, transactionID string) error {
	if err := requireID("transaction id", transactionID); err != nil {
		return err
	}
	return doEmpty(ctx, s.client, http.MethodPost, joinPath("transaction", transactionID, "void"), nil, nil)
}

// Refund returns funds from a settled transaction. Multiple partial refunds
// are allowed up to the settled amount. A nil req refunds in full.
func (s *TransactionsService) Refund(ctx context.Context, transactionID string, req *RefundRequest) (*Transaction, error) {
	if err := requireID("transaction id", transactionID); err != nil {
		return nil, err
	}
	if req == nil {
		req = &RefundRequest{}
	}
	return doOne[Transaction](ctx, s.client, http.MethodPost, joinPath("transaction", transactionID, "refund"), nil, req)
}
