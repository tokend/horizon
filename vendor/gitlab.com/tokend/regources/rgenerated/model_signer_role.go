package rgenerated

type SignerRole struct {
	Key
	Attributes    *SignerRoleAttributes    `json:"attributes,omitempty"`
	Relationships *SignerRoleRelationships `json:"relationships,omitempty"`
}
type SignerRoleResponse struct {
	Data     SignerRole `json:"data"`
	Included Included   `json:"included"`
}

type SignerRolesResponse struct {
	Data     []SignerRole `json:"data"`
	Included Included     `json:"included"`
	Links    *Links       `json:"links"`
}

// MustSignerRole - returns SignerRole from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustSignerRole(key Key) *SignerRole {
	var signerRole SignerRole
	if c.tryFindEntry(key, &signerRole) {
		return &signerRole
	}
	return nil
}
