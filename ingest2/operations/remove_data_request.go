package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type removeDataRequestOpHandler struct {
	effectsProvider
}

func (h *removeDataRequestOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	response := opRes.MustRemoveDataRequestResult().MustRemoveDataRequestResponse()

	return history2.OperationDetails{
		Type: xdr.OperationTypeRemoveDataRequest,
		RemoveDataRequest: &history2.RemoveDataRequestDetails{
			RequestID: uint64(response.RequestId),
		},
	}, nil
}

func (h *removeDataRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {

	return []history2.ParticipantEffect{h.Participant(sourceAccountID)}, nil
}
