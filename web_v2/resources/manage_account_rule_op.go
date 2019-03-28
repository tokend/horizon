package resources

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/rgenerated"
)

func newManageAccountRule(id int64, details history2.ManageAccountRuleDetails,
) *rgenerated.ManageAccountRuleOp {
	switch details.Action {
	case xdr.ManageAccountRuleActionCreate:
		return &rgenerated.ManageAccountRuleOp{
			Key: rgenerated.NewKeyInt64(id, rgenerated.OPERATIONS_CREATE_ACCOUNT_RULE),
			Attributes: &rgenerated.ManageAccountRuleOpAttributes{
				Resource: details.CreateDetails.Resource,
				Action:   details.CreateDetails.Action,
				Forbids:  details.CreateDetails.IsForbid,
				Details:  details.CreateDetails.Details,
			},
			Relationships: &rgenerated.ManageAccountRuleOpRelationships{
				Rule: NewAccountRoleKey(details.RuleID).AsRelation(),
			},
		}
	case xdr.ManageAccountRuleActionUpdate:
		return &rgenerated.ManageAccountRuleOp{
			Key: rgenerated.NewKeyInt64(id, rgenerated.OPERATIONS_UPDATE_ACCOUNT_RULE),
			Attributes: &rgenerated.ManageAccountRuleOpAttributes{
				Resource: details.UpdateDetails.Resource,
				Action:   details.UpdateDetails.Action,
				Forbids:  details.UpdateDetails.IsForbid,
				Details:  details.UpdateDetails.Details,
			},
			Relationships: &rgenerated.ManageAccountRuleOpRelationships{
				Rule: NewAccountRoleKey(details.RuleID).AsRelation(),
			},
		}
	case xdr.ManageAccountRuleActionRemove:
		return &rgenerated.ManageAccountRuleOp{
			Key: rgenerated.NewKeyInt64(id, rgenerated.OPERATIONS_REMOVE_ACCOUNT_RULE),
			Relationships: &rgenerated.ManageAccountRuleOpRelationships{
				Rule: NewAccountRoleKey(details.RuleID).AsRelation(),
			},
		}
	default:
		panic(errors.New("unexpected manage account role action"))
	}
}
