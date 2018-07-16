package resource

import (
	"time"

	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/tokend/regources"
)

type ExternalSystemAccountID struct {
	Type      regources.Flag `json:"type"`
	Data      string         `json:"data"`
	AssetCode string         `json:"asset_code,omitempty"`
	ExpiresAt *string        `json:"expires_at,omitempty"`
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

	// check out actions_account.go to find ExpiresAt default value
	if coreRecord.ExpiresAt != nil {
		expiresAt := time.Unix(*coreRecord.ExpiresAt, 0).Format(time.RFC3339)
		id.ExpiresAt = &expiresAt
	}
}
