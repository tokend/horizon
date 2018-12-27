package operations

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type withdrawHandler struct {
	balanceProvider balanceProvider
}

//ParticipantsEffects - returns effect due to review of withdrawal request
func (h *withdrawHandler) ParticipantsEffects(op xdr.ReviewRequestOp,
	res xdr.ReviewRequestSuccessResult, request xdr.ReviewableRequestEntry,
	source history2.ParticipantEffect, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	details := request.Body.MustWithdrawalRequest()

	effect := history2.Effect{
		Type: history2.EffectTypeChargedFromLocked,
		ChargedFromLocked: &history2.BalanceChangeEffect{
			Amount: amount.StringU(uint64(details.Amount)),
			Fee: history2.Fee{
				Fixed:             amount.StringU(uint64(details.Fee.Fixed)),
				CalculatedPercent: amount.StringU(uint64(details.Fee.Percent)),
			},
		},
	}

	if op.Action != xdr.ReviewRequestOpActionApprove {
		effect = history2.Effect{
			Type: history2.EffectTypeWithdrawn,
			Withdrawn: &history2.BalanceChangeEffect{
				Amount: amount.StringU(uint64(details.Amount)),
				Fee: history2.Fee{
					Fixed:             amount.StringU(uint64(details.Fee.Fixed)),
					CalculatedPercent: amount.StringU(uint64(details.Fee.Percent)),
				},
			},
		}
	}

	return populateEffects(h.balanceProvider.MustBalance(details.Balance), effect, source), nil
}
