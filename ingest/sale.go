package ingest

import (
	"encoding/json"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon/db2"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"time"
)

func saleCreate(is *Session, ledgerEntry *xdr.LedgerEntry) error {
	rawSale := ledgerEntry.Data.MustSale()
	sale, err := convertSale(rawSale)
	if err != nil {
		return errors.Wrap(err, "failed to convert sale")
	}

	histSale, err := is.Ingestion.HistoryQ().Sales().ByID(sale.ID)
	if err != nil {
		return errors.Wrap(err, "failed to get sale from History DB")
	}

	if histSale != nil {
		err = is.Ingestion.HistoryQ().Sales().Update(*sale)
		if err != nil {
			return errors.Wrap(err, "failed to update sale")
		}
		return nil
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
		return errors.Wrap(err, "failed to update sale")
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

	rawState := xdr.SaleStateNone
	if raw.Ext.StatableSaleExt != nil {
		rawState = raw.Ext.StatableSaleExt.State
	}

	state, err := convertSaleState(rawState)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert sale state")
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
		State:          state,
		SaleType:       saleType,
	}, nil
}

func convertSaleState(state xdr.SaleState) (history.SaleState, error) {
	switch state {
	case xdr.SaleStateNone:
		return history.SaleStateOpen, nil
	case xdr.SaleStatePromotion:
		return history.SaleStatePromotion, nil
	case xdr.SaleStateVoting:
		return history.SaleStateVoting, nil
	default:
		return history.SaleStateOpen, errors.From(errors.New("unepxected sale of the state"), map[string]interface{}{
			"state": state,
		})
	}
}
