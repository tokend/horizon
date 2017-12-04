package core

type Asset struct {
	Code                 string `db:"code"`
	Policies             int32  `db:"policies"`
	Owner                string `db:"owner"`
	AvailableForIssuance uint64 `db:"available_for_issueance"`
	Name                 string `db:"name"`
	PreissuedAssetSigner string `db:"preissued_asset_signer"`
	Description          string `db:"description"`
	ExternalResourceLink string `db:"external_resource_link"`
	MaxIssuanceAmount    uint64 `db:"max_issuance_amount"`
	Issued               uint64 `db:"issued"`
}

func (a *Asset) IsVisibleForUser(account *Account) bool {
	return true
}
