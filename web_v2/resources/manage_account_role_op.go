package resources

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/v2"
)

func newManageAccountRole(id int64, details history2.ManageAccountRoleDetails,
) *regources.ManageAccountRole {
	result := &regources.ManageAccountRole{
		Key: regources.NewKeyInt64(id, regources.TypeManageAccountRole),
		Attributes: regources.ManageAccountRoleAttrs{
			Action: details.Action,
			RoleID: details.RoleID,
		},
	}

	switch details.Action {
	case xdr.ManageAccountRoleActionCreate:
		createDetails := regources.UpdateAccountRoleAttrs(*details.CreateDetails)
		result.Attributes.CreateAttrs = &createDetails
	case xdr.ManageAccountRoleActionUpdate:
		updateDetails := regources.UpdateAccountRoleAttrs(*details.UpdateDetails)
		result.Attributes.UpdateAttrs = &updateDetails
	}

	return result
}
