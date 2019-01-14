package resource

import (
	"gitlab.com/tokend/go/xdr"
	core "gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/horizon/resource/base"
)

// Account - resource object representing AccountEntry
type Account struct {
	ID           string       `jsonapi:"primary,accounts"`
	Role         *AccountRole `jsonapi:"attr,role,omitempty"`
	BlockReasons *Mask        `jsonapi:"attr,block_reasons,omitempty"`
	IsBlocked    bool         `jsonapi:"attr,is_blocked"`
	Balances     []*Balance   `jsonapi:"relation,balances,omitempty"`
}

//NewAccount - creates new account using core.Account
func NewAccount(record *core.Account) Account {
	return Account{
		ID: record.Address,
		Role: &AccountRole{
			// TODO: must use account role
			ID:   int64(record.AccountType),
			Name: xdr.AccountType(record.AccountType).ShortString(),
		},
		BlockReasons: &Mask{
			Mask:  record.BlockReasons,
			Flags: base.FlagFromXdrBlockReasons(record.BlockReasons, xdr.BlockReasonsAll),
		},
		IsBlocked: record.BlockReasons != 0,
	}
}

// AccountRole - represents account role which defines actions allowed to be performed by this account
type AccountRole struct {
	ID   int64  `jsonapi:"attr,id"`
	Name string `jsonapi:"attr,name"`
}
