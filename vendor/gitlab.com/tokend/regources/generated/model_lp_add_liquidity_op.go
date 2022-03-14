/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package regources

import "encoding/json"

type LpAddLiquidityOp struct {
	Key
	Attributes    *LpAddLiquidityOpAttributes   `json:"attributes,omitempty"`
	Relationships LpAddLiquidityOpRelationships `json:"relationships"`
}
type LpAddLiquidityOpResponse struct {
	Data     LpAddLiquidityOp `json:"data"`
	Included Included         `json:"included"`
}

type LpAddLiquidityOpListResponse struct {
	Data     []LpAddLiquidityOp `json:"data"`
	Included Included           `json:"included"`
	Links    *Links             `json:"links"`
	Meta     json.RawMessage    `json:"meta,omitempty"`
}

func (r *LpAddLiquidityOpListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *LpAddLiquidityOpListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustLpAddLiquidityOp - returns LpAddLiquidityOp from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustLpAddLiquidityOp(key Key) *LpAddLiquidityOp {
	var lpAddLiquidityOp LpAddLiquidityOp
	if c.tryFindEntry(key, &lpAddLiquidityOp) {
		return &lpAddLiquidityOp
	}
	return nil
}
