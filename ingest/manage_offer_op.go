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

func (is *Session) deriveOfferBaseAsset (changes xdr.LedgerEntryChanges, saleId xdr.Uint64) xdr.AssetCode {
	for i := range changes {
		ledgerData := changes[i].Updated.Data
		if ledgerData.Type == xdr.LedgerEntryTypeSale && ledgerData.Sale.SaleId == saleId {
			return ledgerData.Sale.BaseAsset
		}
	}
	return xdr.AssetCode("")
}
