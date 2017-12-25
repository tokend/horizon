package core

import "encoding/json"

type Asset struct {
	Code                 string `db:"code"`
	Policies             int32  `db:"policies"`
	Owner                string `db:"owner"`
	AvailableForIssuance uint64 `db:"available_for_issueance"`
	PreissuedAssetSigner string `db:"preissued_asset_signer"`
	MaxIssuanceAmount    uint64 `db:"max_issuance_amount"`
	Issued               uint64 `db:"issued"`
	LockedIssuance       uint64 `db:"locked_issuance"`
	Details              []byte `db:"details"`
}

type AssetDetails struct {
	Name                 string `json:"name"`
	Description          string `json:"description"`
	ExternalResourceLink string `json:"external_resource_link"`
	LogoID               string `json:"logo_id"`
}

func (a Asset) GetDetails() AssetDetails {
	var result AssetDetails
	// error is ignored on purpose
	_ = json.Unmarshal(a.Details, &result)
	return result
}
