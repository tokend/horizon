package resource

import (
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/swarmfund/horizon/resource/base"
)

type ExternalSystemAccountID struct {
	Type      base.Flag `json:"type"`
	Data      string    `json:"data"`
	AssetCode string    `json:"asset_code"`
}

func (id *ExternalSystemAccountID) Populate(coreRecord core.ExternalSystemAccountID) {
	switch coreRecord.ExternalSystemType {
	case 1:
		id.Type.Name = "Bitcoin"
		id.AssetCode = "BTC"
	case 2:
		id.Type.Name = "Ethereum"
		id.AssetCode = "ETH"
	}
	id.Type.Value = int32(coreRecord.ExternalSystemType)
	id.Data = coreRecord.Data
}
