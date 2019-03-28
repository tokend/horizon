package resources

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	regources "gitlab.com/tokend/regources/v2/generated"
)

func newManageAccountRule(id int64, details history2.ManageAccountRuleDetails,
) *regources.ManageAccountRuleOp {
	switch details.Action {
	case xdr.ManageAccountRuleActionCreate:
		return &regources.ManageAccountRuleOp{
			Key: regources.NewKeyInt64(id, regources.OPERATIONS_CREATE_ACCOUNT_RULE),
			Attributes: &regources.ManageAccountRuleOpAttributes{
				Resource: details.CreateDetails.Resource,
				Action:   details.CreateDetails.Action,
				Forbids:  details.CreateDetails.IsForbid,
				Details:  details.CreateDetails.Details,
			},
			Relationships: &regources.ManageAccountRuleOpRelationships{
				Rule: NewAccountRoleKey(details.RuleID).AsRelation(),
			},
		}
	case xdr.ManageAccountRuleActionUpdate:
		return &regources.ManageAccountRuleOp{
			Key: regources.NewKeyInt64(id, regources.OPERATIONS_UPDATE_ACCOUNT_RULE),
			Attributes: &regources.ManageAccountRuleOpAttributes{
				Resource: details.UpdateDetails.Resource,
				Action:   details.UpdateDetails.Action,
				Forbids:  details.UpdateDetails.IsForbid,
				Details:  details.UpdateDetails.Details,
			},
			Relationships: &regources.ManageAccountRuleOpRelationships{
				Rule: NewAccountRoleKey(details.RuleID).AsRelation(),
			},
		}
	case xdr.ManageAccountRuleActionRemove:
		return &regources.ManageAccountRuleOp{
			Key: regources.NewKeyInt64(id, regources.OPERATIONS_REMOVE_ACCOUNT_RULE),
			Relationships: &regources.ManageAccountRuleOpRelationships{
				Rule: NewAccountRoleKey(details.RuleID).AsRelation(),
			},
		}
	default:
		panic(errors.New("unexpected manage account role action"))
	}
}
