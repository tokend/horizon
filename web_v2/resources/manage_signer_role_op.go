package resources

import (
	"errors"

	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/v2"
)

func newManageSignerRole(id int64, details history2.ManageSignerRoleDetails,
) *regources.ManageSignerRole {
	switch details.Action {
	case xdr.ManageSignerRoleActionCreate:
		return &regources.ManageSignerRole{
			Key: regources.NewKeyInt64(id, regources.TypeCreateSignerRole),
			Attributes: &regources.ManageSignerRoleAttrs{
				Details:    details.CreateDetails.Details,
				IsReadOnly: details.CreateDetails.IsReadOnly,
			},
			Relationships: regources.ManageSignerRoleRelation{
				Rules: idsAsRelations(details.CreateDetails.RuleIDs, regources.TypeSignerRules),
				Role:  NewSignerRoleKey(details.RoleID).AsRelation(),
			},
		}
	case xdr.ManageSignerRoleActionUpdate:
		return &regources.ManageSignerRole{
			Key: regources.NewKeyInt64(id, regources.TypeUpdateSignerRole),
			Attributes: &regources.ManageSignerRoleAttrs{
				Details: details.UpdateDetails.Details,
			},
			Relationships: regources.ManageSignerRoleRelation{
				Rules: idsAsRelations(details.UpdateDetails.RuleIDs, regources.TypeSignerRules),
				Role:  NewSignerRoleKey(details.RoleID).AsRelation(),
			},
		}
	case xdr.ManageSignerRoleActionRemove:
		return &regources.ManageSignerRole{
			Key: regources.NewKeyInt64(id, regources.TypeRemoveSignerRole),
			Relationships: regources.ManageSignerRoleRelation{
				Role: NewSignerRoleKey(details.RoleID).AsRelation(),
			},
		}
	default:
		panic(errors.New("unexpected manage account role action"))
	}
}
