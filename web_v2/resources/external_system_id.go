package resources

import (
	"github.com/spf13/cast"
	"gitlab.com/tokend/horizon/db2/core2"
	regources "gitlab.com/tokend/regources/v2/generated"
)

// NewExternalSystemID creates new instance of ExternalSystemID from provided one
func NewExternalSystemID(extSysID core2.ExternalSystemID) *regources.ExternalSystemId {
	return &regources.ExternalSystemId{
		Key: regources.Key{
			ID:   cast.ToString(extSysID.ID),
			Type: regources.EXTERNAL_SYSTEM_IDS,
		},
		Attributes: regources.ExternalSystemIdAttributes{
			ExternalSystemType: extSysID.ExternalSystemType,
			Data:               extSysID.Data,
			IsDeleted:          extSysID.IsDeleted,
			ExpiresAt:          extSysID.ExpiresAt,
			BindedAt:           extSysID.BindedAt,
		},
	}
}
