/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package regources

import "encoding/json"

type CancelDataOwnerUpdateRequestOp struct {
	Key
	Relationships CancelDataOwnerUpdateRequestOpRelationships `json:"relationships"`
}
type CancelDataOwnerUpdateRequestOpResponse struct {
	Data     CancelDataOwnerUpdateRequestOp `json:"data"`
	Included Included                       `json:"included"`
}

type CancelDataOwnerUpdateRequestOpListResponse struct {
	Data     []CancelDataOwnerUpdateRequestOp `json:"data"`
	Included Included                         `json:"included"`
	Links    *Links                           `json:"links"`
	Meta     json.RawMessage                  `json:"meta,omitempty"`
}

func (r *CancelDataOwnerUpdateRequestOpListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *CancelDataOwnerUpdateRequestOpListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustCancelDataOwnerUpdateRequestOp - returns CancelDataOwnerUpdateRequestOp from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustCancelDataOwnerUpdateRequestOp(key Key) *CancelDataOwnerUpdateRequestOp {
	var cancelDataOwnerUpdateRequestOp CancelDataOwnerUpdateRequestOp
	if c.tryFindEntry(key, &cancelDataOwnerUpdateRequestOp) {
		return &cancelDataOwnerUpdateRequestOp
	}
	return nil
}
