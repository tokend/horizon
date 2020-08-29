package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ingest2/internal"
)

type createDataRemoveRequestHandler struct {
	effectsProvider
}

func (h *createDataRemoveRequestHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	oper := op.Body.MustCreateDataRemoveRequestOp()
	response := opRes.MustCreateDataRemoveRequestResult().Success

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreateDataRemoveRequest,
		CreateDataRemoveRequest: &history2.CreateDataRemoveRequest{
			ID:             uint64(oper.DataRemoveRequest.Id),
			RequestID:      uint64(response.RequestId),
			CreatorDetails: internal.MarshalCustomDetails(oper.DataRemoveRequest.CreatorDetails),
		},
	}, nil
}

func (h *createDataRemoveRequestHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	participants := []history2.ParticipantEffect{h.Participant(sourceAccountID)}

	return participants, nil
}
