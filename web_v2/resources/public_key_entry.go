package resources

import (
	"gitlab.com/tokend/regources/rgenerated"
)

// NewPublicKeyEntry creates new instance of PublicKeyEntry resource from provided publicKey
func NewPublicKeyEntry(publicKey string) rgenerated.PublicKeyEntry {
	return rgenerated.PublicKeyEntry{
		Key: rgenerated.Key{
			ID:   publicKey,
			Type: rgenerated.PUBLIC_KEY_ENTRIES,
		},
	}
}
