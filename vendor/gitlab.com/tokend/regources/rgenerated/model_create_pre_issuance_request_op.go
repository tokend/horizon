package rgenerated

type CreatePreIssuanceRequestOp struct {
	Key
	Attributes    *CreatePreIssuanceRequestOpAttributes    `json:"attributes,omitempty"`
	Relationships *CreatePreIssuanceRequestOpRelationships `json:"relationships,omitempty"`
}
type CreatePreIssuanceRequestOpResponse struct {
	Data     CreatePreIssuanceRequestOp `json:"data"`
	Included Included                   `json:"included"`
}

type CreatePreIssuanceRequestOpsResponse struct {
	Data     []CreatePreIssuanceRequestOp `json:"data"`
	Included Included                     `json:"included"`
	Links    *Links                       `json:"links"`
}

// MustCreatePreIssuanceRequestOp - returns CreatePreIssuanceRequestOp from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustCreatePreIssuanceRequestOp(key Key) *CreatePreIssuanceRequestOp {
	var createPreIssuanceRequestOp CreatePreIssuanceRequestOp
	if c.tryFindEntry(key, &createPreIssuanceRequestOp) {
		return &createPreIssuanceRequestOp
	}
	return nil
}
