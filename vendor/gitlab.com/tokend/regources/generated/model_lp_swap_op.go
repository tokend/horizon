/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package regources

import "encoding/json"

type LpSwapOp struct {
	Key
	Attributes    LpSwapOpAttributes    `json:"attributes"`
	Relationships LpSwapOpRelationships `json:"relationships"`
}
type LpSwapOpResponse struct {
	Data     LpSwapOp `json:"data"`
	Included Included `json:"included"`
}

type LpSwapOpListResponse struct {
	Data     []LpSwapOp      `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
	Meta     json.RawMessage `json:"meta,omitempty"`
}

func (r *LpSwapOpListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *LpSwapOpListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustLpSwapOp - returns LpSwapOp from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustLpSwapOp(key Key) *LpSwapOp {
	var lpSwapOp LpSwapOp
	if c.tryFindEntry(key, &lpSwapOp) {
		return &lpSwapOp
	}
	return nil
}
