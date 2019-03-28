package resources

import (
	"errors"

	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/rgenerated"
)

func newManageSignerRole(id int64, details history2.ManageSignerRoleDetails,
) *rgenerated.ManageSignerRoleOp {
	switch details.Action {
	case xdr.ManageSignerRoleActionCreate:
		return &rgenerated.ManageSignerRoleOp{
			Key: rgenerated.NewKeyInt64(id, rgenerated.OPERATIONS_CREATE_SIGNER_ROLE),
			Attributes: &rgenerated.ManageSignerRoleOpAttributes{
				Details:    details.CreateDetails.Details,
				IsReadOnly: details.CreateDetails.IsReadOnly,
			},
			Relationships: rgenerated.ManageSignerRoleOpRelationships{
				Rules: idsAsRelations(details.CreateDetails.RuleIDs, rgenerated.SIGNER_RULES),
				Role:  NewSignerRoleKey(details.RoleID).AsRelation(),
			},
		}
	case xdr.ManageSignerRoleActionUpdate:
		return &rgenerated.ManageSignerRoleOp{
			Key: rgenerated.NewKeyInt64(id, rgenerated.OPERATIONS_UPDATE_SIGNER_ROLE),
			Attributes: &rgenerated.ManageSignerRoleOpAttributes{
				Details: details.UpdateDetails.Details,
			},
			Relationships: rgenerated.ManageSignerRoleOpRelationships{
				Rules: idsAsRelations(details.UpdateDetails.RuleIDs, rgenerated.SIGNER_RULES),
				Role:  NewSignerRoleKey(details.RoleID).AsRelation(),
			},
		}
	case xdr.ManageSignerRoleActionRemove:
		return &rgenerated.ManageSignerRoleOp{
			Key: rgenerated.NewKeyInt64(id, rgenerated.OPERATIONS_REMOVE_SIGNER_ROLE),
			Relationships: rgenerated.ManageSignerRoleOpRelationships{
				Role: NewSignerRoleKey(details.RoleID).AsRelation(),
			},
		}
	default:
		panic(errors.New("unexpected manage account role action"))
	}
}
