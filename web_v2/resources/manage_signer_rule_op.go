package resources

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/rgenerated"
)

func newManageSignerRule(id int64, details history2.ManageSignerRuleDetails,
) *rgenerated.ManageSignerRuleOp {
	switch details.Action {
	case xdr.ManageSignerRuleActionCreate:
		return &rgenerated.ManageSignerRuleOp{
			Key: rgenerated.NewKeyInt64(id, rgenerated.OPERATIONS_CREATE_SIGNER_RULE),
			Attributes: &rgenerated.ManageSignerRuleOpAttributes{
				Resource:   details.CreateDetails.Resource,
				Action:     details.CreateDetails.Action,
				Forbids:    details.CreateDetails.IsForbid,
				IsDefault:  details.CreateDetails.IsDefault,
				IsReadOnly: details.CreateDetails.IsReadOnly,
				Details:    details.CreateDetails.Details,
			},
			Relationships: &rgenerated.ManageSignerRuleOpRelationships{
				Rule: NewSignerRoleKey(details.RuleID).AsRelation(),
			},
		}
	case xdr.ManageSignerRuleActionUpdate:
		return &rgenerated.ManageSignerRuleOp{
			Key: rgenerated.NewKeyInt64(id, rgenerated.OPERATIONS_UPDATE_SIGNER_RULE),
			Attributes: &rgenerated.ManageSignerRuleOpAttributes{
				Resource:  details.UpdateDetails.Resource,
				Action:    details.UpdateDetails.Action,
				Forbids:   details.UpdateDetails.IsForbid,
				IsDefault: details.UpdateDetails.IsDefault,
				Details:   details.UpdateDetails.Details,
			},
			Relationships: &rgenerated.ManageSignerRuleOpRelationships{
				Rule: NewSignerRoleKey(details.RuleID).AsRelation(),
			},
		}
	case xdr.ManageSignerRuleActionRemove:
		return &rgenerated.ManageSignerRuleOp{
			Key: rgenerated.NewKeyInt64(id, rgenerated.OPERATIONS_REMOVE_SIGNER_RULE),
			Relationships: &rgenerated.ManageSignerRuleOpRelationships{
				Rule: NewSignerRoleKey(details.RuleID).AsRelation(),
			},
		}
	default:
		panic(errors.New("unexpected manage account role action"))
	}
}
