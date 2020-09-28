package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	regources "gitlab.com/tokend/regources/generated"
)

type closeDeferredPaymentHandler struct {
	effectsProvider
	defPaymentProvider
}

func (h *closeDeferredPaymentHandler) Fulfilled(details requestDetails) ([]history2.ParticipantEffect, error) {
	closeDeferredPaymentRequest := details.Request.Body.MustCloseDeferredPaymentRequest()
	res := details.Result.TypeExt.MustCloseDeferredPaymentResult()
	funded := h.effectsProvider.BalanceEffect(res.DestinationBalance,
		&history2.Effect{
			Type: history2.EffectTypeFunded,
			Unlocked: &history2.BalanceChangeEffect{
				Amount: regources.Amount(closeDeferredPaymentRequest.Amount),
			},
		})

	dp := h.MustDeferredPayment(int64(closeDeferredPaymentRequest.DeferredPaymentId))
	var sb xdr.BalanceId
	sb.SetString(dp.SourceBalance)

	charged := h.effectsProvider.BalanceEffect(sb,
		&history2.Effect{
			Type: history2.EffectTypeChargedFromLocked,
			Unlocked: &history2.BalanceChangeEffect{
				Amount: regources.Amount(closeDeferredPaymentRequest.Amount),
			},
		})
	return []history2.ParticipantEffect{funded, charged}, nil
}

func (h *closeDeferredPaymentHandler) PermanentReject(details requestDetails) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{h.effectsProvider.Participant(details.SourceAccountID)}, nil
}
