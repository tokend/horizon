package regources

// AssetCreationRequest - represents details of the `asset create` reviewable request
type AssetCreateRequest struct {
	Key
	Attributes AssetCreateRequestAttrs `json:"attributes"`
}

// AssetCreationRequestAttrs - attributes of the `asset create` reviewable request
type AssetCreateRequestAttrs struct {
	Asset                  string  `json:"asset"`
	Policies               int32   `json:"policies"`
	PreIssuanceAssetSigner string  `json:"pre_issuance_asset_signer"`
	MaxIssuanceAmount      string  `json:"max_issuance_amount"`
	InitialPreissuedAmount string  `json:"initial_preissued_amount"`
	CreatorDetails         Details `json:"creator_details"`
}
