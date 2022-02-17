package operations

import (
	"gitlab.com/tokend/horizon/db2/history2"

	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"
)

type lpRemoveLiquidityHandler struct {
	effectsProvider
}

func (h *lpRemoveLiquidityHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	lpRemoveRes := opRes.MustLpRemoveLiquidityResult().MustSuccess()
	lpRemoveOp := op.Body.MustLpRemoveLiquidityOp()

	return history2.OperationDetails{
		Type: xdr.OperationTypeLpRemoveLiquidity,
		LiquidityPoolRemoveLiquidity: &history2.LiquidityPoolManageLiquidity{
			LiquidityPoolID:   uint64(lpRemoveRes.LiquidityPoolId),
			FirstBalance:      lpRemoveRes.SourceFirstAssetBalanceId.AsString(),
			SecondBalance:     lpRemoveRes.SourceSecondAssetBalanceId.AsString(),
			FirstAssetAmount:  regources.Amount(lpRemoveRes.FirstAssetAmount),
			SecondAssetAmount: regources.Amount(lpRemoveRes.SecondAssetAmount),
			LPTokensAmount:    regources.Amount(lpRemoveOp.LpTokensAmount),
		},
	}, nil
}

func (h *lpRemoveLiquidityHandler) ParticipantsEffects(opBody xdr.OperationBody, opRes xdr.OperationResultTr,
	sourceID xdr.AccountId, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	return h.participantEffects(opBody.MustLpRemoveLiquidityOp(), opRes.MustLpRemoveLiquidityResult().MustSuccess())
}

func (h *lpRemoveLiquidityHandler) participantEffects(op xdr.LpRemoveLiquidityOp, res xdr.LpRemoveLiquiditySuccess,
) ([]history2.ParticipantEffect, error) {
	lpFirstBalanceEffect := h.BalanceEffect(res.LpFirstAssetBalanceId, &history2.Effect{
		Type: history2.EffectTypeCharged,
		Charged: &history2.BalanceChangeEffect{
			Amount: regources.Amount(res.FirstAssetAmount),
			Fee:    regources.Fee{},
		},
	})

	lpSecondBalanceEffect := h.BalanceEffect(res.LpSecondAssetBalanceId, &history2.Effect{
		Type: history2.EffectTypeCharged,
		Charged: &history2.BalanceChangeEffect{
			Amount: regources.Amount(res.SecondAssetAmount),
			Fee:    regources.Fee{},
		},
	})

	sourceFirstBalanceEffect := h.BalanceEffect(res.SourceFirstAssetBalanceId, &history2.Effect{
		Type: history2.EffectTypeFunded,
		Funded: &history2.BalanceChangeEffect{
			Amount: regources.Amount(res.FirstAssetAmount),
			Fee:    regources.Fee{},
		},
	})

	sourceSecondBalanceEffect := h.BalanceEffect(res.SourceSecondAssetBalanceId, &history2.Effect{
		Type: history2.EffectTypeFunded,
		Funded: &history2.BalanceChangeEffect{
			Amount: regources.Amount(res.SecondAssetAmount),
			Fee:    regources.Fee{},
		},
	})

	lpTokensBalanceEffect := h.BalanceEffect(op.LpTokenBalance, &history2.Effect{
		Type: history2.EffectTypeCharged,
		Charged: &history2.BalanceChangeEffect{
			Amount: regources.Amount(op.LpTokensAmount),
			Fee:    regources.Fee{},
		},
	})

	return []history2.ParticipantEffect{
		lpFirstBalanceEffect,
		lpSecondBalanceEffect,
		sourceFirstBalanceEffect,
		sourceSecondBalanceEffect,
		lpTokensBalanceEffect,
	}, nil
}
