package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ingest2/internal"
	"gitlab.com/tokend/regources/v2"
)

type issuanceHandler struct {
	balanceProvider balanceProvider
}

//ParticipantsEffects - returns effect for receiver of the funds and source of the op
func (h *issuanceHandler) ParticipantsEffects(op xdr.ReviewRequestOp,
	res xdr.ExtendedResult, request xdr.ReviewableRequestEntry,
	source history2.ParticipantEffect, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	if op.Action != xdr.ReviewRequestOpActionApprove {
		return []history2.ParticipantEffect{source}, nil
	}

	details := request.Body.MustIssuanceRequest()

	effect := history2.Effect{
		Type: history2.EffectTypeIssued,
		Issued: &history2.BalanceChangeEffect{
			Amount: regources.Amount(details.Amount),
			Fee:    internal.FeeFromXdr(details.Fee),
		},
	}

	return populateEffects(h.balanceProvider.MustBalance(details.Receiver), effect, source), nil
}
