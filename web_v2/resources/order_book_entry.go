package resources

import (
	"gitlab.com/tokend/horizon/db2/core2"
	regources "gitlab.com/tokend/regources/v2/generated"
)

// NewOrderBookEntryKey - returns new instance of OrderBookEntryKey
func NewOrderBookEntryKey(id string) regources.Key {
	return regources.Key{
		ID:   id,
		Type: regources.ORDER_BOOK_ENTRIES,
	}
}

// NewOrderBookEntry - returns new instance of OrderBookEntry
func NewOrderBookEntry(record core2.OrderBookEntry) regources.OrderBookEntry {
	return regources.OrderBookEntry{
		Key: NewOrderBookEntryKey(record.ID),
		Attributes: regources.OrderBookEntryAttributes{
			IsBuy:       record.IsBuy,
			Price:       regources.Amount(record.Price),
			BaseAmount:  regources.Amount(record.BaseAmount),
			QuoteAmount: regources.Amount(record.QuoteAmount),
			CreatedAt:   record.CreatedAt,
		},
		Relationships: regources.OrderBookEntryRelationships{
			BaseAsset:  NewAssetKey(record.BaseAssetCode).AsRelation(),
			QuoteAsset: NewAssetKey(record.QuoteAssetCode).AsRelation(),
		},
	}
}
