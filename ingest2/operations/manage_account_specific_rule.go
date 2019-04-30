package operations

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type manageAccountSpecificRuleOpHandler struct {
	effectsProvider
}

// Details returns details about bind external system account operation
func (h *manageAccountSpecificRuleOpHandler) Details(op rawOperation,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	manageRuleOp := op.Body.MustManageAccountSpecificRuleOp()

	opDetails := history2.OperationDetails{
		Type: xdr.OperationTypeManageAccountSpecificRule,
		ManageAccountSpecificRule: &history2.ManageAccountSpecificRuleDetails{
			Action: manageRuleOp.Data.Action,
		},
	}

	switch manageRuleOp.Data.Action {
	case xdr.ManageAccountSpecificRuleActionCreate:
		opDetails.ManageAccountSpecificRule.RuleID =
			uint64(opRes.MustManageAccountSpecificRuleResult().MustSuccess().RuleId)
		details := manageRuleOp.Data.MustCreateData()

		creationDetails := &history2.CreateAccountSpecificRuleDetails{
			LedgerKey: details.LedgerKey,
			Forbids:   details.Forbids,
			AccountID: details.AccountId.Address(), // Address() smart
		}

		opDetails.ManageAccountSpecificRule.CreateDetails = creationDetails
	case xdr.ManageAccountSpecificRuleActionRemove:
		opDetails.ManageAccountSpecificRule.RuleID = uint64(manageRuleOp.Data.MustRemoveData().RuleId)
	default:
		return history2.OperationDetails{}, errors.New("Unexpected action on manage account specific rule")
	}

	return opDetails, nil
}
