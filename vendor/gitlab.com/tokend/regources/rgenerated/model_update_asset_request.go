package rgenerated

type UpdateAssetRequest struct {
	Key
	Attributes    *UpdateAssetRequestAttributes    `json:"attributes,omitempty"`
	Relationships *UpdateAssetRequestRelationships `json:"relationships,omitempty"`
}
type UpdateAssetRequestResponse struct {
	Data     UpdateAssetRequest `json:"data"`
	Included Included           `json:"included"`
}

type UpdateAssetRequestsResponse struct {
	Data     []UpdateAssetRequest `json:"data"`
	Included Included             `json:"included"`
	Links    *Links               `json:"links"`
}

// MustUpdateAssetRequest - returns UpdateAssetRequest from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustUpdateAssetRequest(key Key) *UpdateAssetRequest {
	var updateAssetRequest UpdateAssetRequest
	if c.tryFindEntry(key, &updateAssetRequest) {
		return &updateAssetRequest
	}
	return nil
}
