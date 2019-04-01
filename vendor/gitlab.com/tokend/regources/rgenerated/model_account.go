package rgenerated

type Account struct {
	Key
	Relationships *AccountRelationships `json:"relationships,omitempty"`
}
type AccountResponse struct {
	Data     Account  `json:"data"`
	Included Included `json:"included"`
}

type AccountsResponse struct {
	Data     []Account `json:"data"`
	Included Included  `json:"included"`
	Links    *Links    `json:"links"`
}

// MustAccount - returns Account from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustAccount(key Key) *Account {
	var account Account
	if c.tryFindEntry(key, &account) {
		return &account
	}
	return nil
}
