package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/rgenerated"
)

type createAtomicSwapRequestOpHandler struct {
	effectsProvider
}

// Details returns details about create atomic swap request operation
func (h *createAtomicSwapRequestOpHandler) Details(op rawOperation,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	aSwapRequest := op.Body.MustCreateASwapRequestOp().Request
	successRes := opRes.MustCreateASwapRequestResult().MustSuccess()

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreateAswapRequest,
		CreateAtomicSwapRequest: &history2.CreateAtomicSwapRequestDetails{
			BidID:      int64(aSwapRequest.BidId),
			BaseAmount: rgenerated.Amount(aSwapRequest.BaseAmount),
			QuoteAsset: string(aSwapRequest.QuoteAsset),
			RequestDetails: history2.RequestDetails{
				RequestID: int64(successRes.RequestId),
			},
		},
	}, nil
}
