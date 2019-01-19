package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type stubOpHandler struct {
}

//Details - used as temporary solution for not handled operations
func (h *stubOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr) (history2.OperationDetails, error) {
	return history2.OperationDetails{
		Type: op.Body.Type,
	}, nil
}

//ParticipantsEffects - used as temroary solution for not handled operations. Returns only source participant
func (h *stubOpHandler) ParticipantsEffects(opBody xdr.OperationBody, opRes xdr.OperationResultTr,
	source history2.ParticipantEffect, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{source}, nil
}
