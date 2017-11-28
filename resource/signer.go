package resource

import (
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon/db2/core"
)

type SignerType struct {
	Name  string `json:"name"`
	Value int32  `json:"value"`
}

// Signer represents one of an account's signers.
type Signer struct {
	PublicKey      string       `json:"public_key"`
	Weight         int32        `json:"weight"`
	SignerTypeI    int32        `json:"signer_type_i"`
	SignerTypes    []SignerType `json:"signer_types"`
	SignerIdentity int32        `json:"signer_identity"`
	SignerName     string       `json:"signer_name"`
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
	for _, signerType := range xdr.SignerTypeAll {
		if (int32(signerType) & s.SignerTypeI) != 0 {
			s.SignerTypes = append(s.SignerTypes, SignerType{
				Value: int32(signerType),
				Name:  signerType.String(),
			})
		}
	}
}
