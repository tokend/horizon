package resources

import (
	"gitlab.com/tokend/horizon/db2/core2"
	regources "gitlab.com/tokend/regources/generated"
)

//NewAccount - creates new instance of account
func NewAccount(core core2.Account, accountStatus regources.KYCRecoveryStatus) regources.Account {
	return regources.Account{
		Key: regources.Key{
			ID:   core.Address,
			Type: regources.ACCOUNTS,
		},
		Attributes: regources.AccountAttributes{
			KycRecoveryStatus: accountStatus.String(),
		},
	}
}

//NewAccountKey - creates account key from address
func NewAccountKey(address string) regources.Key {
	return regources.Key{
		ID:   address,
		Type: regources.ACCOUNTS,
	}
}
