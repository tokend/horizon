package regources

import "time"

// OrderBookEntriesResponse - represents the order book response
type OrderBookEntriesResponse struct {
	Links    *Links           `json:"links"`
	Data     []OrderBookEntry `json:"data"`
	Included Included         `json:"included"`
}

// OrderBookEntry - represents the order book entry
type OrderBookEntry struct {
	Key
	Attributes    OrderBookEntryAttrs     `json:"attributes"`
	Relationships OrderBookEntryRelations `json:"relationships"`
}

// OrderBookEntryAttrs - represents the order book entry attributes
type OrderBookEntryAttrs struct {
	IsBuy       bool      `json:"is_buy"`
	Price       string    `json:"price"`
	BaseAmount  string    `json:"base_amount"`
	QuoteAmount string    `json:"quote_amount"`
	CreatedAt   time.Time `json:"created_at"`
}

// OrderBookEntryRelations - represents the order book entry relationships
type OrderBookEntryRelations struct {
	Offer      *Relation `json:"offer"`
	BaseAsset  *Relation `json:"base_asset"`
	QuoteAsset *Relation `json:"quote_asset"`
}
