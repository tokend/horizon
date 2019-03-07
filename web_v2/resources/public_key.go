package resources

import (
	"gitlab.com/tokend/regources/v2"
)

// NewPublicKeyResource creates new instance of PublicKey resource from provided publicKey
func NewPublicKeyResource(publicKey string) regources.PublicKey {
	return regources.PublicKey{
		Key: regources.Key{
			ID:   publicKey,
			Type: regources.TypePublicKey,
		},
	}
}

