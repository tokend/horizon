package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ingest2/internal"
	"gitlab.com/tokend/regources/v2"
)

type withdrawHandler struct {
	balanceProvider balanceProvider
}

//ParticipantsEffects - returns effect due to review of withdrawal request
func (h *withdrawHandler) ParticipantsEffects(op xdr.ReviewRequestOp,
	res xdr.ExtendedResult, request xdr.ReviewableRequestEntry,
	source history2.ParticipantEffect, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	details := request.Body.MustWithdrawalRequest()

	effect := history2.Effect{
		Type: history2.EffectTypeChargedFromLocked,
		ChargedFromLocked: &history2.BalanceChangeEffect{
			Amount: regources.Amount(details.Amount),
			Fee:    internal.FeeFromXdr(details.Fee),
		},
	}

	if op.Action != xdr.ReviewRequestOpActionApprove {
		effect = history2.Effect{
			Type: history2.EffectTypeWithdrawn,
			Withdrawn: &history2.BalanceChangeEffect{
				Amount: regources.Amount(details.Amount),
				Fee:    internal.FeeFromXdr(details.Fee),
			},
		}
	}

	return populateEffects(h.balanceProvider.MustBalance(details.Balance), effect, source), nil
}
