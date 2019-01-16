package regources

// Asset - resource object representing AssetEntry
type Asset struct {
	ID                     string                 `jsonapi:"primary,assets"`
	PreIssuanceAssetSigner string                 `jsonapi:"attr,pre_issuance_asset_signer" `
	Details                map[string]interface{} `jsonapi:"attr,details"`
	MaxIssuanceAmount      Amount                 `jsonapi:"attr,max_issuance_amount"`
	AvailableForIssuance   Amount                 `jsonapi:"attr,available_for_issuance"`
	Issued                 Amount                 `jsonapi:"attr,issued"`
	PendingIssuance        Amount                 `jsonapi:"attr,pending_issuance"`
	Policies               Mask                   `jsonapi:"attr,policies"`
	TrailingDigits         int64                  `jsonapi:"attr,trailing_digits"`
	Owner                  *Account               `jsonapi:"relation,owner,omitempty"`
}
