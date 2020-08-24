package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ingest2/internal"
)

type createDataUpdateRequestHandler struct {
	effectsProvider
}

func (h *createDataUpdateRequestHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	oper := op.Body.MustCreateDataUpdateRequestOp()
	response := opRes.MustCreateDataUpdateRequestResult().CreateDataUpdateRequestResponse

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreateDataUpdateRequest,
		CreateDataUpdateRequest: &history2.CreateDataUpdateRequest{
			ID:             uint64(oper.DataUpdateRequest.Id),
			Value:          internal.MarshalCustomDetails(response.Value),
			RequestID:      uint64(response.RequestId),
			CreatorDetails: internal.MarshalCustomDetails(oper.DataUpdateRequest.CreatorDetails),
		},
	}, nil
}

func (h *createDataUpdateRequestHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	participants := []history2.ParticipantEffect{h.Participant(sourceAccountID)}

	return participants, nil
}
