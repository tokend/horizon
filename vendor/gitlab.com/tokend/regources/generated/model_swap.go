/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package regources

type Swap struct {
	Key
	Attributes    SwapAttributes    `json:"attributes"`
	Relationships SwapRelationships `json:"relationships"`
}
type SwapResponse struct {
	Data     Swap     `json:"data"`
	Included Included `json:"included"`
}

type SwapListResponse struct {
	Data     []Swap   `json:"data"`
	Included Included `json:"included"`
	Links    *Links   `json:"links"`
}

// MustSwap - returns Swap from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustSwap(key Key) *Swap {
	var swap Swap
	if c.tryFindEntry(key, &swap) {
		return &swap
	}
	return nil
}
