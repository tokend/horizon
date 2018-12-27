package operations

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type deprecatedReviewRequestHandler struct {
}

//ParticipantsEffects - always returns errors, as deprecated request must not occur in the core
func (h *deprecatedReviewRequestHandler) ParticipantsEffects(op xdr.ReviewRequestOp, res xdr.ReviewRequestSuccessResult,
	request xdr.ReviewableRequestEntry, source history2.ParticipantEffect, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	return nil, errors.From(errors.New("tried to ingest deprecated reviewable request"), logan.F{
		"request_type": request.Body.Type,
	})
}
