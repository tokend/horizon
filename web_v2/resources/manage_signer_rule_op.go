package resources

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/v2"
)

func newManageSignerRule(id int64, details history2.ManageSignerRuleDetails,
) *regources.ManageSignerRule {
	switch details.Action {
	case xdr.ManageSignerRuleActionCreate:
		return &regources.ManageSignerRule{
			Key: regources.NewKeyInt64(id, regources.TypeCreateSignerRule),
			Attributes: &regources.ManageSignerRuleAttrs{
				Resource:   details.CreateDetails.Resource,
				Action:     details.CreateDetails.Action,
				IsForbid:   details.CreateDetails.IsForbid,
				IsDefault:  details.CreateDetails.IsDefault,
				IsReadOnly: details.CreateDetails.IsReadOnly,
				Details:    details.CreateDetails.Details,
			},
			Relationships: &regources.ManageSignerRuleRelation{
				Rule: NewSignerRoleKey(details.RuleID).AsRelation(),
			},
		}
	case xdr.ManageSignerRuleActionUpdate:
		return &regources.ManageSignerRule{
			Key: regources.NewKeyInt64(id, regources.TypeUpdateSignerRule),
			Attributes: &regources.ManageSignerRuleAttrs{
				Resource:  details.UpdateDetails.Resource,
				Action:    details.UpdateDetails.Action,
				IsForbid:  details.UpdateDetails.IsForbid,
				IsDefault: details.UpdateDetails.IsDefault,
				Details:   details.UpdateDetails.Details,
			},
			Relationships: &regources.ManageSignerRuleRelation{
				Rule: NewSignerRoleKey(details.RuleID).AsRelation(),
			},
		}
	case xdr.ManageSignerRuleActionRemove:
		return &regources.ManageSignerRule{
			Key: regources.NewKeyInt64(id, regources.TypeRemoveSignerRule),
			Relationships: &regources.ManageSignerRuleRelation{
				Rule: NewSignerRoleKey(details.RuleID).AsRelation(),
			},
		}
	default:
		panic(errors.New("unexpected manage account role action"))
	}
}
