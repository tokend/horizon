package operations

import (
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ingest2/internal"
	regources "gitlab.com/tokend/regources/generated"
)

type issuanceHandler struct {
	effectsProvider
}

//Fulfilled - returns slice of effects for participants of the operation
func (h *issuanceHandler) Fulfilled(details requestDetails) ([]history2.ParticipantEffect, error) {
	request := details.Request.Body.MustIssuanceRequest()

	effect := history2.Effect{
		Type: history2.EffectTypeIssued,
		Issued: &history2.BalanceChangeEffect{
			Amount: regources.Amount(request.Amount),
			Fee:    internal.FeeFromXdr(request.Fee),
		},
	}

	return h.BalanceEffectWithAccount(details.SourceAccountID, request.Receiver, &effect), nil
}

//PermanentReject - returns participants of fully rejected request
func (h *issuanceHandler) PermanentReject(details requestDetails) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{h.Participant(details.SourceAccountID)}, nil
}
