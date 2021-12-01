package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ingest2/internal"
)

type manageCreateDataOpHandler struct {
	effectsProvider
}

func (h *manageCreateDataOpHandler) ParticipantsEffects(opBody xdr.OperationBody, opRes xdr.OperationResultTr,
	sourceID xdr.AccountId, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{h.Participant(sourceID)}, nil
}

func (h *manageCreateDataOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	createDataOp := op.Body.MustCreateDataOp()
	res := opRes.MustCreateDataResult().MustSuccess()

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreateData,
		CreateData: &history2.CreateDataDetails{
			Type:  uint64(createDataOp.Type),
			Value: internal.MarshalCustomDetails(createDataOp.Value),
			Owner: op.Source.Address(),
			ID:    uint64(res.DataId),
		},
	}, nil
}
