package operaitons

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type manageExternalSystemPoolOpHandler struct {
}

// OperationDetails returns details about manage external system pool operation
func (h *manageExternalSystemPoolOpHandler) OperationDetails(op rawOperation,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	manageExternalSystemPoolOp := op.Body.MustManageExternalSystemAccountIdPoolEntryOp()

	operationDetails := history2.OperationDetails{
		Type: xdr.OperationTypeManageExternalSystemAccountIdPoolEntry,
		ManageExternalSystemPool: &history2.ManageExternalSystemPoolDetails{
			Action: manageExternalSystemPoolOp.ActionInput.Action,
		},
	}

	switch operationDetails.ManageExternalSystemPool.Action {
	case xdr.ManageExternalSystemAccountIdPoolEntryActionCreate:
		creationDetails := manageExternalSystemPoolOp.ActionInput.MustCreateExternalSystemAccountIdPoolEntryActionInput()

		operationDetails.ManageExternalSystemPool.Create = &history2.CreateExternalSystemPoolDetails{
			Data:               string(creationDetails.Data),
			ExternalSystemType: int32(creationDetails.ExternalSystemType),
			Parent:             uint64(creationDetails.Parent),
			PoolID:             uint64(opRes.MustManageExternalSystemAccountIdPoolEntryResult().MustSuccess().PoolEntryId),
		}
	case xdr.ManageExternalSystemAccountIdPoolEntryActionRemove:
		operationDetails.ManageExternalSystemPool.Remove = &history2.RemoveExternalSystemPoolDetails{
			PoolID: uint64(manageExternalSystemPoolOp.ActionInput.MustDeleteExternalSystemAccountIdPoolEntryActionInput().PoolEntryId),
		}
	default:
		return history2.OperationDetails{}, errors.From(
			errors.New("unexpected manage external system pool action"), map[string]interface{}{
				"action": int32(operationDetails.ManageExternalSystemPool.Action),
			})
	}

	return operationDetails, nil
}

func (h *manageExternalSystemPoolOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{source}, nil
}
