package resources

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/regources/v2"
)

//NewAsset - creates new instance of Asset from provided one. Returns nil if record is nil
func NewAsset(record *core2.Asset) *regources.Asset {
	if record == nil {
		return nil
	}

	return &regources.Asset{
		ID: record.Code,
		PreIssuanceAssetSigner: record.PreIssuanceAssetSigner,
		Details:                record.Details,
		MaxIssuanceAmount:      regources.Amount(record.MaxIssuanceAmount),
		AvailableForIssuance:   regources.Amount(record.AvailableForIssuance),
		Issued:                 regources.Amount(record.Issued),
		PendingIssuance:        regources.Amount(record.PendingIssuance),
		Policies: regources.Mask{
			Mask:  record.Policies,
			Flags: FlagFromXdrAssetPolicy(record.Policies, xdr.AssetPolicyAll),
		},
		TrailingDigits: record.TrailingDigits,
	}
}
