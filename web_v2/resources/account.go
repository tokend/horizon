package resources

import (
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/regources/rgenerated"
)

//NewAccount - creates new instance of account
func NewAccount(core core2.Account) rgenerated.Account {
	return rgenerated.Account{
		Key: rgenerated.Key{
			ID:   core.Address,
			Type: rgenerated.ACCOUNTS,
		},
	}
}

//NewAccountKey - creates account key from address
func NewAccountKey(address string) rgenerated.Key {
	return rgenerated.Key{
		ID:   address,
		Type: rgenerated.ACCOUNTS,
	}
}
