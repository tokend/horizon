/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package regources

import "encoding/json"

type OrderBook struct {
	Key
	Relationships OrderBookRelationships `json:"relationships"`
}
type OrderBookResponse struct {
	Data     OrderBook `json:"data"`
	Included Included  `json:"included"`
}

type OrderBookListResponse struct {
	Data     []OrderBook     `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
	Meta     json.RawMessage `json:"meta,omitempty"`
}

func (r *OrderBookListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *OrderBookListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustOrderBook - returns OrderBook from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustOrderBook(key Key) *OrderBook {
	var orderBook OrderBook
	if c.tryFindEntry(key, &orderBook) {
		return &orderBook
	}
	return nil
}
