/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package regources

import "encoding/json"

type LiquidityPool struct {
	Key
	Attributes LiquidityPoolAttributes `json:"attributes"`
}
type LiquidityPoolResponse struct {
	Data     LiquidityPool `json:"data"`
	Included Included      `json:"included"`
}

type LiquidityPoolListResponse struct {
	Data     []LiquidityPool `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
	Meta     json.RawMessage `json:"meta,omitempty"`
}

func (r *LiquidityPoolListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *LiquidityPoolListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustLiquidityPool - returns LiquidityPool from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustLiquidityPool(key Key) *LiquidityPool {
	var liquidityPool LiquidityPool
	if c.tryFindEntry(key, &liquidityPool) {
		return &liquidityPool
	}
	return nil
}
