package operations

import "gitlab.com/tokend/horizon/db2/history2"

type createDataHandler struct {
	effectsProvider
}

func (h *createDataHandler) PermanentReject(details requestDetails) ([]history2.ParticipantEffect, error) {

}

func (h *createDataHandler) Fulfilled(details requestDetails) ([]history2.ParticipantEffect, error) {

}
