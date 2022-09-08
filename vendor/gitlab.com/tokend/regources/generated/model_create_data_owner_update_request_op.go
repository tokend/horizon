/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package regources

import "encoding/json"

type CreateDataOwnerUpdateRequestOp struct {
	Key
	Attributes    CreateDataOwnerUpdateRequestOpAttributes    `json:"attributes"`
	Relationships CreateDataOwnerUpdateRequestOpRelationships `json:"relationships"`
}
type CreateDataOwnerUpdateRequestOpResponse struct {
	Data     CreateDataOwnerUpdateRequestOp `json:"data"`
	Included Included                       `json:"included"`
}

type CreateDataOwnerUpdateRequestOpListResponse struct {
	Data     []CreateDataOwnerUpdateRequestOp `json:"data"`
	Included Included                         `json:"included"`
	Links    *Links                           `json:"links"`
	Meta     json.RawMessage                  `json:"meta,omitempty"`
}

func (r *CreateDataOwnerUpdateRequestOpListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *CreateDataOwnerUpdateRequestOpListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustCreateDataOwnerUpdateRequestOp - returns CreateDataOwnerUpdateRequestOp from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustCreateDataOwnerUpdateRequestOp(key Key) *CreateDataOwnerUpdateRequestOp {
	var createDataOwnerUpdateRequestOp CreateDataOwnerUpdateRequestOp
	if c.tryFindEntry(key, &createDataOwnerUpdateRequestOp) {
		return &createDataOwnerUpdateRequestOp
	}
	return nil
}
