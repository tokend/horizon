package regources

type AccountResponse struct {
	Data     Account  `json:"data"`
	Included Included `json:"included"`
}

// Account - Resource object representing AccountEntry
type Account struct {
	Key
	Relationships AccountRelationships `json:"relationships"`
}

type AccountRelationships struct {
	Role     *Relation           `json:"role,omitempty"`
	Balances *RelationCollection `json:"balances,omitempty"`
	Referrer *Relation           `json:"referrer,omitempty"`
	Limits   *Relation           `json:"limits,omitempty"`
}
