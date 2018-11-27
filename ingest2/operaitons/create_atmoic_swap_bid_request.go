package operaitons

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources"
)

type createAtomicSwapBidRequestOpHandler struct {
}

func (h *createAtomicSwapBidRequestOpHandler) OperationDetails(op rawOperation, opRes xdr.OperationResultTr) (history2.OperationDetails, error) {
	aSwapBidRequest := op.Body.MustCreateASwapBidCreationRequestOp().Request

	var quoteAssets []regources.SaleQuoteAsset
	for _, quoteAsset := range aSwapBidRequest.QuoteAssets {
		quoteAssets = append(quoteAssets, regources.SaleQuoteAsset{
			QuoteAsset: string(quoteAsset.QuoteAsset),
			Price:      regources.Amount(quoteAsset.Price),
		})
	}

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreateAswapBidRequest,
		CreateAtomicSwapBidRequest: &history2.CreateAtomicSwapBidRequestDetails{
			Amount:      amount.StringU(uint64(aSwapBidRequest.Amount)),
			BaseBalance: aSwapBidRequest.BaseBalance.AsString(),
			QuoteAssets: quoteAssets,
			Details:     customDetailsUnmarshal([]byte(aSwapBidRequest.Details)),
		},
	}, nil
}

func (h *createAtomicSwapBidRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	aSwapBidRequest := opBody.MustCreateASwapBidCreationRequestOp().Request

	return nil, nil
}
