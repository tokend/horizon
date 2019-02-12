package operations

import (
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/v2"
)

type amlAlertHandler struct {
	effectsProvider
}

//Fulfilled - returns participant of fully approved request
func (h *amlAlertHandler) Fulfilled(details requestDetails) ([]history2.ParticipantEffect, error) {
	amlRequest := details.Request.Body.MustAmlAlertRequest()
	effect := history2.Effect{
		Type: history2.EffectTypeWithdrawn,
		Withdrawn: &history2.BalanceChangeEffect{
			Amount: regources.Amount(amlRequest.Amount),
		},
	}
	return h.effectsProvider.BalanceEffectWithAccount(details.SourceAccountID, amlRequest.BalanceId, &effect), nil
}

//PermanentReject - returns participants of fully rejected request
func (h *amlAlertHandler) PermanentReject(details requestDetails) ([]history2.ParticipantEffect, error) {
	amlRequest := details.Request.Body.MustAmlAlertRequest()
	effect := history2.Effect{
		Type: history2.EffectTypeUnlocked,
		Unlocked: &history2.BalanceChangeEffect{
			Amount: regources.Amount(amlRequest.Amount),
		},
	}
	return h.effectsProvider.BalanceEffectWithAccount(details.SourceAccountID, amlRequest.BalanceId, &effect), nil
}
