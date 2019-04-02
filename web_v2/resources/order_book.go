package resources

import (
	"gitlab.com/tokend/regources/generated"
)

// TODO: tmp here, remove and generate from docs
type OrderBookResponse struct {
	Data     OrderBook          `json:"data"`
	Included regources.Included `json:"included"`
}

type OrderBook struct {
	regources.Key
	Relationships OrderBookRelations `json:"relationships"`
}

type OrderBookRelations struct {
	BuyEntries  *regources.RelationCollection `json:"buy_entries"`
	SellEntries *regources.RelationCollection `json:"sell_entries"`
}

// NewOrderBookKey - creates new Key for OrderBook
func NewOrderBookKey(id int64) regources.Key {
	return regources.NewKeyInt64(id, "order-books")
}

// NewOrderBook - creates new instance of OrderBook
func NewOrderBook(id int64) OrderBook {
	return OrderBook{
		Key: NewOrderBookKey(id),
	}
}
