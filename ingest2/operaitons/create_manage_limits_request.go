package operaitons

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type createManageLimitsRequestOpHandler struct {
}

func (h *createManageLimitsRequestOpHandler) OperationDetails(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	createManageLimitsRequestOp := opBody.MustCreateManageLimitsRequestOp()

	var data map[string]interface{}
	rawData, ok := createManageLimitsRequestOp.ManageLimitsRequest.Ext.GetDetails()
	if ok {
		data = customDetailsUnmarshal([]byte(rawData))
	}

	var requestID int64
	if rawID, ok := createManageLimitsRequestOp.Ext.GetRequestId(); ok {
		requestID = int64(rawID)
	}

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreateManageLimitsRequest,
		CreateManageLimitsRequest: &history2.CreateManageLimitsRequestDetails{
			Data:      data,
			RequestID: requestID,
		},
	}, nil
}

func (h *createManageLimitsRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{source}, nil
}
