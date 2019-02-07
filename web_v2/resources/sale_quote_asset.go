package resources

import (
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/v2"
)

// NewSaleQuoteAssetKey - returns new instance of SaleQuoteAssetKey
func NewSaleQuoteAssetKey(assetCode string) regources.Key {
	return regources.Key{
		ID:   assetCode,
		Type: regources.TypeSaleQuoteAssets,
	}
}

// NewSaleQuoteAsset - returns new instance of SaleQuoteAsset
func NewSaleQuoteAsset(qa history2.SaleQuoteAsset) regources.SaleQuoteAsset {
	return regources.SaleQuoteAsset{
		Key: NewSaleQuoteAssetKey(qa.Asset),
		Attributes: regources.SaleQuoteAssetAttrs{
			Price:           qa.Price,
			CurrentCap:      qa.CurrentCap,
			TotalCurrentCap: qa.TotalCurrentCap,
			HardCap:         qa.HardCap,
		},
		Relationships: regources.SaleQuoteAssetRelations{
			Asset: NewAssetKey(qa.Asset).AsRelation(),
		},
	}
}
