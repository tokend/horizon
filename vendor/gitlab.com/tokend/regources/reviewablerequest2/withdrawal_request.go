package reviewablerequest2

type WithdrawalRequest struct {
	BalanceID              string                 `json:"balance_id"`
	Amount                 string                 `json:"amount"`
	FixedFee               string                 `json:"fixed_fee"`
	PercentFee             string                 `json:"percent_fee"`
	PreConfirmationDetails map[string]interface{} `json:"pre_confirmation_details"`
	ExternalDetails        map[string]interface{} `json:"external_details"`
	DestAssetCode          string                 `json:"dest_asset_code"`
	DestAssetAmount        string                 `json:"dest_asset_amount"`
	ReviewerDetails        map[string]interface{} `json:"reviewer_details"`
}
