package operaitons

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type createSaleRequestOpHandler struct {
}

func (h *createSaleRequestOpHandler) OperationDetails(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	createSaleRequest := op.Body.MustCreateSaleCreationRequestOp().Request

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreateSaleRequest,
		CreateSaleRequest: &history2.CreateSaleRequestDetails{
			RequestID:         int64(opRes.MustCreateSaleCreationRequestResult().MustSuccess().RequestId),
			BaseAsset:         createSaleRequest.BaseAsset,
			DefaultQuoteAsset: createSaleRequest.DefaultQuoteAsset,
			StartTime:         int64(createSaleRequest.StartTime),
			EndTime:           int64(createSaleRequest.EndTime),
			HardCap:           amount.StringU(uint64(createSaleRequest.HardCap)),
			SoftCap:           amount.StringU(uint64(createSaleRequest.SoftCap)),
			Details:           customDetailsUnmarshal([]byte(createSaleRequest.Details)),
		},
	}, nil
}

func (h *createSaleRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{source}, nil
}
