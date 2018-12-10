package operaitons

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type manageExternalSystemPoolOpHandler struct {
}

// OperationDetails returns details about manage external system pool operation
func (h *manageExternalSystemPoolOpHandler) OperationDetails(op rawOperation,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	//manageExternalSystemPoolOp := opBody.MustManageExternalSystemAccountIdPoolEntryOp()

	return history2.OperationDetails{}, nil
}

func (h *manageExternalSystemPoolOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{source}, nil
}
