package operations

import (
	"gitlab.com/tokend/horizon/db2/history2"

	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"
)

type lpAddLiquidityOpHandler struct {
	effectsProvider
}

func (h *lpAddLiquidityOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	addLiquidityOpRes := opRes.MustLpAddLiquidityResult().MustSuccess()

	return history2.OperationDetails{
		Type: xdr.OperationTypeLpAddLiquidity,
		LiquidityPoolAddLiquidity: &history2.LiquidityPoolManageLiquidity{
			LiquidityPoolID:   uint64(addLiquidityOpRes.LiquidityPoolId),
			FirstBalance:      addLiquidityOpRes.SourceFirstAssetBalanceId.AsString(),
			SecondBalance:     addLiquidityOpRes.SourceSecondAssetBalanceId.AsString(),
			FirstAssetAmount:  regources.Amount(addLiquidityOpRes.FirstAssetAmount),
			SecondAssetAmount: regources.Amount(addLiquidityOpRes.SecondAssetAmount),
			LPTokensAmount:    regources.Amount(addLiquidityOpRes.LpTokensAmount),
		},
	}, nil
}

func (h *lpAddLiquidityOpHandler) ParticipantsEffects(opBody xdr.OperationBody, opRes xdr.OperationResultTr,
	sourceID xdr.AccountId, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	return h.participantEffects(opBody.MustLpAddLiquidityOp(), opRes.MustLpAddLiquidityResult().MustSuccess())
}

func (h *lpAddLiquidityOpHandler) participantEffects(op xdr.LpAddLiquidityOp, res xdr.LpAddLiquiditySuccess,
) ([]history2.ParticipantEffect, error) {
	lpFirstBalanceEffect := h.BalanceEffect(res.LpFirstAssetBalanceId, &history2.Effect{
		Type: history2.EffectTypeFunded,
		Funded: &history2.BalanceChangeEffect{
			Amount: regources.Amount(res.FirstAssetAmount),
			Fee:    regources.Fee{},
		},
	})

	lpSecondBalanceEffect := h.BalanceEffect(res.LpSecondAssetBalanceId, &history2.Effect{
		Type: history2.EffectTypeFunded,
		Funded: &history2.BalanceChangeEffect{
			Amount: regources.Amount(res.SecondAssetAmount),
			Fee:    regources.Fee{},
		},
	})

	sourceFirstBalanceEffect := h.BalanceEffect(res.SourceFirstAssetBalanceId, &history2.Effect{
		Type: history2.EffectTypeCharged,
		Charged: &history2.BalanceChangeEffect{
			Amount: regources.Amount(res.FirstAssetAmount),
			Fee:    regources.Fee{},
		},
	})

	sourceSecondBalanceEffect := h.BalanceEffect(res.SourceSecondAssetBalanceId, &history2.Effect{
		Type: history2.EffectTypeCharged,
		Charged: &history2.BalanceChangeEffect{
			Amount: regources.Amount(res.SecondAssetAmount),
			Fee:    regources.Fee{},
		},
	})

	lpTokensBalanceEffect := h.BalanceEffect(res.LpTokensBalanceId, &history2.Effect{
		Type: history2.EffectTypeIssued,
		Issued: &history2.BalanceChangeEffect{
			Amount: regources.Amount(res.LpTokensAmount),
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
