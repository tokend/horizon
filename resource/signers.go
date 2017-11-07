package resource

import (
	"bullioncoin.githost.io/development/horizon/db2/core"
)

type Signers struct {
	Signers []Signer `json:"signers"`
}

func (s *Signers) Populate(signers []core.Signer) {
	s.Signers = make([]Signer, len(signers))
	for i := range signers {
		s.Signers[i].Populate(signers[i])
	}
}
