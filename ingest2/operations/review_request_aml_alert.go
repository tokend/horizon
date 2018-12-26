package operations

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type amlAlertHandler struct {
	balanceProvider balanceProvider
}

func (h *amlAlertHandler) ParticipantsEffects(op xdr.ReviewRequestOp,
	res xdr.ReviewRequestSuccessResult, request xdr.ReviewableRequestEntry,
	source history2.ParticipantEffect, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	details := request.Body.MustAmlAlertRequest()

	effect := history2.Effect{
		Type: history2.EffectTypeWithdrawn,
		Withdrawn: &history2.BalanceChangeEffect{
			Amount: amount.StringU(uint64(details.Amount)),
		},
	}

	if op.Action != xdr.ReviewRequestOpActionApprove {
		effect = history2.Effect{
			Type: history2.EffectTypeUnlocked,
			Unlocked: &history2.BalanceChangeEffect{
				Amount: amount.StringU(uint64(details.Amount)),
			},
		}
	}

	return populateEffects(h.balanceProvider.MustBalance(details.BalanceId), effect, source), nil
}
