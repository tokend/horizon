package operations

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/v2"
)

type manageExternalSystemPoolOpHandler struct {
}

// Details returns details about manage external system pool operation
func (h *manageExternalSystemPoolOpHandler) Details(op rawOperation,
	opRes xdr.OperationResultTr,
) (regources.OperationDetails, error) {
	manageExternalSystemPoolOp := op.Body.MustManageExternalSystemAccountIdPoolEntryOp()

	operationDetails := regources.OperationDetails{
		Type: xdr.OperationTypeManageExternalSystemAccountIdPoolEntry,
		ManageExternalSystemPool: &regources.ManageExternalSystemPoolDetails{
			Action: manageExternalSystemPoolOp.ActionInput.Action,
		},
	}

	switch operationDetails.ManageExternalSystemPool.Action {
	case xdr.ManageExternalSystemAccountIdPoolEntryActionCreate:
		creationDetails := manageExternalSystemPoolOp.ActionInput.
			MustCreateExternalSystemAccountIdPoolEntryActionInput()

		operationDetails.ManageExternalSystemPool.Create = &regources.CreateExternalSystemPoolDetails{
			Data:               string(creationDetails.Data),
			ExternalSystemType: int32(creationDetails.ExternalSystemType),
			Parent:             uint64(creationDetails.Parent),
			PoolID: uint64(opRes.MustManageExternalSystemAccountIdPoolEntryResult().MustSuccess().
				PoolEntryId),
		}
	case xdr.ManageExternalSystemAccountIdPoolEntryActionRemove:
		operationDetails.ManageExternalSystemPool.Remove = &regources.RemoveExternalSystemPoolDetails{
			PoolID: uint64(manageExternalSystemPoolOp.ActionInput.
				MustDeleteExternalSystemAccountIdPoolEntryActionInput().PoolEntryId),
		}
	default:
		return regources.OperationDetails{}, errors.From(
			errors.New("unexpected manage external system pool action"), map[string]interface{}{
				"action": int32(operationDetails.ManageExternalSystemPool.Action),
			})
	}

	return operationDetails, nil
}

//ParticipantsEffects - returns source of the operation
func (h *manageExternalSystemPoolOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	_ xdr.OperationResultTr, source history2.ParticipantEffect, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{source}, nil
}
