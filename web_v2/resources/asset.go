package resources

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/regources/rgenerated"
)

//NewAsset - creates new instance of Asset from provided one.
func NewAsset(record core2.Asset) rgenerated.Asset {
	return rgenerated.Asset{
		Key: rgenerated.Key{
			ID:   record.Code,
			Type: rgenerated.ASSETS,
		},
		Attributes: rgenerated.AssetAttributes{
			PreIssuanceAssetSigner: record.PreIssuanceAssetSigner,
			Details:                record.Details.AsRegourcesDetails(),
			MaxIssuanceAmount:      rgenerated.Amount(record.MaxIssuanceAmount),
			AvailableForIssuance:   rgenerated.Amount(record.AvailableForIssuance),
			Issued:                 rgenerated.Amount(record.Issued),
			PendingIssuance:        rgenerated.Amount(record.PendingIssuance),
			Policies:               xdr.AssetPolicy(record.Policies),
			TrailingDigits:         record.TrailingDigits,
			Type:                   record.Type,
		},
		Relationships: rgenerated.AssetRelationships{
			Owner: NewAccountKey(record.Owner).AsRelation(),
		},
	}
}

//NewAssetKey - creates new Key for asset
func NewAssetKey(assetCode string) rgenerated.Key {
	return rgenerated.Key{
		ID:   assetCode,
		Type: rgenerated.ASSETS,
	}
}
