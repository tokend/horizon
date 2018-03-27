package history

import (
	"time"

	"gitlab.com/swarmfund/go/xdr"
)

type AssetCreationRequest struct {
	Asset                  string                 `json:"asset"`
	Policies               int32                  `json:"policies"`
	PreIssuedAssetSigner   string                 `json:"pre_issued_asset_signer"`
	MaxIssuanceAmount      string                 `json:"max_issuance_amount"`
	InitialPreissuedAmount string                 `json:"initial_preissued_amount"`
	Details                map[string]interface{} `json:"details"`
}

type AssetUpdateRequest struct {
	Asset    string                 `json:"asset"`
	Policies int32                  `json:"policies"`
	Details  map[string]interface{} `json:"details"`
}

type PreIssuanceRequest struct {
	Asset     string `json:"asset"`
	Amount    string `json:"amount"`
	Signature string `json:"signature"`
	Reference string `json:"reference"`
}

type IssuanceRequest struct {
	Asset           string                 `json:"asset"`
	Amount          string                 `json:"amount"`
	Receiver        string                 `json:"receiver"`
	ExternalDetails map[string]interface{} `json:"external_details"`
}

type WithdrawalRequest struct {
	BalanceID              string                 `json:"balance_id"`
	Amount                 string                 `json:"amount"`
	FixedFee               string                 `json:"fixed_fee"`
	PercentFee             string                 `json:"percent_fee"`
	ExternalDetails        map[string]interface{} `json:"external_details"`
	DestAssetCode          string                 `json:"dest_asset_code"`
	DestAssetAmount        string                 `json:"dest_asset_amount"`
	ReviewerDetails        map[string]interface{} `json:"reviewer_details"`
	PreConfirmationDetails map[string]interface{} `json:"pre_confirmation_details"`
}

type SaleRequest struct {
	BaseAsset         string                 `json:"base_asset"`
	DefaultQuoteAsset string                 `json:"quote_asset"`
	StartTime         time.Time              `json:"start_time"`
	EndTime           time.Time              `json:"end_time"`
	SoftCap           string                 `json:"soft_cap"`
	HardCap           string                 `json:"hard_cap"`
	Details           map[string]interface{} `json:"details"`
	QuoteAssets       []SaleQuoteAsset       `json:"quote_assets"`
	SaleType          xdr.SaleType           `json:"sale_type"`
}

type SaleQuoteAsset struct {
	QuoteAsset string `json:"quote_asset"`
	Price      string `json:"price"`
}

type LimitsUpdateRequest struct {
	DocumentHash string `json:"document_hash"`
}

type AmlAlertRequest struct {
	BalanceID string `json:"balance_id"`
	Amount    string `json:"amount"`
	Reason    string `json:"reason"`
}

type UpdateKYCRequest struct {
	AccountToUpdateKYC string                   `json:"updated_account_id"`
	AccountTypeToSet   xdr.AccountType          `json:"account_type_to_set"`
	KYCLevel           uint32                   `json:"kyc_level"`
	KYCData            map[string]interface{}   `json:"kyc_data"`
	AllTasks           uint32                   `json:"all_tasks"`
	PendingTasks       uint32                   `json:"pending_tasks"`
	SequenceNumber     uint32                   `json:"sequence_number"`
	ExternalDetails    []map[string]interface{} `json:"extrenal_details"`
}
