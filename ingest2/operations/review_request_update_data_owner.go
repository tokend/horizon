package operations

import "gitlab.com/tokend/horizon/db2/history2"

type updateDataOwnerHandler struct {
	effectsProvider
}

func (h *updateDataOwnerHandler) PermanentReject(details requestDetails) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{
		h.Participant(details.SourceAccountID),
		h.Participant(details.Request.Body.MustDataOwnerUpdateRequest().UpdateDataOwnerOp.NewOwner),
	}, nil
}

func (h *updateDataOwnerHandler) Fulfilled(details requestDetails) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{
		h.Participant(details.SourceAccountID),
		h.Participant(details.Request.Body.MustDataOwnerUpdateRequest().UpdateDataOwnerOp.NewOwner),
	}, nil
}
