package rgenerated

type Limits struct {
	Key
	Attributes    *LimitsAttributes    `json:"attributes,omitempty"`
	Relationships *LimitsRelationships `json:"relationships,omitempty"`
}
type LimitsResponse struct {
	Data     Limits   `json:"data"`
	Included Included `json:"included"`
}

type LimitssResponse struct {
	Data     []Limits `json:"data"`
	Included Included `json:"included"`
	Links    *Links   `json:"links"`
}

// MustLimits - returns Limits from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustLimits(key Key) *Limits {
	var limits Limits
	if c.tryFindEntry(key, &limits) {
		return &limits
	}
	return nil
}
