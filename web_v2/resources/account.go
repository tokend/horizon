package resources

import (
	"gitlab.com/tokend/horizon/db2/core2"
	regources "gitlab.com/tokend/regources/generated"
)

//NewAccount - creates new instance of account
func NewAccount(core core2.Account, accountStatus *regources.KYCRecoveryStatus, accountSigners ...core2.Signer) regources.Account {

	account := regources.Account{
		Key: regources.Key{
			ID:   core.Address,
			Type: regources.ACCOUNTS,
		},
		Relationships: regources.AccountRelationships{
			Signers: &regources.RelationCollection{
				Data: make([]regources.Key, 0, len(accountSigners)),
			},
		},
	}
	if accountStatus != nil {
		account.Attributes = regources.AccountAttributes{
			KycRecoveryStatus: accountStatus,
		}
	}

	for _, s := range accountSigners {
		account.Relationships.Signers.Data = append(account.Relationships.Signers.Data, NewSignerKey(s.PublicKey))
	}

	return account
}

//NewAccountKey - creates account key from address
func NewAccountKey(address string) regources.Key {
	return regources.Key{
		ID:   address,
		Type: regources.ACCOUNTS,
	}
}
