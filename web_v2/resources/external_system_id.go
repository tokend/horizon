package resources

import (
	"github.com/spf13/cast"
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/regources/rgenerated"
)

// NewExternalSystemID creates new instance of ExternalSystemID from provided one
func NewExternalSystemID(extSysID core2.ExternalSystemID) *rgenerated.ExternalSystemId {
	return &rgenerated.ExternalSystemId{
		Key: rgenerated.Key{
			ID:   cast.ToString(extSysID.ID),
			Type: rgenerated.EXTERNAL_SYSTEM_IDS,
		},
		Attributes: rgenerated.ExternalSystemIdAttributes{
			ExternalSystemType: extSysID.ExternalSystemType,
			Data:               extSysID.Data,
			IsDeleted:          extSysID.IsDeleted,
			ExpiresAt:          extSysID.ExpiresAt,
			BindedAt:           extSysID.BindedAt,
		},
	}
}
