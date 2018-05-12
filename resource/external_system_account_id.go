package resource

import (
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/swarmfund/horizon/resource/base"
)

type ExternalSystemAccountID struct {
	Type base.Flag `json:"type"`
	Data string    `json:"data"`
}

func (id *ExternalSystemAccountID) Populate(coreRecord core.ExternalSystemAccountID) {
	switch coreRecord.ExternalSystemType{
	case 1:
		id.Type.Name = "Bitcoin"
	case 2:
		id.Type.Name = "Ethereum"
	}
	id.Type.Value = int32(coreRecord.ExternalSystemType)
	id.Data = coreRecord.Data
}
