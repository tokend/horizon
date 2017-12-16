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
	id.Type.Name = coreRecord.ExternalSystemType.ShortString()
	id.Type.Value = int32(coreRecord.ExternalSystemType)
	id.Data = coreRecord.Data
}
