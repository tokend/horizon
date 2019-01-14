package resource

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/horizon/resource/base"
)

// Asset - resource object representing AssetEntry
type Asset struct {
	ID                     string                 `jsonapi:"primary, assets"`
	PreIssuanceAssetSigner string                 `jsonapi:"attr,pre_issuance_asset_signer" `
	Details                map[string]interface{} `jsonapi:"attr,details"`
	MaxIssuanceAmount      string                 `jsonapi:"attr,max_issuance_amount"`
	AvailableForIssuance   string                 `jsonapi:"attr,available_for_issuance"`
	Issued                 string                 `jsonapi:"attr,issued"`
	PendingIssuance        string                 `jsonapi:"attr,pending_issuance"`
	Policies               Mask                   `jsonapi:"attr,policies"`
	TrailingDigits         int64                  `jsonapi:"attr,trailing_digits"`
	Owner                  *Account               `jsonapi:"relation,owner"`
}

func NewAsset(core *core2.Asset) *Asset {
	if core == nil {
		return nil
	}
	return &Asset{
		ID: core.Code,
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
	}
}
