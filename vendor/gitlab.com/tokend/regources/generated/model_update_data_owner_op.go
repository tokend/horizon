/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package regources

import "encoding/json"

type UpdateDataOwnerOp struct {
	Key
	Relationships UpdateDataOwnerOpRelationships `json:"relationships"`
}
type UpdateDataOwnerOpResponse struct {
	Data     UpdateDataOwnerOp `json:"data"`
	Included Included          `json:"included"`
}

type UpdateDataOwnerOpListResponse struct {
	Data     []UpdateDataOwnerOp `json:"data"`
	Included Included            `json:"included"`
	Links    *Links              `json:"links"`
	Meta     json.RawMessage     `json:"meta,omitempty"`
}

func (r *UpdateDataOwnerOpListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *UpdateDataOwnerOpListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustUpdateDataOwnerOp - returns UpdateDataOwnerOp from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustUpdateDataOwnerOp(key Key) *UpdateDataOwnerOp {
	var updateDataOwnerOp UpdateDataOwnerOp
	if c.tryFindEntry(key, &updateDataOwnerOp) {
		return &updateDataOwnerOp
	}
	return nil
}
