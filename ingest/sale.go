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
	var quoteAssets []history.QuoteAsset
	for i := range raw.QuoteAssets {
		quoteAssets = append(quoteAssets, history.QuoteAsset{
			Asset:          string(raw.QuoteAssets[i].QuoteAsset),
			Price:          amount.StringU(uint64(raw.QuoteAssets[i].Price)),
			QuoteBalanceID: raw.QuoteAssets[i].QuoteBalance.AsString(),
			CurrentCap:     amount.StringU(uint64(raw.QuoteAssets[i].CurrentCap)),
		})
	}

	var saleDetails db2.Details
	_ = json.Unmarshal([]byte(raw.Details), &saleDetails)

	saleType := xdr.SaleTypeBasicSale
	if raw.Ext.SaleTypeExt != nil {
		saleType = raw.Ext.SaleTypeExt.TypedSale.SaleType
	}

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
		QuoteAssets: history.QuoteAssets{
			QuoteAssets: quoteAssets,
		},
		BaseCurrentCap: int64(raw.CurrentCapInBase),
		BaseHardCap:    int64(raw.MaxAmountToBeSold),
		State:          history.SaleStateOpen,
		SaleType:       saleType,
	}, nil
}
