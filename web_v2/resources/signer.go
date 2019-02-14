package resources

import (
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/regources/v2"
)

//NewSigner - creates new instance of signer
func NewSigner(signer core2.Signer) regources.Signer {
	return regources.Signer{
		Key: regources.Key{
			ID:   signer.ID,
			Type: regources.TypeSigners,
		},
		Attributes: regources.SignerAttrs{
			Weight:   signer.Weight,
			Identity: signer.Identity,
			// TODO: FIX ME after roles
			Details: map[string]interface{}{
				"name": signer.Details,
			},
		},
	}
}
