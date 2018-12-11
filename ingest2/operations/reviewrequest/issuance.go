package reviewrequest

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type issuanceHandler struct {
	effectHelper effectHelper
}

func (h *issuanceHandler) specificParticipantsEffects(op xdr.ReviewRequestOp,
	res xdr.ReviewRequestSuccessResult, request xdr.ReviewableRequestEntry,
	source history2.ParticipantEffect, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	if op.Action != xdr.ReviewRequestOpActionApprove {
		return []history2.ParticipantEffect{source}, nil
	}

	details := request.Body.MustIssuanceRequest()

	effect := history2.Effect{
		Type: history2.EffectTypeIssued,
		Issued: &history2.FundedEffect{
			Amount: amount.StringU(uint64(details.Amount)),
			FeePaid: history2.FeePaid{
				Fixed:             amount.StringU(uint64(details.Fee.Fixed)),
				CalculatedPercent: amount.StringU(uint64(details.Fee.Percent)),
			},
		},
	}

	return h.effectHelper.getParticipantEffectByBalanceID(details.Receiver, effect, source), nil
}
