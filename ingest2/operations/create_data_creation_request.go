package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ingest2/internal"
)

type createDataCreationRequestHandler struct {
	effectsProvider
}

func (h *createDataCreationRequestHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	createDataRequest := op.Body.MustCreateDataCreationRequestOp()
	response := opRes.MustCreateDataCreationRequestResult().CreateDataCreationRequestResponse

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreateDataCreationRequest,
		CreateDataCreationRequest: &history2.CreateDataCreationRequest{
			ID:        uint64(response.RequestId),
			Type:      uint64(response.Type),
			RequestID: uint64(response.RequestId),
			Value:     internal.MarshalCustomDetails(response.Value),
			CreatorDetails: internal.MarshalCustomDetails(createDataRequest.DataCreationRequest.CreatorDetails),
			Owner:     createDataRequest.DataCreationRequest.Owner.Address(),
		},
	}, nil
}

func (h *createDataCreationRequestHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	participants := []history2.ParticipantEffect{h.Participant(sourceAccountID)}

	owner := opBody.MustCreateDataCreationRequestOp().DataCreationRequest.Owner
	if !sourceAccountID.Equals(owner) {
		participants = append(participants, h.Participant(owner))
	}

	return participants, nil
}
