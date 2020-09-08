package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ingest2/internal"
)

type manageUpdateDataOpHandler struct {
	effectsProvider
}

func (h *manageUpdateDataOpHandler) ParticipantsEffects(opBody xdr.OperationBody, opRes xdr.OperationResultTr,
	sourceID xdr.AccountId, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{h.Participant(sourceID)}, nil
}

func (h *manageUpdateDataOpHandler) Details(op rawOperation, _ xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	updateDataOp := op.Body.MustUpdateDataOp()

	return history2.OperationDetails{
		Type: xdr.OperationTypeUpdateData,
		UpdateData: &history2.UpdateDataDetails{
			Value: internal.MarshalCustomDetails(updateDataOp.Value),
			ID:    uint64(updateDataOp.DataId),
		},
	}, nil
}
