package resources

import (
	"errors"

	"strconv"

	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/rgenerated"
)

func newManageAccountRole(id int64, details history2.ManageAccountRoleDetails,
) *rgenerated.ManageAccountRoleOp {
	switch details.Action {
	case xdr.ManageAccountRoleActionCreate:
		return &rgenerated.ManageAccountRoleOp{
			Key: rgenerated.NewKeyInt64(id, rgenerated.OPERATIONS_CREATE_ACCOUNT_ROLE),
			Attributes: &rgenerated.ManageAccountRoleOpAttributes{
				Details: details.CreateDetails.Details,
			},
			Relationships: rgenerated.ManageAccountRoleOpRelationships{
				Rules: idsAsRelations(details.CreateDetails.RuleIDs, rgenerated.ACCOUNT_RULES),
				Role:  NewAccountRoleKey(details.RoleID).AsRelation(),
			},
		}
	case xdr.ManageAccountRoleActionUpdate:
		return &rgenerated.ManageAccountRoleOp{
			Key: rgenerated.NewKeyInt64(id, rgenerated.OPERATIONS_UPDATE_ACCOUNT_ROLE),
			Attributes: &rgenerated.ManageAccountRoleOpAttributes{
				Details: details.UpdateDetails.Details,
			},
			Relationships: rgenerated.ManageAccountRoleOpRelationships{
				Rules: idsAsRelations(details.UpdateDetails.RuleIDs, rgenerated.ACCOUNT_RULES),
				Role:  NewAccountRoleKey(details.RoleID).AsRelation(),
			},
		}
	case xdr.ManageAccountRoleActionRemove:
		return &rgenerated.ManageAccountRoleOp{
			Key: rgenerated.NewKeyInt64(id, rgenerated.OPERATIONS_REMOVE_ACCOUNT_ROLE),
			Relationships: rgenerated.ManageAccountRoleOpRelationships{
				Role: NewAccountRoleKey(details.RoleID).AsRelation(),
			},
		}
	default:
		panic(errors.New("unexpected manage account role action"))
	}
}

func idsAsRelations(ids []uint64, resourceType rgenerated.ResourceType) *rgenerated.RelationCollection {
	keys := make([]rgenerated.Key, 0, len(ids))
	for _, id := range ids {
		keys = append(keys, rgenerated.Key{
			ID:   strconv.FormatUint(id, 10),
			Type: resourceType,
		})
	}

	return &rgenerated.RelationCollection{
		Data: keys,
	}
}
