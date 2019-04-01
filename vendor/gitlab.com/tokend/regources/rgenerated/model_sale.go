package rgenerated

type Sale struct {
	Key
	Attributes    *SaleAttributes    `json:"attributes,omitempty"`
	Relationships *SaleRelationships `json:"relationships,omitempty"`
}
type SaleResponse struct {
	Data     Sale     `json:"data"`
	Included Included `json:"included"`
}

type SalesResponse struct {
	Data     []Sale   `json:"data"`
	Included Included `json:"included"`
	Links    *Links   `json:"links"`
}

// MustSale - returns Sale from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustSale(key Key) *Sale {
	var sale Sale
	if c.tryFindEntry(key, &sale) {
		return &sale
	}
	return nil
}
