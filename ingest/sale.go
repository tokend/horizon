package ingest

import (
	"encoding/json"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon/db2"
	"gitlab.com/swarmfund/horizon/db2/history"
	"time"
)

func saleCreate(is *Session, ledgerEntry *xdr.LedgerEntry) error {
	rawSale := ledgerEntry.Data.MustSale()
	sale, err := convertSale(rawSale)
	if err != nil {
		return errors.Wrap(err, "failed to convert sale")
	}

	err = is.Ingestion.HistoryQ().Sales().Insert(*sale)
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

	err = is.Ingestion.HistoryQ().Sales().Update(*sale)
	if err != nil {
		return errors.Wrap(err, "faied to update sale")
	}

	return nil
}

func convertSale(raw xdr.SaleEntry) (*history.Sale, error) {
	var quoteAssets []interface{}
	for i := range raw.QuoteAssets {
		quoteAssets = append(quoteAssets, map[string]interface{}{
			"asset":            raw.QuoteAssets[i].QuoteAsset,
			"price":            amount.StringU(uint64(raw.QuoteAssets[i].Price)),
			"quote_balance_id": raw.QuoteAssets[i].QuoteBalance.AsString(),
			"current_cap":      amount.StringU(uint64(raw.QuoteAssets[i].CurrentCap)),
		})
	}

	var saleDetails db2.Details
	_ = json.Unmarshal([]byte(raw.Details), &saleDetails)

	return &history.Sale{
		ID:                uint64(raw.SaleId),
		OwnerID:           raw.OwnerId.Address(),
		BaseAsset:         string(raw.BaseAsset),
		DefaultQuoteAsset: string(raw.DefaultQuoteAsset),
		StartTime:         time.Unix(int64(raw.StartTime), 0).UTC(),
		EndTime:           time.Unix(int64(raw.EndTime), 0).UTC(),
		SoftCap:           uint64(raw.SoftCap),
		HardCap:           uint64(raw.HardCap),
		Details:           saleDetails,
		QuoteAssets: map[string]interface{}{
			"quote_assets": quoteAssets,
		},
		State: history.SaleStateOpen,
	}, nil
}
