package resource

import (
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/swarmfund/horizon/resource/base"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/regources/valueflag"
)

// Signer represents one of an account's signers.
type Signer struct {
	PublicKey      string           `json:"public_key"`
	Weight         int32            `json:"weight"`
	SignerTypeI    int32            `json:"signer_type_i"`
	SignerTypes    []valueflag.Flag `json:"signer_types"`
	SignerIdentity int32            `json:"signer_identity"`
	SignerName     string           `json:"signer_name"`
}

// Populate fills out the fields of the signer, using one of an account's
// secondary signers.
func (s *Signer) Populate(row core.Signer) {
	s.populate(row.Publickey, row.Weight, row.SignerType, row.Identity, row.Name)
}

func (s *Signer) populate(publicKey string, weight, rawSignerType, identity int32, name string) {
	s.PublicKey = publicKey
	s.Weight = weight
	s.SignerTypeI = rawSignerType
	s.SignerIdentity = identity
	s.SignerName = name
	s.SignerTypes = base.FlagFromXdrSignerType(rawSignerType, xdr.SignerTypeAll)
}

func (s *Signer) FromXDR(xSigner xdr.Signer) {
	s.PublicKey = xSigner.PubKey.Address()
	s.Weight = int32(xSigner.Weight)
	s.SignerTypeI = int32(xSigner.SignerType)
	s.SignerIdentity = int32(xSigner.Identity)
	s.SignerName = string(xSigner.Name)
	s.SignerTypes = base.FlagFromXdrSignerType(s.SignerTypeI, xdr.SignerTypeAll)

}
