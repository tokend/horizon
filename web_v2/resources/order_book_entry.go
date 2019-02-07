package resources

import (
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/regources/v2"
)

// NewOrderBookEntryKey - returns new instance of OrderBookEntryKey
func NewOrderBookEntryKey(id string) regources.Key {
	return regources.Key{
		ID:   id,
		Type: regources.TypeOrderBookEntries,
	}
}

// NewOrderBookEntry - returns new instance of OrderBookEntry
func NewOrderBookEntry(record core2.OrderBookEntry) regources.OrderBookEntry {
	return regources.OrderBookEntry{
		Key: NewOrderBookEntryKey(record.ID),
		Attributes: regources.OrderBookEntryAttrs{
			IsBuy:       record.IsBuy,
			Price:       regources.Amount(record.Price),
			BaseAmount:  regources.Amount(record.BaseAmount),
			QuoteAmount: regources.Amount(record.QuoteAmount),
			CreatedAt:   record.CreatedAt,
		},
		Relationships: regources.OrderBookEntryRelations{
			BaseAsset:  NewAssetKey(record.BaseAssetCode).AsRelation(),
			QuoteAsset: NewAssetKey(record.QuoteAssetCode).AsRelation(),
		},
	}
}
