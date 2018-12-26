package operations

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type issuanceHandler struct {
	balanceProvider balanceProvider
}

func (h *issuanceHandler) ParticipantsEffects(op xdr.ReviewRequestOp,
	res xdr.ReviewRequestSuccessResult, request xdr.ReviewableRequestEntry,
	source history2.ParticipantEffect, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	if op.Action != xdr.ReviewRequestOpActionApprove {
		return []history2.ParticipantEffect{source}, nil
	}

	details := request.Body.MustIssuanceRequest()

	effect := history2.Effect{
		Type: history2.EffectTypeIssued,
		Issued: &history2.BalanceChangeEffect{
			Amount: amount.StringU(uint64(details.Amount)),
			Fee: history2.Fee{
				Fixed:             amount.StringU(uint64(details.Fee.Fixed)),
				CalculatedPercent: amount.StringU(uint64(details.Fee.Percent)),
			},
		},
	}

	return populateEffects(h.balanceProvider.MustBalance(details.Receiver), effect, source), nil
}
