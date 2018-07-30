package ingest

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
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

func (is *Session) getOfferBaseAsset (changes xdr.LedgerEntryChanges, saleId xdr.Uint64) xdr.AssetCode {
	for _, change := range changes {
		if change.Type != xdr.LedgerEntryChangeTypeUpdated {
			continue
		}
		data := change.Updated.Data
		if data.Type == xdr.LedgerEntryTypeSale && data.Sale.SaleId == saleId {
			return data.Sale.BaseAsset
		}
	}
	return xdr.AssetCode("")
}
