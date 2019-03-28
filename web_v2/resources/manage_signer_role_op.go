package resources

import (
	"errors"

	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	regources "gitlab.com/tokend/regources/v2/generated"
)

func newManageSignerRole(id int64, details history2.ManageSignerRoleDetails,
) *regources.ManageSignerRoleOp {
	switch details.Action {
	case xdr.ManageSignerRoleActionCreate:
		return &regources.ManageSignerRoleOp{
			Key: regources.NewKeyInt64(id, regources.OPERATIONS_CREATE_SIGNER_ROLE),
			Attributes: &regources.ManageSignerRoleOpAttributes{
				Details:    details.CreateDetails.Details,
				IsReadOnly: details.CreateDetails.IsReadOnly,
			},
			Relationships: regources.ManageSignerRoleOpRelationships{
				Rules: idsAsRelations(details.CreateDetails.RuleIDs, regources.SIGNER_RULES),
				Role:  NewSignerRoleKey(details.RoleID).AsRelation(),
			},
		}
	case xdr.ManageSignerRoleActionUpdate:
		return &regources.ManageSignerRoleOp{
			Key: regources.NewKeyInt64(id, regources.OPERATIONS_UPDATE_SIGNER_ROLE),
			Attributes: &regources.ManageSignerRoleOpAttributes{
				Details: details.UpdateDetails.Details,
			},
			Relationships: regources.ManageSignerRoleOpRelationships{
				Rules: idsAsRelations(details.UpdateDetails.RuleIDs, regources.SIGNER_RULES),
				Role:  NewSignerRoleKey(details.RoleID).AsRelation(),
			},
		}
	case xdr.ManageSignerRoleActionRemove:
		return &regources.ManageSignerRoleOp{
			Key: regources.NewKeyInt64(id, regources.OPERATIONS_REMOVE_SIGNER_ROLE),
			Relationships: regources.ManageSignerRoleOpRelationships{
				Role: NewSignerRoleKey(details.RoleID).AsRelation(),
			},
		}
	default:
		panic(errors.New("unexpected manage account role action"))
	}
}
