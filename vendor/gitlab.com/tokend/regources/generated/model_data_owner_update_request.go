/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package regources

import "encoding/json"

type DataOwnerUpdateRequest struct {
	Key
	Attributes    DataOwnerUpdateRequestAttributes    `json:"attributes"`
	Relationships DataOwnerUpdateRequestRelationships `json:"relationships"`
}
type DataOwnerUpdateRequestResponse struct {
	Data     DataOwnerUpdateRequest `json:"data"`
	Included Included               `json:"included"`
}

type DataOwnerUpdateRequestListResponse struct {
	Data     []DataOwnerUpdateRequest `json:"data"`
	Included Included                 `json:"included"`
	Links    *Links                   `json:"links"`
	Meta     json.RawMessage          `json:"meta,omitempty"`
}

func (r *DataOwnerUpdateRequestListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *DataOwnerUpdateRequestListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustDataOwnerUpdateRequest - returns DataOwnerUpdateRequest from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustDataOwnerUpdateRequest(key Key) *DataOwnerUpdateRequest {
	var dataOwnerUpdateRequest DataOwnerUpdateRequest
	if c.tryFindEntry(key, &dataOwnerUpdateRequest) {
		return &dataOwnerUpdateRequest
	}
	return nil
}
