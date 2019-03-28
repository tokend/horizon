package resources

import (
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/rgenerated"
)

// NewSaleQuoteAssetKey - returns new instance of SaleQuoteAssetKey
func NewSaleQuoteAssetKey(assetCode string) rgenerated.Key {
	return rgenerated.Key{
		ID:   assetCode,
		Type: rgenerated.SALE_QUOTE_ASSETS,
	}
}

// NewSaleQuoteAsset - returns new instance of SaleQuoteAsset
func NewSaleQuoteAsset(qa history2.SaleQuoteAsset) rgenerated.SaleQuoteAsset {
	return rgenerated.SaleQuoteAsset{
		Key: NewSaleQuoteAssetKey(qa.Asset),
		Attributes: rgenerated.SaleQuoteAssetAttributes{
			Price:           qa.Price,
			CurrentCap:      qa.CurrentCap,
			TotalCurrentCap: qa.TotalCurrentCap,
			HardCap:         qa.HardCap,
		},
		Relationships: rgenerated.SaleQuoteAssetRelationships{
			Asset: NewAssetKey(qa.Asset).AsRelation(),
		},
	}
}
