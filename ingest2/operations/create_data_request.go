package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ingest2/internal"
)

type createDataRequestOpHandler struct {
	effectsProvider
}

func (h *createDataRequestOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	createDataRequest := op.Body.MustCreateDataRequestOp()
	response := opRes.MustCreateDataRequestResult().MustCreateDataRequestResponse()

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreateDataRequest,
		CreateDataRequest: &history2.CreateDataRequestDetails{
			ID:    uint64(response.RequestId),
			Type:  uint64(response.Type),
			Value: internal.MarshalCustomDetails(response.Value),
			Owner: createDataRequest.CreateDataRequest.Owner.Address(),
		},
	}, nil
}

func (h *createDataRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {

	return []history2.ParticipantEffect{h.Participant(sourceAccountID)}, nil
}
