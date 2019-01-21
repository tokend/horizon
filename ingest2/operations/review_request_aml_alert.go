package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/v2"
)

type amlAlertHandler struct {
	balanceProvider balanceProvider
}

//ParticipantsEffects - returns source participant and effects for balance for which AML Alert was created
func (h *amlAlertHandler) ParticipantsEffects(op xdr.ReviewRequestOp,
	res xdr.ExtendedResult, request xdr.ReviewableRequestEntry,
	source history2.ParticipantEffect, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	details := request.Body.MustAmlAlertRequest()

	effect := regources.Effect{
		Type: regources.EffectTypeWithdrawn,
		Withdrawn: &regources.BalanceChangeEffect{
			Amount: regources.Amount(details.Amount),
		},
	}

	if op.Action != xdr.ReviewRequestOpActionApprove {
		effect = regources.Effect{
			Type: regources.EffectTypeUnlocked,
			Unlocked: &regources.BalanceChangeEffect{
				Amount: regources.Amount(details.Amount),
			},
		}
	}

	return populateEffects(h.balanceProvider.MustBalance(details.BalanceId), effect, source), nil
}
