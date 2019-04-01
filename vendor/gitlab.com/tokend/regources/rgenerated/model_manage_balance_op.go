package rgenerated

type ManageBalanceOp struct {
	Key
	Attributes    *ManageBalanceOpAttributes    `json:"attributes,omitempty"`
	Relationships *ManageBalanceOpRelationships `json:"relationships,omitempty"`
}
type ManageBalanceOpResponse struct {
	Data     ManageBalanceOp `json:"data"`
	Included Included        `json:"included"`
}

type ManageBalanceOpsResponse struct {
	Data     []ManageBalanceOp `json:"data"`
	Included Included          `json:"included"`
	Links    *Links            `json:"links"`
}

// MustManageBalanceOp - returns ManageBalanceOp from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustManageBalanceOp(key Key) *ManageBalanceOp {
	var manageBalanceOp ManageBalanceOp
	if c.tryFindEntry(key, &manageBalanceOp) {
		return &manageBalanceOp
	}
	return nil
}
