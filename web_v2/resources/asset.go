package resources

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/core2"
	regources "gitlab.com/tokend/regources/v2/generated"
)

//NewAsset - creates new instance of Asset from provided one.
func NewAsset(record core2.Asset) regources.Asset {
	return regources.Asset{
		Key: regources.Key{
			ID:   record.Code,
			Type: regources.ASSETS,
		},
		Attributes: regources.AssetAttributes{
			PreIssuanceAssetSigner: record.PreIssuanceAssetSigner,
			Details:                record.Details.AsRegourcesDetails(),
			MaxIssuanceAmount:      regources.Amount(record.MaxIssuanceAmount),
			AvailableForIssuance:   regources.Amount(record.AvailableForIssuance),
			Issued:                 regources.Amount(record.Issued),
			PendingIssuance:        regources.Amount(record.PendingIssuance),
			Policies:               xdr.AssetPolicy(record.Policies),
			TrailingDigits:         record.TrailingDigits,
			Type:                   record.Type,
		},
		Relationships: regources.AssetRelationships{
			Owner: NewAccountKey(record.Owner).AsRelation(),
		},
	}
}

//NewAssetKey - creates new Key for asset
func NewAssetKey(assetCode string) regources.Key {
	return regources.Key{
		ID:   assetCode,
		Type: regources.ASSETS,
	}
}
