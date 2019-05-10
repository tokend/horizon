package resources

import (
	"fmt"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/generated"
)

// NewMatchKey - creates new Key for OrderBook
func NewMatchKey(id string) regources.Key {
	return regources.Key{
		ID:   id,
		Type: regources.ORDER_BOOKS,
	}
}

// NewMatch - returns new instance of `Match` resource
func NewMatch(record history2.Match) regources.Match {
	return regources.Match{
		Key: NewMatchKey(record.ID),
		Attributes: regources.MatchAttributes{
			BaseAmount:  regources.Amount(record.BaseAmount),
			OrderBookId: fmt.Sprintf("%d", record.OrderBookID),
			Price:       regources.Amount(record.Price),
			QuoteAmount: regources.Amount(record.QuoteAmount),
		},
		Relationships: regources.MatchRelationships{
			BaseAsset:  *NewAssetKey(record.BaseAsset).AsRelation(),
			QuoteAsset: *NewAssetKey(record.QuoteAsset).AsRelation(),
		},
	}
}
