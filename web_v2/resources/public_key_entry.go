package resources

import (
	"gitlab.com/tokend/regources/v2"
)

// NewPublicKeyEntry creates new instance of PublicKeyEntry resource from provided publicKey
func NewPublicKeyEntry(publicKey string) regources.PublicKeyEntry {
	return regources.PublicKeyEntry{
		Key: regources.Key{
			ID:   publicKey,
			Type: regources.TypePublicKeyEntries,
		},
	}
}
