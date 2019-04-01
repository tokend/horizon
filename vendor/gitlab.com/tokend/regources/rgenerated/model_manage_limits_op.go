package rgenerated

type ManageLimitsOp struct {
	Key
	Attributes *ManageLimitsOpAttributes `json:"attributes,omitempty"`
}
type ManageLimitsOpResponse struct {
	Data     ManageLimitsOp `json:"data"`
	Included Included       `json:"included"`
}

type ManageLimitsOpsResponse struct {
	Data     []ManageLimitsOp `json:"data"`
	Included Included         `json:"included"`
	Links    *Links           `json:"links"`
}

// MustManageLimitsOp - returns ManageLimitsOp from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustManageLimitsOp(key Key) *ManageLimitsOp {
	var manageLimitsOp ManageLimitsOp
	if c.tryFindEntry(key, &manageLimitsOp) {
		return &manageLimitsOp
	}
	return nil
}
