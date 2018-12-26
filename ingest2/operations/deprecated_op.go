package operations

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type deprecatedOpHandler struct {
}

func (h *deprecatedOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr) (history2.OperationDetails, error) {
	return history2.OperationDetails{}, errors.From(errors.New("Tried to ingest deprecated operation"), logan.F{
		"op_type": op.Body.Type,
	})
}

func (h *deprecatedOpHandler) ParticipantsEffects(opBody xdr.OperationBody, opRes xdr.OperationResultTr,
source history2.ParticipantEffect, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	return nil, errors.From(errors.New("Tried to ingest deprecated operation"), logan.F{
		"op_type": opBody.Type,
	})
}
