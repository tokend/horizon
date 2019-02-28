package resources

import (
	"github.com/spf13/cast"
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/regources/v2"
)

// NewExternalSystemID creates new instance of ExternalSystemID from provided one
func NewExternalSystemID(extSysID core2.ExternalSystemID) *regources.ExternalSystemID {
	return &regources.ExternalSystemID{
		Key: regources.Key{
			ID:   cast.ToString(extSysID.ID),
			Type: regources.TypeExternalSystemID,
		},
		Attributes: regources.ExternalSystemIDAttr{
			ExternalSystemType: extSysID.ExternalSystemType,
			Data:               extSysID.Data,
			IsDeleted:          extSysID.IsDeleted,
			ExpiresAt:          extSysID.ExpiresAt,
			BindedAt:           extSysID.BindedAt,
		},
		//Relationships: regources.ExternalSystemIDRelations{
		//	Account: NewAccountKey(extSysID.AccountID).AsRelation(),
		//},
	}
}
