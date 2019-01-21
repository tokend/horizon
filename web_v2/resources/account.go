package resources

import (
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/regources/v2"
)

//NewAccount - creates new instance of account
func NewAccount(core core2.Account) regources.Account {
	return regources.Account{
		Key: regources.Key{
			ID:   core.Address,
			Type: regources.TypeAccounts,
		},
	}
}

//NewAccountKey - creates account key from address
func NewAccountKey(address string) regources.Key {
	return regources.Key{
		ID:   address,
		Type: regources.TypeAccounts,
	}
}
