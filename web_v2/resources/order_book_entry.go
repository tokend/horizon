package resources

import (
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/regources/rgenerated"
)

// NewOrderBookEntryKey - returns new instance of OrderBookEntryKey
func NewOrderBookEntryKey(id string) rgenerated.Key {
	return rgenerated.Key{
		ID:   id,
		Type: rgenerated.ORDER_BOOK_ENTRIES,
	}
}

// NewOrderBookEntry - returns new instance of OrderBookEntry
func NewOrderBookEntry(record core2.OrderBookEntry) rgenerated.OrderBookEntry {
	return rgenerated.OrderBookEntry{
		Key: NewOrderBookEntryKey(record.ID),
		Attributes: rgenerated.OrderBookEntryAttributes{
			IsBuy:       record.IsBuy,
			Price:       rgenerated.Amount(record.Price),
			BaseAmount:  rgenerated.Amount(record.BaseAmount),
			QuoteAmount: rgenerated.Amount(record.QuoteAmount),
			CreatedAt:   record.CreatedAt,
		},
		Relationships: rgenerated.OrderBookEntryRelationships{
			BaseAsset:  NewAssetKey(record.BaseAssetCode).AsRelation(),
			QuoteAsset: NewAssetKey(record.QuoteAssetCode).AsRelation(),
		},
	}
}
