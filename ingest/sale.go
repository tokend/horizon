package ingest

import (
	"encoding/json"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon/db2/history"
	"time"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon/db2"
)

func saleCreate(is *Session, ledgerEntry *xdr.LedgerEntry) error {
	rawSale := ledgerEntry.Data.MustSale()
	sale, err := convertSale(rawSale)
	if err != nil {
		return errors.Wrap(err, "failed to convert sale")
	}

	err = is.Ingestion.HistoryQ.Sales().Insert(*sale)
	if err != nil {
		return errors.Wrap(err, "failed to insert sale")
	}

	return nil
}

func saleUpdate(is *Session, ledgerEntry *xdr.LedgerEntry) error {
	rawSale := ledgerEntry.Data.MustSale()
	sale, err := convertSale(rawSale)
	if err != nil {
		return errors.Wrap(err, "failed to convert sale")
	}

	err = is.Ingestion.HistoryQ.Sales().Update(*sale)
	if err != nil {
		return errors.Wrap(err, "faied to update sale")
	}

	return nil
}

func convertSale(raw xdr.SaleEntry) (*history.Sale, error) {
	var saleDetails db2.Details
	_ = json.Unmarshal([]byte(raw.Details), &saleDetails)

	return &history.Sale{
		ID:         uint64(raw.SaleId),
		OwnerID:    raw.OwnerId.Address(),
		BaseAsset:  string(raw.BaseAsset),
		QuoteAsset: string(raw.QuoteAsset),
		StartTime:  time.Unix(int64(raw.StartTime), 0).UTC(),
		EndTime:    time.Unix(int64(raw.EndTime), 0).UTC(),
		Price:      uint64(raw.Price),
		SoftCap:    uint64(raw.SoftCap),
		HardCap:    uint64(raw.HardCap),
		CurrentCap: uint64(raw.CurrentCap),
		Details:    saleDetails,
		State:      history.SaleStateOpen,
	}, nil
}
