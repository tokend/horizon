package operations

import (
	"gitlab.com/tokend/horizon/db2/history2"

	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"
)

type lpSwapOpHandler struct {
	effectsProvider
}

func (h *lpSwapOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	lpSwapRes := opRes.MustLpSwapResult().MustSuccess()
	lpSwapOp := op.Body.MustLpSwapOp()

	return history2.OperationDetails{
		Type: xdr.OperationTypeLpSwap,
		LiquidityPoolSwap: &history2.LiquidityPoolSwap{
			LiquidityPoolID:  uint64(lpSwapRes.LiquidityPoolId),
			SourceInBalance:  lpSwapRes.SourceInBalanceId.AsString(),
			SourceOutBalance: lpSwapRes.SourceOutBalanceId.AsString(),
			InAmount:         regources.Amount(lpSwapRes.SwapInAmount),
			OutAmount:        regources.Amount(lpSwapRes.SwapOutAmount),
			SwapType:         lpSwapOp.LpSwapRequest.Type.String(),
		},
	}, nil
}

func (h *lpSwapOpHandler) ParticipantsEffects(opBody xdr.OperationBody, opRes xdr.OperationResultTr,
	sourceID xdr.AccountId, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	return h.participantEffects(opBody.MustLpSwapOp(), opRes.MustLpSwapResult().MustSuccess())
}

func (h *lpSwapOpHandler) participantEffects(op xdr.LpSwapOp, res xdr.LpSwapSuccess,
) ([]history2.ParticipantEffect, error) {
	lpInBalanceEffect := h.BalanceEffect(res.LpInBalanceId, &history2.Effect{
		Type: history2.EffectTypeFunded,
		Funded: &history2.BalanceChangeEffect{
			Amount: regources.Amount(res.SwapInAmount),
			Fee:    regources.Fee{},
		},
	})

	lpOutBalanceEffect := h.BalanceEffect(res.LpOutBalanceId, &history2.Effect{
		Type: history2.EffectTypeCharged,
		Charged: &history2.BalanceChangeEffect{
			Amount: regources.Amount(res.SwapOutAmount),
			Fee:    regources.Fee{},
		},
	})

	sourceInBalanceEffect := h.BalanceEffect(op.FromBalance, &history2.Effect{
		Type: history2.EffectTypeCharged,
		Charged: &history2.BalanceChangeEffect{
			Amount: regources.Amount(res.SwapInAmount),
			Fee:    regources.Fee{},
		},
	})

	sourceOutBalanceEffect := h.BalanceEffect(op.ToBalance, &history2.Effect{
		Type: history2.EffectTypeFunded,
		Funded: &history2.BalanceChangeEffect{
			Amount: regources.Amount(res.SwapOutAmount),
			Fee:    regources.Fee{},
		},
	})

	return []history2.ParticipantEffect{
		lpInBalanceEffect,
		lpOutBalanceEffect,
		sourceInBalanceEffect,
		sourceOutBalanceEffect,
	}, nil
}
