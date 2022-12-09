/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package regources

import "encoding/json"

type LpManageLiquidityOp struct {
	Key
	Attributes    LpManageLiquidityOpAttributes    `json:"attributes"`
	Relationships LpManageLiquidityOpRelationships `json:"relationships"`
}
type LpManageLiquidityOpResponse struct {
	Data     LpManageLiquidityOp `json:"data"`
	Included Included            `json:"included"`
}

type LpManageLiquidityOpListResponse struct {
	Data     []LpManageLiquidityOp `json:"data"`
	Included Included              `json:"included"`
	Links    *Links                `json:"links"`
	Meta     json.RawMessage       `json:"meta,omitempty"`
}

func (r *LpManageLiquidityOpListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *LpManageLiquidityOpListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustLpManageLiquidityOp - returns LpManageLiquidityOp from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustLpManageLiquidityOp(key Key) *LpManageLiquidityOp {
	var lpManageLiquidityOp LpManageLiquidityOp
	if c.tryFindEntry(key, &lpManageLiquidityOp) {
		return &lpManageLiquidityOp
	}
	return nil
}
