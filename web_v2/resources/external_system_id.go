package resources

import (
	"github.com/spf13/cast"
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/regources/v2"
)

func NewExternalSystemIDs(extSysID core2.ExternalSystemID) *regources.ExternalSystemIDs {
	return &regources.ExternalSystemIDs{
		Key: NewExternalSystemIDsKey(extSysID.ID),
		Attributes: regources.ExternalSystemIDsAttr{
			AccountID:          extSysID.AccountID,
			ExternalSystemType: extSysID.ExternalSystemType,
			Data:               extSysID.Data,
			IsDeleted:          extSysID.IsDeleted,
			ExpiresAt:          extSysID.ExpiresAt,
			BindedAt:           extSysID.BindedAt,
		},
	}
}

func NewExternalSystemIDsKey(extSysIDsID uint64) regources.Key {
	return regources.Key{
		ID:   cast.ToString(extSysIDsID),
		Type: regources.TypeExternalSystemIDs,
	}
}
