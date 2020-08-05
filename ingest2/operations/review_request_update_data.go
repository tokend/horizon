package operations

import "gitlab.com/tokend/horizon/db2/history2"

type updateDataHandler struct {
	*manageUpdateDataOpHandler
}

func (h *updateDataHandler) PermanentReject(details requestDetails) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{h.Participant(details.SourceAccountID)}, nil
}

func (h *updateDataHandler) Fulfilled(details requestDetails) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{h.Participant(details.SourceAccountID)}, nil
}
