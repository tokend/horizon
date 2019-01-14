package resource

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/horizon/resource/base"
)

// Asset - resource object representing AssetEntry
type Asset struct {
	Key
	Attributes    AssetAttributes    `json:"attributes"`
	Relationships AssetRelationships `json:"relationships"`
}

func NewAsset(core *core2.Asset) Asset {
	return Asset{
		Key: Key{
			ID:   core.Code,
			Type: typeAssets,
		},
		Attributes: AssetAttributes{
			PreIssuanceAssetSigner: core.PreIssuanceAssetSigner,
			Details:                core.Details,
			MaxIssuanceAmount:      amount.String(core.MaxIssuanceAmount),
			AvailableForIssuance:   amount.String(core.AvailableForIssuance),
			Issued:                 amount.String(core.Issued),
			PendingIssuance:        amount.String(core.PendingIssuance),
			Policies: Mask{
				Mask:  core.Policies,
				Flags: base.FlagFromXdrAssetPolicy(core.Policies, xdr.AssetPolicyAll),
			},
			TrailingDigits: core.TrailingDigits,
		},
		Relationships: AssetRelationships{
			Owner: Key{
				ID:   core.Owner,
				Type: typeAccounts,
			},
		},
	}
}

//AssetAttributes - represents info about asset
type AssetAttributes struct {
	PreIssuanceAssetSigner string                 `json:"pre_issuance_asset_signer"`
	Details                map[string]interface{} `json:"details"`
	MaxIssuanceAmount      string                 `json:"max_issuance_amount"`
	AvailableForIssuance   string                 `json:"available_for_issuance"`
	Issued                 string                 `json:"issued"`
	PendingIssuance        string                 `json:"pending_issuance"`
	Policies               Mask                   `json:"policies"`
	TrailingDigits         int64                  `json:"trailing_digits"`
}

// AssetRelationships - represents references from account to other resource objects
type AssetRelationships struct {
	Owner Key `json:"owner,omitempty"`
}
