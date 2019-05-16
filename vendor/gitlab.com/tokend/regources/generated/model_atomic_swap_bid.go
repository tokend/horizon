/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package regources

type AtomicSwapBid struct {
	Key
	Attributes    AtomicSwapBidAttributes    `json:"attributes"`
	Relationships AtomicSwapBidRelationships `json:"relationships"`
}
type AtomicSwapBidResponse struct {
	Data     AtomicSwapBid `json:"data"`
	Included Included      `json:"included"`
}

type AtomicSwapBidListResponse struct {
	Data     []AtomicSwapBid `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
}

// MustAtomicSwapBid - returns AtomicSwapBid from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustAtomicSwapBid(key Key) *AtomicSwapBid {
	var atomicSwapBid AtomicSwapBid
	if c.tryFindEntry(key, &atomicSwapBid) {
		return &atomicSwapBid
	}
	return nil
}
