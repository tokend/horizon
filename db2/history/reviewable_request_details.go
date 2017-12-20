package history

type AssetCreationRequest struct {
	Asset                string `json:"asset"`
	Description          string `json:"description"`
	ExternalResourceLink string `json:"external_resource_link"`
	Policies             int32  `json:"policies"`
	Name                 string `json:"name"`
	PreIssuedAssetSigner string `json:"pre_issued_asset_signer"`
	MaxIssuanceAmount    string `json:"max_issuance_amount"`
	LogoID               string `json:"logo_id"`
}

type AssetUpdateRequest struct {
	Asset                string `json:"asset"`
	Description          string `json:"description"`
	ExternalResourceLink string `json:"external_resource_link"`
	Policies             int32  `json:"policies"`
	LogoID               string `json:"logo_id"`
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
	ExternalDetails string `json:"external_details"`
}

type WithdrawalRequest struct {
	BalanceID       string `json:"balance_id"`
	Amount          string `json:"amount"`
	FixedFee        string `json:"fixed_fee"`
	PercentFee      string `json:"percent_fee"`
	ExternalDetails string `json:"external_details"`
	DestAssetCode   string `json:"dest_asset_code"`
	DestAssetAmount string `json:"dest_asset_amount"`
}
