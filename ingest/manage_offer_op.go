package ingest

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (is *Session) updateOfferState(offerID, state uint64) {
	if is.Err != nil {
		return
	}

	err := is.Ingestion.UpdateOfferState(offerID, state)
	if err != nil {
		is.Err = errors.Wrap(err, "failed to update offer state")
		return
	}
}
