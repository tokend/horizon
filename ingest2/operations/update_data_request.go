package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ingest2/internal"
)

type updateDataRequestOpHandler struct {
	effectsProvider
}

func (h *updateDataRequestOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	response := opRes.MustUpdateDataRequestResult().MustUpdateDataRequestResponse()

	return history2.OperationDetails{
		Type: xdr.OperationTypeUpdateDataRequest,
		UpdateDataRequest: &history2.UpdateDataRequestDetails{
			DataID:    uint64(response.DataId),
			RequestID: uint64(response.RequestId),
			Type:      uint64(response.Type),
			Value:     internal.MarshalCustomDetails(response.Value),
			Owner:     response.Owner.Address(),
		},
	}, nil
}

func (h *updateDataRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {

	return []history2.ParticipantEffect{h.Participant(sourceAccountID)}, nil
}
