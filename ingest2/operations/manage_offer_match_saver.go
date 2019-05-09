package operations

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2/history2"
)

type manageOfferMatchSaver struct {
	storage matchesStorage
}

// Handle - Handles manage offer op by storing result matches
func (h *manageOfferMatchSaver) Handle(op operation) error {
	manageOfferOp := op.Operation().Body.MustManageOfferOp()
	manageOfferOpResult := op.Result().MustManageOfferResult().MustSuccess()

	for _, atom := range manageOfferOpResult.OffersClaimed {
		err := h.storage.Insert(history2.NewMatch(
			manageOfferOpResult.BaseAsset,
			manageOfferOpResult.QuoteAsset,
			manageOfferOp.OrderBookId,
			atom,
		))
		if err != nil {
			return errors.New("failed to insert match")
		}
	}

	return nil
}
