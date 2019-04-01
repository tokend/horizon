package rgenerated

type CreateAtomicSwapBidRequestOp struct {
	Key
	Attributes    *CreateAtomicSwapBidRequestOpAttributes    `json:"attributes,omitempty"`
	Relationships *CreateAtomicSwapBidRequestOpRelationships `json:"relationships,omitempty"`
}
type CreateAtomicSwapBidRequestOpResponse struct {
	Data     CreateAtomicSwapBidRequestOp `json:"data"`
	Included Included                     `json:"included"`
}

type CreateAtomicSwapBidRequestOpsResponse struct {
	Data     []CreateAtomicSwapBidRequestOp `json:"data"`
	Included Included                       `json:"included"`
	Links    *Links                         `json:"links"`
}

// MustCreateAtomicSwapBidRequestOp - returns CreateAtomicSwapBidRequestOp from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustCreateAtomicSwapBidRequestOp(key Key) *CreateAtomicSwapBidRequestOp {
	var createAtomicSwapBidRequestOp CreateAtomicSwapBidRequestOp
	if c.tryFindEntry(key, &createAtomicSwapBidRequestOp) {
		return &createAtomicSwapBidRequestOp
	}
	return nil
}
