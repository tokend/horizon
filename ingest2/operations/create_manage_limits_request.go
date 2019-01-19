package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ingest2/internal"
)

type createManageLimitsRequestOpHandler struct {
}

// Details returns details about create limits request operation
func (h *createManageLimitsRequestOpHandler) Details(op rawOperation,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	createManageLimitsRequestOp := op.Body.MustCreateManageLimitsRequestOp()

	var data map[string]interface{}
	rawData, ok := createManageLimitsRequestOp.ManageLimitsRequest.Ext.GetDetails()
	if ok {
		data = internal.MarshalCustomDetails(rawData)
	}

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreateManageLimitsRequest,
		CreateManageLimitsRequest: &history2.CreateManageLimitsRequestDetails{
			Data:      data,
			RequestID: int64(opRes.MustCreateManageLimitsRequestResult().MustSuccess().ManageLimitsRequestId),
		},
	}, nil
}

//ParticipantsEffects returns source participant effect
func (h *createManageLimitsRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{source}, nil
}
