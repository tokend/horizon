package resources

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/v2"
)

func newManageAccountRule(id int64, details history2.ManageAccountRuleDetails,
) *regources.ManageAccountRule {
	switch details.Action {
	case xdr.ManageAccountRuleActionCreate:
		return &regources.ManageAccountRule{
			Key: regources.NewKeyInt64(id, regources.TypeCreateAccountRule),
			Attributes: &regources.ManageAccountRuleAttrs{
				Resource: details.CreateDetails.Resource,
				Action:   details.CreateDetails.Action,
				IsForbid: details.CreateDetails.IsForbid,
				Details:  details.CreateDetails.Details,
			},
			Relationships: &regources.ManageAccountRuleRelation{
				Rule: NewAccountRoleKey(details.RuleID).AsRelation(),
			},
		}
	case xdr.ManageAccountRuleActionUpdate:
		return &regources.ManageAccountRule{
			Key: regources.NewKeyInt64(id, regources.TypeUpdateAccountRule),
			Attributes: &regources.ManageAccountRuleAttrs{
				Resource: details.UpdateDetails.Resource,
				Action:   details.UpdateDetails.Action,
				IsForbid: details.UpdateDetails.IsForbid,
				Details:  details.UpdateDetails.Details,
			},
			Relationships: &regources.ManageAccountRuleRelation{
				Rule: NewAccountRoleKey(details.RuleID).AsRelation(),
			},
		}
	case xdr.ManageAccountRuleActionRemove:
		return &regources.ManageAccountRule{
			Key: regources.NewKeyInt64(id, regources.TypeRemoveAccountRule),
			Relationships: &regources.ManageAccountRuleRelation{
				Rule: NewAccountRoleKey(details.RuleID).AsRelation(),
			},
		}
	default:
		panic(errors.New("unexpected manage account role action"))
	}
}
