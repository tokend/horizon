package resources

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	regources "gitlab.com/tokend/regources/generated"
)

func newManageSignerRule(id int64, details history2.ManageSignerRuleDetails,
) *regources.ManageSignerRuleOp {
	switch details.Action {
	case xdr.ManageSignerRuleActionCreate:
		return &regources.ManageSignerRuleOp{
			Key: regources.NewKeyInt64(id, regources.OPERATIONS_CREATE_SIGNER_RULE),
			Attributes: &regources.ManageSignerRuleOpAttributes{
				Resource:   details.CreateDetails.Resource,
				Action:     details.CreateDetails.Action,
				Forbids:    details.CreateDetails.IsForbid,
				IsDefault:  details.CreateDetails.IsDefault,
				IsReadOnly: details.CreateDetails.IsReadOnly,
				Details:    details.CreateDetails.Details,
			},
			Relationships: &regources.ManageSignerRuleOpRelationships{
				Rule: NewSignerRoleKey(details.RuleID).AsRelation(),
			},
		}
	case xdr.ManageSignerRuleActionUpdate:
		return &regources.ManageSignerRuleOp{
			Key: regources.NewKeyInt64(id, regources.OPERATIONS_UPDATE_SIGNER_RULE),
			Attributes: &regources.ManageSignerRuleOpAttributes{
				Resource:  details.UpdateDetails.Resource,
				Action:    details.UpdateDetails.Action,
				Forbids:   details.UpdateDetails.IsForbid,
				IsDefault: details.UpdateDetails.IsDefault,
				Details:   details.UpdateDetails.Details,
			},
			Relationships: &regources.ManageSignerRuleOpRelationships{
				Rule: NewSignerRoleKey(details.RuleID).AsRelation(),
			},
		}
	case xdr.ManageSignerRuleActionRemove:
		return &regources.ManageSignerRuleOp{
			Key: regources.NewKeyInt64(id, regources.OPERATIONS_REMOVE_SIGNER_RULE),
			Relationships: &regources.ManageSignerRuleOpRelationships{
				Rule: NewSignerRoleKey(details.RuleID).AsRelation(),
			},
		}
	default:
		panic(errors.New("unexpected manage account role action"))
	}
}
