/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type CreateAccountOp struct {
	Key
	Relationships CreateAccountOpRelationships `json:"relationships"`
}
type CreateAccountOpResponse struct {
	Data     CreateAccountOp `json:"data"`
	Included Included        `json:"included"`
}

type CreateAccountOpsResponse struct {
	Data     []CreateAccountOp `json:"data"`
	Included Included          `json:"included"`
	Links    *Links            `json:"links"`
}

// MustCreateAccountOp - returns CreateAccountOp from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustCreateAccountOp(key Key) *CreateAccountOp {
	var createAccountOp CreateAccountOp
	if c.tryFindEntry(key, &createAccountOp) {
		return &createAccountOp
	}
	return nil
}
