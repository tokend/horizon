package resource

import (
	"gitlab.com/distributed_lab/tokend/horizon/db2/core"
)

type Trusts struct {
	Trusts []Trust `json:"trusts"`
}

func (ct *Trusts) Populate(trusts []core.Trust) {
	ct.Trusts = make([]Trust, len(trusts))
	for i := range trusts {
		ct.Trusts[i].Populate(trusts[i])
	}
}
