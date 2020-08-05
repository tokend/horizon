package operations

import "gitlab.com/tokend/horizon/db2/history2"

type createDataHandler struct {
	*manageCreateDataOpHandler
}

func (h *createDataHandler) PermanentReject(details requestDetails) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{h.Participant(details.SourceAccountID)}, nil
}

func (h *createDataHandler) Fulfilled(details requestDetails) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{h.Participant(details.SourceAccountID)}, nil
}
