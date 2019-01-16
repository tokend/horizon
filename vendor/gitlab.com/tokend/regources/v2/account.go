package regources

// Account - resource object representing AccountEntry
type Account struct {
	ID       string        `jsonapi:"primary,accounts"`
	State    *AccountState `jsonapi:"relation,state,omitempty"`
	Role     *Role         `jsonapi:"relation,role,omitempty"`
	Balances []*Balance    `jsonapi:"relation,balances,omitempty"`
	Referrer *Account      `jsonapi:"relation,referrer,omitempty"`
}

type AccountState struct {
	ID           string `jsonapi:"primary,account_states"`
	BlockReasons *Mask  `jsonapi:"attr,block_reasons,omitempty"`
	IsBlocked    bool   `jsonapi:"attr,is_blocked"`
}
