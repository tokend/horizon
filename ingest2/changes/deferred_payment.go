package changes

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	history "gitlab.com/tokend/horizon/db2/history2"
	regources "gitlab.com/tokend/regources/generated"
)

//go:generate mockery -case underscore -name deferredPaymentPairStorage -inpkg -testonly
type deferredPaymentStorage interface {
	Insert(deferredPayment history.DeferredPayment) error
	Update(deferredPayment history.DeferredPayment) error
	Remove(deferredPaymentID int64) error
}

type deferredPaymentHandler struct {
	storage deferredPaymentStorage
}

func newDeferredPaymentHandler(storage deferredPaymentStorage) *deferredPaymentHandler {
	return &deferredPaymentHandler{
		storage: storage,
	}
}

func (h *deferredPaymentHandler) Created(lc ledgerChange) error {
	rawDeferredPayment := lc.LedgerChange.MustCreated().Data.MustDeferredPayment()
	deferredPayment := h.convertDeferredPayment(rawDeferredPayment)
	if err := h.storage.Insert(deferredPayment); err != nil {
		return errors.Wrap(err, "failed to insert from created")
	}
	return nil
}

func (h *deferredPaymentHandler) Updated(lc ledgerChange) error {
	rawDeferredPayment := lc.LedgerChange.MustUpdated().Data.MustDeferredPayment()
	deferredPayment := h.convertDeferredPayment(rawDeferredPayment)
	if err := h.storage.Update(deferredPayment); err != nil {
		return errors.Wrap(err, "failed to update from updated")
	}
	return nil
}

func (h *deferredPaymentHandler) Removed(lc ledgerChange) error {
	id := lc.LedgerChange.MustRemoved().MustDeferredPayment().Id
	if err := h.storage.Remove(int64(id)); err != nil {
		return errors.Wrap(err, "failed to remove deferredPayment by id")
	}

	return nil
}

func (h *deferredPaymentHandler) convertDeferredPayment(raw xdr.DeferredPaymentEntry) history.DeferredPayment {
	return history.DeferredPayment{
		ID:                    int64(raw.Id),
		Amount:                regources.Amount(raw.Amount),
		SourceAccount:         raw.Source.Address(),
		SourceBalance:         raw.SourceBalance.AsString(),
		DestinationAccount:    raw.Destination.Address(),
		SourcePaysForDest:     raw.FeeData.SourcePaysForDest,
		SourceFixedFee:        regources.Amount(raw.FeeData.SourceFee.Fixed),
		SourcePercentFee:      regources.Amount(raw.FeeData.SourceFee.Percent),
		DestinationFixedFee:   regources.Amount(raw.FeeData.DestinationFee.Fixed),
		DestinationPercentFee: regources.Amount(raw.FeeData.DestinationFee.Percent),
	}
}
