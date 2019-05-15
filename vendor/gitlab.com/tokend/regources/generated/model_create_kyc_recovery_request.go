/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package regources

type CreateKycRecoveryRequest struct {
	Key
	Attributes    CreateKycRecoveryRequestAttributes    `json:"attributes"`
	Relationships CreateKycRecoveryRequestRelationships `json:"relationships"`
}
type CreateKycRecoveryRequestResponse struct {
	Data     CreateKycRecoveryRequest `json:"data"`
	Included Included                 `json:"included"`
}

type CreateKycRecoveryRequestsResponse struct {
	Data     []CreateKycRecoveryRequest `json:"data"`
	Included Included                   `json:"included"`
	Links    *Links                     `json:"links"`
}

// MustCreateKycRecoveryRequest - returns CreateKycRecoveryRequest from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustCreateKycRecoveryRequest(key Key) *CreateKycRecoveryRequest {
	var createKYCRecoveryRequest CreateKycRecoveryRequest
	if c.tryFindEntry(key, &createKYCRecoveryRequest) {
		return &createKYCRecoveryRequest
	}
	return nil
}
