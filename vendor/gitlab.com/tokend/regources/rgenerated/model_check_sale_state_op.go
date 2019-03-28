/*
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package rgenerated

type CheckSaleStateOp struct {
	Key
	Attributes    CheckSaleStateOpAttributes    `json:"attributes"`
	Relationships CheckSaleStateOpRelationships `json:"relationships"`
}
type CheckSaleStateOpResponse struct {
	Data     CheckSaleStateOp `json:"data"`
	Included Included         `json:"included"`
}

type CheckSaleStateOpsResponse struct {
	Data     []CheckSaleStateOp `json:"data"`
	Included Included           `json:"included"`
	Links    *Links             `json:"links"`
}

// MustCheckSaleStateOp - returns CheckSaleStateOp from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustCheckSaleStateOp(key Key) *CheckSaleStateOp {
	var checkSaleStateOp CheckSaleStateOp
	if c.tryFindEntry(key, &checkSaleStateOp) {
		return &checkSaleStateOp
	}
	return nil
}
