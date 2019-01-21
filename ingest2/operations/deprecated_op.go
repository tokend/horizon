package operations

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/v2"
)

type deprecatedOpHandler struct {
}

//Details - returns error as deprecated op should never occur
func (h *deprecatedOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr) (regources.OperationDetails, error) {
	return regources.OperationDetails{}, errors.From(errors.New("Tried to ingest deprecated operation"), logan.F{
		"op_type": op.Body.Type,
	})
}

//ParticipantsEffects - returns errors as deprecated op should never occur
func (h *deprecatedOpHandler) ParticipantsEffects(opBody xdr.OperationBody, opRes xdr.OperationResultTr,
	source history2.ParticipantEffect, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	return nil, errors.From(errors.New("Tried to ingest deprecated operation"), logan.F{
		"op_type": opBody.Type,
	})
}
