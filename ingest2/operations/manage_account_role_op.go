package operations

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ingest2/internal"
)

type manageAccountRoleOpHandler struct {
	effectsProvider
}

// Details returns details about bind external system account operation
func (h *manageAccountRoleOpHandler) Details(op rawOperation,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	manageAccountRoleOp := op.Body.MustManageAccountRoleOp()

	opDetails := history2.OperationDetails{
		Type: xdr.OperationTypeManageAccountRole,
		ManageAccountRole: &history2.ManageAccountRoleDetails{
			Action: manageAccountRoleOp.Data.Action,
		},
	}

	switch manageAccountRoleOp.Data.Action {
	case xdr.ManageAccountRoleActionCreate:
		opDetails.ManageAccountRole.RoleID =
			uint64(opRes.MustManageAccountRoleResult().MustSuccess().RoleId)
		details := manageAccountRoleOp.Data.MustCreateData()

		creationDetails := &history2.UpdateAccountRoleDetails{
			Details: internal.MarshalCustomDetails(details.Details),
		}

		for _, id := range details.AccountRuleIDs {
			creationDetails.RuleIDs = append(creationDetails.RuleIDs, uint64(id))
		}

		opDetails.ManageAccountRole.CreateDetails = creationDetails
	case xdr.ManageAccountRoleActionUpdate:
		details := manageAccountRoleOp.Data.MustUpdateData()
		opDetails.ManageAccountRule.RuleID = uint64(details.RoleId)

		updateDetails := &history2.UpdateAccountRoleDetails{
			Details: internal.MarshalCustomDetails(details.Details),
		}

		for _, id := range details.AccountRuleIDs {
			updateDetails.RuleIDs = append(updateDetails.RuleIDs, uint64(id))
		}

		opDetails.ManageAccountRole.UpdateDetails = updateDetails
	case xdr.ManageAccountRoleActionRemove:
		opDetails.ManageAccountRule.RuleID = uint64(manageAccountRoleOp.Data.MustRemoveData().AccountRoleId)
	default:
		return history2.OperationDetails{}, errors.New("Unexpected action on manage account role")
	}

	return opDetails, nil
}

// ParticipantsEffects returns only source without effects
func (h *manageAccountRoleOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{h.Participant(sourceAccountID)}, nil
}
