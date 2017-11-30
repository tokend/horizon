package history

type AssetCreationRequest struct {
	Code                 string `json:"code"`
	Description          string `json:"description"`
	ExternalResourceLink string `json:"external_resource_link"`
	Policies             uint32 `json:"policies"`
	Name                 string `json:"name"`
	PreIssuedAssetSigner string `json:"pre_issued_asset_signer"`
	MaxIssuanceAmount    string `json:"max_issuance_amount"`
}

type AssetUpdateRequest struct {
	Code                 string `json:"code"`
	Description          string `json:"description"`
	ExternalResourceLink string `json:"external_resource_link"`
	Policies             uint32 `json:"policies"`
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
