package rgenerated

type AccountRole struct {
	Key
	Attributes    *AccountRoleAttributes    `json:"attributes,omitempty"`
	Relationships *AccountRoleRelationships `json:"relationships,omitempty"`
}
type AccountRoleResponse struct {
	Data     AccountRole `json:"data"`
	Included Included    `json:"included"`
}

type AccountRolesResponse struct {
	Data     []AccountRole `json:"data"`
	Included Included      `json:"included"`
	Links    *Links        `json:"links"`
}

// MustAccountRole - returns AccountRole from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustAccountRole(key Key) *AccountRole {
	var accountRole AccountRole
	if c.tryFindEntry(key, &accountRole) {
		return &accountRole
	}
	return nil
}
