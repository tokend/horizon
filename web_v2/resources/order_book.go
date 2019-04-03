package resources

import (
	"gitlab.com/tokend/regources/generated"
)

// NewOrderBookKey - creates new Key for OrderBook
func NewOrderBookKey(id int64) regources.Key {
	return regources.NewKeyInt64(id, "order-books")
}

// NewOrderBook - creates new instance of OrderBook
func NewOrderBook(id int64) regources.OrderBook {
	return regources.OrderBook{
		Key: NewOrderBookKey(id),
	}
}
