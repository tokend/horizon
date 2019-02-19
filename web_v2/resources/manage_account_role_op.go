package resources

import (
	"errors"

	"strconv"

	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/v2"
)

func newManageAccountRole(id int64, details history2.ManageAccountRoleDetails,
) *regources.ManageAccountRole {
	switch details.Action {
	case xdr.ManageAccountRoleActionCreate:
		return &regources.ManageAccountRole{
			Key: regources.NewKeyInt64(id, regources.TypeCreateAccountRole),
			Attributes: &regources.ManageAccountRoleAttrs{
				Details: details.CreateDetails.Details,
			},
			Relationships: regources.ManageAccountRoleRelation{
				Rules: idsAsRelations(details.CreateDetails.RuleIDs, regources.TypeAccountRules),
				Role:  NewAccountRoleKey(details.RoleID).AsRelation(),
			},
		}
	case xdr.ManageAccountRoleActionUpdate:
		return &regources.ManageAccountRole{
			Key: regources.NewKeyInt64(id, regources.TypeUpdateAccountRole),
			Attributes: &regources.ManageAccountRoleAttrs{
				Details: details.UpdateDetails.Details,
			},
			Relationships: regources.ManageAccountRoleRelation{
				Rules: idsAsRelations(details.UpdateDetails.RuleIDs, regources.TypeAccountRules),
				Role:  NewAccountRoleKey(details.RoleID).AsRelation(),
			},
		}
	case xdr.ManageAccountRoleActionRemove:
		return &regources.ManageAccountRole{
			Key: regources.NewKeyInt64(id, regources.TypeRemoveAccountRole),
			Relationships: regources.ManageAccountRoleRelation{
				Role: NewAccountRoleKey(details.RoleID).AsRelation(),
			},
		}
	default:
		panic(errors.New("unexpected manage account role action"))
	}
}

func idsAsRelations(ids []uint64, resourceType regources.ResourceType) *regources.RelationCollection {
	keys := make([]regources.Key, 0, len(ids))
	for _, id := range ids {
		keys = append(keys, regources.Key{
			ID:   strconv.FormatUint(id, 10),
			Type: resourceType,
		})
	}

	return &regources.RelationCollection{
		Data: keys,
	}
}
