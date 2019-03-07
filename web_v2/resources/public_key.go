package resources

import (
	"gitlab.com/tokend/regources/v2"
)

// NewPublicKeyResource creates new instance of PublicKey resource from provided publicKey
func NewPublicKeyEntry(publicKey string) regources.PublicKeyEntry {
	return regources.PublicKeyEntry{
		Key: regources.Key{
			ID:   publicKey,
			Type: regources.TypePublicKeyEntries,
		},
	}
}

