/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package regources

type CloseSwapOp struct {
	Key
	Attributes    CloseSwapOpAttributes    `json:"attributes"`
	Relationships CloseSwapOpRelationships `json:"relationships"`
}
type CloseSwapOpResponse struct {
	Data     CloseSwapOp `json:"data"`
	Included Included    `json:"included"`
}

type CloseSwapOpListResponse struct {
	Data     []CloseSwapOp `json:"data"`
	Included Included      `json:"included"`
	Links    *Links        `json:"links"`
}

// MustCloseSwapOp - returns CloseSwapOp from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustCloseSwapOp(key Key) *CloseSwapOp {
	var closeSwapOp CloseSwapOp
	if c.tryFindEntry(key, &closeSwapOp) {
		return &closeSwapOp
	}
	return nil
}
