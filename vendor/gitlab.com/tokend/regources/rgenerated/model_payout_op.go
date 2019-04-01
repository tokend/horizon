package rgenerated

type PayoutOp struct {
	Key
	Attributes    *PayoutOpAttributes    `json:"attributes,omitempty"`
	Relationships *PayoutOpRelationships `json:"relationships,omitempty"`
}
type PayoutOpResponse struct {
	Data     PayoutOp `json:"data"`
	Included Included `json:"included"`
}

type PayoutOpsResponse struct {
	Data     []PayoutOp `json:"data"`
	Included Included   `json:"included"`
	Links    *Links     `json:"links"`
}

// MustPayoutOp - returns PayoutOp from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustPayoutOp(key Key) *PayoutOp {
	var payoutOp PayoutOp
	if c.tryFindEntry(key, &payoutOp) {
		return &payoutOp
	}
	return nil
}
