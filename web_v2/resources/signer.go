package resources

import (
	"strconv"

	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/regources/rgenerated"
)

//NewSigner - creates new instance of signer
func NewSigner(signer core2.Signer) rgenerated.Signer {
	return rgenerated.Signer{
		Key: NewSignerKey(signer.PublicKey),
		Attributes: rgenerated.SignerAttributes{
			Weight:   signer.Weight,
			Identity: signer.Identity,
			Details:  signer.Details,
		},
		Relationships: rgenerated.SignerRelationships{
			Role: rgenerated.Key{
				ID:   strconv.FormatUint(signer.RoleID, 10),
				Type: rgenerated.SIGNER_ROLES,
			}.AsRelation(),
			Account: rgenerated.Key{
				ID:   signer.AccountID,
				Type: rgenerated.ACCOUNTS,
			}.AsRelation(),
		},
	}
}

//NewSignerKey - creates new key for signer
func NewSignerKey(publicKey string) rgenerated.Key {
	return rgenerated.Key{
		ID:   publicKey,
		Type: rgenerated.SIGNERS,
	}
}
