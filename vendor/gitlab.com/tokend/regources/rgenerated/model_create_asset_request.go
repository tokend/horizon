package rgenerated

type CreateAssetRequest struct {
	Key
	Attributes    *CreateAssetRequestAttributes    `json:"attributes,omitempty"`
	Relationships *CreateAssetRequestRelationships `json:"relationships,omitempty"`
}
type CreateAssetRequestResponse struct {
	Data     CreateAssetRequest `json:"data"`
	Included Included           `json:"included"`
}

type CreateAssetRequestsResponse struct {
	Data     []CreateAssetRequest `json:"data"`
	Included Included             `json:"included"`
	Links    *Links               `json:"links"`
}

// MustCreateAssetRequest - returns CreateAssetRequest from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustCreateAssetRequest(key Key) *CreateAssetRequest {
	var createAssetRequest CreateAssetRequest
	if c.tryFindEntry(key, &createAssetRequest) {
		return &createAssetRequest
	}
	return nil
}
