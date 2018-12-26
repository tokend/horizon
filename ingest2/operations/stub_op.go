package operations

import (
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/go/xdr"
)

type stubOpHandler struct {
}

func (h *stubOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr) (history2.OperationDetails, error) {
	return history2.OperationDetails{
		Type: op.Body.Type,
	}, nil
}

func (h *stubOpHandler) ParticipantsEffects(opBody xdr.OperationBody, opRes xdr.OperationResultTr,
	source history2.ParticipantEffect, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{source}, nil
}
