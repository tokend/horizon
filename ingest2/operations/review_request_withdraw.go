package operations

import (
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ingest2/internal"
	"gitlab.com/tokend/regources/v2"
)

type withdrawHandler struct {
	effectsProvider
}

//Fulfilled - returns slice of effects for participants of the operation
func (h *withdrawHandler) Fulfilled(details requestDetails) ([]history2.ParticipantEffect, error) {
	request := details.Request.Body.MustWithdrawalRequest()

	effect := history2.Effect{
		Type: history2.EffectTypeWithdrawn,
		Withdrawn: &history2.BalanceChangeEffect{
			Amount: regources.Amount(request.Amount),
			Fee:    internal.FeeFromXdr(request.Fee),
		},
	}

	return h.BalanceEffectWithAccount(details.SourceAccountID, request.Balance, &effect), nil
}

//PermanentReject - returns participants of fully rejected request
func (h *withdrawHandler) PermanentReject(details requestDetails) ([]history2.ParticipantEffect, error) {
	request := details.Request.Body.MustWithdrawalRequest()

	effect := history2.Effect{
		Type: history2.EffectTypeUnlocked,
		Unlocked: &history2.BalanceChangeEffect{
			Amount: regources.Amount(request.Amount),
			Fee:    internal.FeeFromXdr(request.Fee),
		},
	}

	return h.effectsProvider.BalanceEffectWithAccount(details.SourceAccountID, request.Balance, &effect), nil
}
