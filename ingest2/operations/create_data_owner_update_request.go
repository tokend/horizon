package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ingest2/internal"
)

type createDataOwnerUpdateRequestHandler struct {
	effectsProvider
}

func (h *createDataOwnerUpdateRequestHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	oper := op.Body.MustCreateDataOwnerUpdateRequestOp()
	response := opRes.MustCreateDataOwnerUpdateRequestResult().Success

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreateDataOwnerUpdateRequest,
		CreateDataOwnerUpdateRequest: &history2.CreateDataOwnerUpdateRequest{
			ID:             uint64(oper.DataOwnerUpdateRequest.UpdateDataOwnerOp.DataId),
			NewOwner:       oper.DataOwnerUpdateRequest.UpdateDataOwnerOp.NewOwner,
			CreatorDetails: internal.MarshalCustomDetails(oper.DataOwnerUpdateRequest.CreatorDetails),
			RequestID:      uint64(response.RequestId),
		},
	}, nil
}

func (h *createDataOwnerUpdateRequestHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	participants := []history2.ParticipantEffect{
		h.Participant(sourceAccountID),
		h.Participant(opBody.UpdateDataOwnerOp.NewOwner),
	}

	return participants, nil
}
