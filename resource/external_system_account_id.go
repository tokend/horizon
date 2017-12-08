package resource

import (
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon/db2/core"
)

type ExternalSystemAccountID struct {
	Type xdr.ExternalSystemType `json:"type"`
	Data string                 `json:"data"`
}

func (id *ExternalSystemAccountID) Populate(coreRecord core.ExternalSystemAccountID) {
	id.Type = coreRecord.ExternalSystemType
	id.Data = coreRecord.Data
}
