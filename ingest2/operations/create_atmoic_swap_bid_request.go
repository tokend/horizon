package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ingest2/internal"
	"gitlab.com/tokend/regources/v2"
)

type createAtomicSwapBidRequestOpHandler struct {
	balanceProvider balanceProvider
}

// CreatorDetails returns details about create atomic swap bid request operation
func (h *createAtomicSwapBidRequestOpHandler) Details(op rawOperation,
	opRes xdr.OperationResultTr) (history2.OperationDetails, error) {

	aSwapBidRequest := op.Body.MustCreateASwapBidCreationRequestOp().Request

	quoteAssets := make([]regources.AssetPrice, 0, len(aSwapBidRequest.QuoteAssets))
	for _, quoteAsset := range aSwapBidRequest.QuoteAssets {
		quoteAssets = append(quoteAssets, regources.AssetPrice{
			Asset: string(quoteAsset.QuoteAsset),
			Price: regources.Amount(quoteAsset.Price),
		})
	}

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreateAswapBidRequest,
		CreateAtomicSwapBidRequest: &history2.CreateAtomicSwapBidRequestDetails{
			Amount:      regources.Amount(aSwapBidRequest.Amount),
			BaseBalance: aSwapBidRequest.BaseBalance.AsString(),
			QuoteAssets: quoteAssets,
			Details:     internal.MarshalCustomDetails(aSwapBidRequest.CreatorDetails),
		},
	}, nil
}

// ParticipantsEffects returns source effect with `locked` effect
func (h *createAtomicSwapBidRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	aSwapBidRequest := opBody.MustCreateASwapBidCreationRequestOp().Request

	balance := h.balanceProvider.MustBalance(aSwapBidRequest.BaseBalance)

	source.BalanceID = &balance.ID
	source.AssetCode = &balance.AssetCode
	source.Effect = &history2.Effect{
		Type: history2.EffectTypeLocked,
		Locked: &history2.BalanceChangeEffect{
			Amount: regources.Amount(aSwapBidRequest.Amount),
		},
	}

	return []history2.ParticipantEffect{source}, nil
}
