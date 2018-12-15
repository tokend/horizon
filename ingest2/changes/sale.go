package changes

import (
	"encoding/json"
	"fmt"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2"
	"gitlab.com/tokend/horizon/db2/history"
)

type saleStorage interface {
	//Inserts sale into DB
	InsertSale(sale history.Sale) error
	//Updates sale
	UpdateSale(sale history.Sale) error
}

type saleHandler struct {
	storage saleStorage
}

func newSaleHandler(storage saleStorage) *saleHandler {
	return &saleHandler{
		storage: storage,
	}
}

func (c *saleHandler) Created(lc ledgerChange) error {
	rawSale := lc.LedgerChange.MustCreated().Data.MustSale()

	sale, err := c.convertSale(rawSale)
	if err != nil {
		return errors.Wrap(err, "failed to convert sale", logan.F{
			"raw_sale":        rawSale,
			"ledger_sequence": lc.LedgerSeq,
		})
	}

	err = c.storage.InsertSale(*sale)
	if err != nil {
		return errors.Wrap(err, "failed to insert sale into DB", logan.F{
			"sale": sale,
		})
	}

	return nil
}

func (c *saleHandler) Updated(lc ledgerChange) error {
	rawSale := lc.LedgerChange.MustUpdated().Data.MustSale()
	sale, err := c.convertSale(rawSale)
	if err != nil {
		return errors.Wrap(err, "failed to convert sale ", logan.F{
			"raw_sale":        rawSale,
			"ledger_sequence": lc.LedgerSeq,
		})
	}

	err = c.storage.UpdateSale(*sale)
	if err != nil {
		return errors.Wrap(err, "failed to update sale", logan.F{
			"sale": sale,
		})
	}
	return nil
}

func (c *saleHandler) convertSale(raw xdr.SaleEntry) (*history.Sale, error) {
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
	err := json.Unmarshal([]byte(raw.Details), &saleDetails)
	if err != nil {
		saleDetails["reason"] = fmt.Sprintf("Expected json, got %v", raw.Details)
		saleDetails["error"] = err.Error()
	}

	saleType := xdr.SaleTypeBasicSale
	rawState := xdr.SaleStateNone
	switch raw.Ext.V {
	case xdr.LedgerVersionEmptyVersion:
	case xdr.LedgerVersionTypedSale:
		saleType = raw.Ext.MustSaleTypeExt().TypedSale.SaleType
	case xdr.LedgerVersionStatableSales:
		ext := raw.Ext.MustStatableSaleExt()
		saleType = ext.SaleTypeExt.TypedSale.SaleType
		rawState = ext.State
	default:
		panic(errors.Wrap(errors.New("Unexpected ledger version in convertSale"),
			"failed to ingest sale", logan.F{
				"actual_ledger_version": raw.Ext.V.ShortString(),
			}))
	}

	state, err := c.convertSaleState(rawState)
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

func (c *saleHandler) convertSaleState(state xdr.SaleState) (history.SaleState, error) {
	switch state {
	case xdr.SaleStateNone:
		return history.SaleStateOpen, nil
	case xdr.SaleStatePromotion:
		return history.SaleStatePromotion, nil
	case xdr.SaleStateVoting:
		return history.SaleStateVoting, nil
	default:
		return history.SaleStateOpen, errors.From(errors.New("unexpected sale of the state"), map[string]interface{}{
			"state": state,
		})
	}
}
