/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package regources

import "encoding/json"

type Investment struct {
	Key
	Attributes InvestmentAttributes `json:"attributes"`
}
type InvestmentResponse struct {
	Data     Investment `json:"data"`
	Included Included   `json:"included"`
}

type InvestmentListResponse struct {
	Data     []Investment    `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
	Meta     json.RawMessage `json:"meta,omitempty"`
}

func (r *InvestmentListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *InvestmentListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustInvestment - returns Investment from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustInvestment(key Key) *Investment {
	var investment Investment
	if c.tryFindEntry(key, &investment) {
		return &investment
	}
	return nil
}
