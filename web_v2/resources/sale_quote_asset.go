package resources

import (
	"gitlab.com/tokend/horizon/db2/history2"
	regources "gitlab.com/tokend/regources/v2/generated"
)

// NewSaleQuoteAssetKey - returns new instance of SaleQuoteAssetKey
func NewSaleQuoteAssetKey(assetCode string) regources.Key {
	return regources.Key{
		ID:   assetCode,
		Type: regources.SALE_QUOTE_ASSETS,
	}
}

// NewSaleQuoteAsset - returns new instance of SaleQuoteAsset
func NewSaleQuoteAsset(qa history2.SaleQuoteAsset) regources.SaleQuoteAsset {
	return regources.SaleQuoteAsset{
		Key: NewSaleQuoteAssetKey(qa.Asset),
		Attributes: regources.SaleQuoteAssetAttributes{
			Price:           qa.Price,
			CurrentCap:      qa.CurrentCap,
			TotalCurrentCap: qa.TotalCurrentCap,
			HardCap:         qa.HardCap,
		},
		Relationships: regources.SaleQuoteAssetRelationships{
			Asset: NewAssetKey(qa.Asset).AsRelation(),
		},
	}
}
