package history

import "time"

type AssetCreationRequest struct {
	Asset                string `json:"asset"`
	Policies             int32  `json:"policies"`
	PreIssuedAssetSigner string `json:"pre_issued_asset_signer"`
	MaxIssuanceAmount    string `json:"max_issuance_amount"`
	Details map[string]interface{} `json:"details"`
}

type AssetUpdateRequest struct {
	Asset                string `json:"asset"`
	Policies             int32  `json:"policies"`
	Details map[string]interface{} `json:"details"`
}

type PreIssuanceRequest struct {
	Asset     string `json:"asset"`
	Amount    string `json:"amount"`
	Signature string `json:"signature"`
	Reference string `json:"reference"`
}

type IssuanceRequest struct {
	Asset    string `json:"asset"`
	Amount   string `json:"amount"`
	Receiver string `json:"receiver"`
}

type WithdrawalRequest struct {
	BalanceID       string                 `json:"balance_id"`
	Amount          string                 `json:"amount"`
	FixedFee        string                 `json:"fixed_fee"`
	PercentFee      string                 `json:"percent_fee"`
	ExternalDetails string                 `json:"external_details"`
	DestAssetCode   string                 `json:"dest_asset_code"`
	DestAssetAmount string                 `json:"dest_asset_amount"`
	ReviewerDetails map[string]interface{} `json:"reviewer_details"`
}

type SaleRequest struct {
	BaseAsset  string                 `json:"base_asset"`
	QuoteAsset string                 `json:"quote_asset"`
	StartTime  time.Time              `json:"start_time"`
	EndTime    time.Time              `json:"end_time"`
	Price      string                 `json:"price"`
	SoftCap    string                 `json:"soft_cap"`
	HardCap    string                 `json:"hard_cap"`
	Details    map[string]interface{} `json:"details"`
}
