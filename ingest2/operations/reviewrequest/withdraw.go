package reviewrequest

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type withdrawHandler struct {
	effectHelper effectHelper
}

func (h *withdrawHandler) specificParticipantsEffects(op xdr.ReviewRequestOp,
	res xdr.ReviewRequestSuccessResult, request xdr.ReviewableRequestEntry,
	source history2.ParticipantEffect, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	details := request.Body.MustWithdrawalRequest()

	effect := history2.Effect{
		Type: history2.EffectTypeChargedFromLocked,
		ChargedFromLocked: &history2.ChargedFromLockedEffect{
			Amount: amount.StringU(uint64(details.Amount)),
			FeePaid: history2.FeePaid{
				Fixed:             amount.StringU(uint64(details.Fee.Fixed)),
				CalculatedPercent: amount.StringU(uint64(details.Fee.Percent)),
			},
		},
	}

	if op.Action != xdr.ReviewRequestOpActionApprove {
		effect = history2.Effect{
			Type: history2.EffectTypeWithdrawn,
			Withdrawn: &history2.ChargedFromLockedEffect{
				Amount: amount.StringU(uint64(details.Amount)),
				FeePaid: history2.FeePaid{
					Fixed:             amount.StringU(uint64(details.Fee.Fixed)),
					CalculatedPercent: amount.StringU(uint64(details.Fee.Percent)),
				},
			},
		}
	}

	return h.effectHelper.populateEffectWithBalanceDetails(details.Balance, effect, source), nil
}
