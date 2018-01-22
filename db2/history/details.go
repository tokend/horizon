package history

import "gitlab.com/swarmfund/go/xdr"

type OperationDetails struct {
	Type                    xdr.OperationType               `json:"type"`
	CreateAccount           *CreateAccountDetails           `json:"create_account,omitempty"`
	Payment                 *PaymentDetails                 `json:"payment,omitempty"`
	SetOptions              *SetOptionsDetails              `json:"set_options,omitempty"`
	SetFees                 *SetFeesDetails                 `json:"set_fees,omitempty"`
	ManageAccount           *ManageAccountDetails           `json:"manage_account,omitempty"`
	CreateWithdrawalRequest *CreateWithdrawalRequestDetails `json:"create_withdrawal_request,omitempty"`
	ManageBalance           *ManageBalanceDetails           `json:"manage_balance,omitempty"`
	ReviewPaymentRequest    *ReviewPaymentRequestDetails    `json:"review_payment_request,omitempty"`
	SetLimits               *SetLimitsDetails               `json:"set_limits,omitempty"`
	DirectDebit             *DirectDebitDetails             `json:"direct_debit,omitempty"`
	ManageInvoice           *ManageInvoiceDetails           `json:"manage_invoice,omitempty"`
	ManageAsset             *ManageAssetDetails             `json:"manage_asset,omitempty"`
	ManagerOffer            *ManagerOfferDetails            `json:"manager_offer,omitempty"`
	ManageAssetPair         *ManageAssetPairDetails         `json:"manage_asset_pair,omitempty"`
	CreateIssuanceRequest   *CreateIssuanceRequestDetails   `json:"create_issuance_request,omitempty"`
	ReviewRequest           *ReviewRequestDetails           `json:"review_request,omitempty"`
}

type CreateAccountDetails struct {
	Funder      string `json:"funder,omitempty"`
	Account     string `json:"account,omitempty"`
	AccountType int32  `json:"account_type"`
}

type BasePayment struct {
	From                  string `json:"from,omitempty"`
	To                    string `json:"to,omitempty"`
	FromBalance           string `json:"from_balance,omitempty"`
	ToBalance             string `json:"to_balance,omitempty"`
	Amount                string `json:"amount"`
	Asset                 string `json:"asset"`
	SourcePaymentFee      string `json:"source_payment_fee"`
	DestinationPaymentFee string `json:"destination_payment_fee"`
	SourceFixedFee        string `json:"source_fixed_fee"`
	DestinationFixedFee   string `json:"destination_fixed_fee"`
	SourcePaysForDest     bool   `json:"source_pays_for_dest"`
}

// Payment is the json resource representing a single operation whose type is
// Payment.
type PaymentDetails struct {
	BasePayment
	Subject    string `json:"subject,omitempty"`
	Reference  string `json:"reference,omitempty"`
	QuoteAsset string `json:"qasset"`
}
type SetOptionsDetails struct {
	HomeDomain    string `json:"home_domain,omitempty"`
	InflationDest string `json:"inflation_dest,omitempty"`

	MasterKeyWeight uint32 `json:"master_key_weight"`
	SignerKey       string `json:"signer_key,omitempty"`
	SignerWeight    uint32 `json:"signer_weight,omitempty"`
	SignerType      uint32 `json:"signer_type,omitempty"`
	SignerIdentity  uint32 `json:"signer_identity,omitempty"`

	SetFlags    []int    `json:"set_flags,omitempty"`
	SetFlagsS   []string `json:"set_flags_s,omitempty"`
	ClearFlags  []int    `json:"clear_flags,omitempty"`
	ClearFlagsS []string `json:"clear_flags_s,omitempty"`

	LowThreshold                    uint32 `json:"low_threshold,omitempty"`
	MedThreshold                    uint32 `json:"med_threshold,omitempty"`
	HighThreshold                   uint32 `json:"high_threshold,omitempty"`
	LimitsUpdateRequestDocumentHash string `json:"limits_update_request_document_hash,omitempty"`
}

