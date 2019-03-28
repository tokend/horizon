package resources

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/rgenerated"
)

//NewSaleKey - creates new Key for asset
func NewSaleKey(id int64) rgenerated.Key {
	return rgenerated.NewKeyInt64(id, rgenerated.SALES)
}

// NewSale - creates new instance of Sale
func NewSale(record history2.Sale) rgenerated.Sale {
	quoteAssets := &rgenerated.RelationCollection{
		Data: make([]rgenerated.Key, 0, len(record.QuoteAssets.QuoteAssets)),
	}

	for _, quoteAsset := range record.QuoteAssets.QuoteAssets {
		quoteAssets.Data = append(quoteAssets.Data, newQuoteAssetKey(quoteAsset.Asset))
	}

	return rgenerated.Sale{
		Key: NewSaleKey(int64(record.ID)),
		Attributes: rgenerated.SaleAttributes{
			StartTime: record.StartTime,
			EndTime:   record.EndTime,
			SaleType:  record.SaleType,
			SaleState: record.State,
			Details:   record.Details,
		},
		Relationships: rgenerated.SaleRelationships{
			Owner:             NewAccountKey(record.OwnerAddress).AsRelation(),
			BaseAsset:         NewAssetKey(record.BaseAsset).AsRelation(),
			QuoteAssets:       quoteAssets,
			DefaultQuoteAsset: newQuoteAssetKey(record.DefaultQuoteAsset).AsRelation(),
		},
	}
}

// NewSaleDefaultQuoteAsset - extracts the default quote asset details from the sale
func NewSaleDefaultQuoteAsset(saleRecord history2.Sale) rgenerated.SaleQuoteAsset {
	return rgenerated.SaleQuoteAsset{
		Key: *NewSaleQuoteAssetKey(saleRecord.DefaultQuoteAsset).GetKeyP(),
		Attributes: rgenerated.SaleQuoteAssetAttributes{
			Price:      amount.One,
			CurrentCap: saleRecord.CurrentCap,
			HardCap:    saleRecord.HardCap,
			SoftCap:    saleRecord.SoftCap,
		},
		Relationships: rgenerated.SaleQuoteAssetRelationships{
			Asset: NewAssetKey(saleRecord.DefaultQuoteAsset).AsRelation(),
		},
	}
}
