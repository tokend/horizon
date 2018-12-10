package operaitons

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources"
)

type createAtomicSwapBidRequestOpHandler struct {
	balanceProvider balanceProvider
}

// OperationDetails returns details about create atomic swap bid request operation
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

// ParticipantsEffects returns source effect with `locked` effect
func (h *createAtomicSwapBidRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	aSwapBidRequest := opBody.MustCreateASwapBidCreationRequestOp().Request

	balance := h.balanceProvider.GetBalanceByID(aSwapBidRequest.BaseBalance)

	source.BalanceID = &balance.ID
	source.AssetCode = &balance.AssetCode
	source.Effect = history2.Effect{
		Type: history2.EffectTypeLocked,
		Locked: &history2.LockedEffect{
			Amount: amount.StringU(uint64(aSwapBidRequest.Amount)),
		},
	}

	return []history2.ParticipantEffect{source}, nil
}