type FeeDetails struct {
	AssetCode   string `json:"asset_code"`
	FixedFee    string `json:"fixed_fee"`
	PercentFee  string `json:"percent_fee"`
	FeeType     int64  `json:"fee_type"`
	AccountID   string `json:"account_id,omitempty"`
	AccountType int32  `json:"account_type"`
	Subtype     int64  `json:"subtype"`
	LowerBound  int64  `json:"lower_bound"`
	UpperBound  int64  `json:"upper_bound"`
}

type SetFeesDetails struct {
	Fee *FeeDetails `json:"fee"`
}

type ManageAccountDetails struct {
	Account              string `json:"account,omitempty"`
	BlockReasonsToAdd    uint32 `json:"block_reasons_to_add,omitempty"`
	BlockReasonsToRemove uint32 `json:"block_reasons_to_remove,omitempty"`
}

type CreateWithdrawalRequestDetails struct {
	Amount          string                 `json:"amount"`
	Balance         string                 `json:"balance"`
	FeeFixed        string                 `json:"fee_fixed"`
	FeePercent      string                 `json:"fee_percent"`
	ExternalDetails map[string]interface{} `json:"external_details"`
	DestAsset       string                 `json:"dest_asset"`
	DestAmount      string                 `json:"dest_amount"`
}

type ManageBalanceDetails struct {
	Destination string `json:"destination"`
	Action      int32  `json:"action"`
}

type ReviewPaymentRequestDetails struct {
	PaymentID    int64  `json:"payment_id"`
	Accept       bool   `json:"accept"`
	RejectReason string `json:"reject_reason"`
}
type SetLimitsDetails struct{}

type ManageInvoiceDetails struct {
	Amount          string  `json:"amount"`
	ReceiverBalance string  `json:"receiver_balance,omitempty"`
	Sender          string  `json:"sender,omitempty"`
	InvoiceID       uint64  `json:"invoice_id"`
	RejectReason    *string `json:"reject_reason,omitempty"`
	Asset           string  `json:"asset"`
}

type ManagerOfferDetails struct {
	IsBuy     bool   `json:"is_buy"`
	Amount    string `json:"amount"`
	Price     string `json:"price"`
	Fee       string `json:"fee"`
	OfferId   uint64 `json:"offer_id"`
	IsDeleted bool   `json:"is_deleted"`
}

type ManageAssetPairDetails struct {
	BaseAsset               string `json:"base_asset"`
	QuoteAsset              string `json:"quote_asset"`
	PhysicalPrice           string `json:"physical_price"`
	PhysicalPriceCorrection string `json:"physical_price_correction"`
	MaxPriceStep            string `json:"max_price_step"`
	Policies                int32  `json:"policies_i"`
}

type CreateIssuanceRequestDetails struct {
	Reference       string                 `json:"reference"`
	Amount          string                 `json:"amount"`
	Asset           string                 `json:"asset"`
	FeeFixed        string                 `json:"fee_fixed"`
	FeePercent      string                 `json:"fee_percent"`
	ExternalDetails map[string]interface{} `json:"external_details"`
	BalanceID       string                 `json:"balance_id"`
}

type DirectDebitDetails struct {
	From                  string `json:"from"`
	To                    string `json:"to"`
	FromBalance           string `json:"from_balance"`
	ToBalance             string `json:"to_balance"`
	Amount                string `json:"amount"`
	SourcePaymentFee      string `json:"source_payment_fee"`
	DestinationPaymentFee string `json:"destination_payment_fee"`
	SourceFixedFee        string `json:"source_fixed_fee"`
	DestinationFixedFee   string `json:"destination_fixed_fee"`
	SourcePaysForDest     bool   `json:"source_pays_for_dest"`
	Subject               string `json:"subject"`
	Reference             string `json:"reference"`
	AssetCode             string `json:"asset"`
}

type ManageAssetDetails struct {
	RequestID uint64 `json:"request_id"`
	Action    int32  `json:"action"`
}

type ReviewRequestDetails struct {
	Action      int32  `json:"action"`
	Reason      string `json:"reason"`
	RequestHash string `json:"request_hash"`
	RequestID   uint64 `json:"request_id"`
	RequestType int32  `json:"request_type"`
}
