package operaitons

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type createAtomicSwapRequestOpHandler struct {
}

func (h *createAtomicSwapRequestOpHandler) OperationDetails(op rawOperation,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	aSwapRequest := op.Body.MustCreateASwapRequestOp().Request
	successRes := opRes.MustCreateASwapRequestResult().MustSuccess()

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreateAswapRequest,
		CreateAtomicSwapRequest: &history2.CreateAtomicSwapRequestDetails{
			BidID:      int64(aSwapRequest.BidId),
			BaseAmount: amount.StringU(uint64(aSwapRequest.BaseAmount)),
			QuoteAsset: aSwapRequest.QuoteAsset,
			RequestDetails: history2.RequestDetails{
				RequestID: int64(successRes.RequestId),
			},
		},
	}, nil
}

func (h *createAtomicSwapRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{source}, nil
}
