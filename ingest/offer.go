package ingest

import (
	"time"

	"gitlab.com/distributed_lab/logan/v3"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history"
)

func offerCreate(is *Session, entry *xdr.LedgerEntry) error {
	op := is.Cursor.Operation().Body.MustManageOfferOp()
	offer := convertOffer(entry.Data.MustOffer(), int64(op.Amount))

	err := is.Ingestion.HistoryQ().Offers().Insert(offer)
	if err != nil {
		return errors.Wrap(err, "failed to insert offer")
	}

	return nil
}

func offerUpdate(is *Session, entry *xdr.LedgerEntry) error {
	offer := entry.Data.MustOffer()

	err := is.Ingestion.HistoryQ().Offers().UpdateBaseAmount(
		int64(offer.BaseAmount),
		int64(offer.OfferId),
	)
	if err != nil {
		return errors.Wrap(err, "failed to update offer current base amount")
	}

	return nil
}

func offerDelete(is *Session, key *xdr.LedgerKey) error {
	op := is.Cursor.Operation().Body.MustManageOfferOp()
	if op.Amount == 0 {
		err := is.Ingestion.HistoryQ().Offers().Cancel(int64(op.OfferId))
		if err != nil {
			return errors.Wrap(err, "failed to cancel offer", logan.F{
				"offer_id": op.OfferId,
			})
		}

		return nil
	}

	err := is.Ingestion.HistoryQ().Offers().UpdateBaseAmount(
		0, int64(key.MustOffer().OfferId))
	if err != nil {
		return errors.Wrap(err, "failed to set offer zero current base amount")
	}

	return nil
}

func convertOffer(raw xdr.OfferEntry, initialBaseAmount int64) history.Offer {
	return history.Offer{
		OfferID:           int64(raw.OfferId),
		BaseAsset:         string(raw.Base),
		QuoteAsset:        string(raw.Quote),
		InitialBaseAmount: initialBaseAmount,
		CurrentBaseAmount: int64(raw.BaseAmount),
		Price:             int64(raw.Price),
		OwnerID:           raw.OwnerId.Address(),
		IsCanceled:        false,
		CreatedAt:         time.Unix(int64(raw.CreatedAt), 0).UTC(),
	}
}

func (is *Session) processCreateMatchedOffer(op xdr.ManageOfferOp, res xdr.ManageOfferSuccessResult) error {
	if (res.Offer.Effect != xdr.ManageOfferEffectDeleted) ||
		(op.Amount == 0) {
		return nil
	}

	newOffer := history.Offer{
		OfferID:           0,
		BaseAsset:         string(res.BaseAsset),
		QuoteAsset:        string(res.QuoteAsset),
		InitialBaseAmount: int64(op.Amount),
		CurrentBaseAmount: 0,
		Price:             int64(op.Price),
		OwnerID:           is.Cursor.OperationSourceAccount().Address(),
		IsCanceled:        false,
		CreatedAt:         is.Cursor.LedgerCloseTime(),
	}

	err := is.Ingestion.HistoryQ().Offers().Insert(newOffer)
	if err != nil {
		return errors.Wrap(err, "failed to insert offer")
	}

	return nil
}
