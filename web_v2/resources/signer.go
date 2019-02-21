package resources

import (
	"strconv"

	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/regources/v2"
)

//NewSigner - creates new instance of signer
func NewSigner(signer core2.Signer) regources.Signer {
	return regources.Signer{
		Key: NewSignerKey(signer.PublicKey),
		Attributes: regources.SignerAttrs{
			Weight:   signer.Weight,
			Identity: signer.Identity,
			Details:  signer.Details,
		},
		Relationships: regources.SignerRelation{
			Role: regources.Key{
				ID:   strconv.FormatUint(signer.RoleID, 10),
				Type: regources.TypeSignerRoles,
			}.AsRelation(),
			Account: regources.Key{
				ID:   signer.AccountID,
				Type: regources.TypeAccounts,
			}.AsRelation(),
		},
	}
}

//NewSignerKey - creates new key for signer
func NewSignerKey(publicKey string) regources.Key {
	return regources.Key{
		ID:   publicKey,
		Type: regources.TypeSigners,
	}
}
