package changes

import (
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	history "gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ingest2/internal"
)

type saleStorage interface {
	//Inserts sale into DB
	Insert(sale history.Sale) error
	//Updates sale
	Update(sale history.Sale) error
}

type saleHandler struct {
	storage saleStorage
}

func newSaleHandler(storage saleStorage) *saleHandler {
	return &saleHandler{
		storage: storage,
	}
}

//Created - handles creation of new sale
func (c *saleHandler) Created(lc ledgerChange) error {
	rawSale := lc.LedgerChange.MustCreated().Data.MustSale()

	sale, err := c.convertSale(rawSale)
	if err != nil {
		return errors.Wrap(err, "failed to convert sale", logan.F{
			"raw_sale":        rawSale,
			"ledger_sequence": lc.LedgerSeq,
		})
	}

	err = c.storage.Insert(*sale)
	if err != nil {
		return errors.Wrap(err, "failed to insert sale into DB", logan.F{
			"sale": sale,
		})
	}

	return nil
}

//Updated - handles update of the sale
func (c *saleHandler) Updated(lc ledgerChange) error {
	rawSale := lc.LedgerChange.MustUpdated().Data.MustSale()
	sale, err := c.convertSale(rawSale)
	if err != nil {
		return errors.Wrap(err, "failed to convert sale ", logan.F{
			"raw_sale":        rawSale,
			"ledger_sequence": lc.LedgerSeq,
		})
	}

	err = c.storage.Update(*sale)
	if err != nil {
		return errors.Wrap(err, "failed to update sale", logan.F{
			"sale": sale,
		})
	}
	return nil
}

func (c *saleHandler) convertSale(raw xdr.SaleEntry) (*history.Sale, error) {
	quoteAssets := make([]history.QuoteAsset, 0, len(raw.QuoteAssets))
	for i := range raw.QuoteAssets {
		quoteAssets = append(quoteAssets, history.QuoteAsset{
			Asset:          string(raw.QuoteAssets[i].QuoteAsset),
			Price:          amount.StringU(uint64(raw.QuoteAssets[i].Price)),
			QuoteBalanceID: raw.QuoteAssets[i].QuoteBalance.AsString(),
			CurrentCap:     amount.StringU(uint64(raw.QuoteAssets[i].CurrentCap)),
		})
	}

	saleType := raw.SaleTypeExt.SaleType

	return &history.Sale{
		ID:                uint64(raw.SaleId),
		OwnerAddress:      raw.OwnerId.Address(),
		BaseAsset:         string(raw.BaseAsset),
		DefaultQuoteAsset: string(raw.DefaultQuoteAsset),
		StartTime:         time.Unix(int64(raw.StartTime), 0).UTC(),
		EndTime:           time.Unix(int64(raw.EndTime), 0).UTC(),
		SoftCap:           uint64(raw.SoftCap),
		HardCap:           uint64(raw.HardCap),
		Details:           internal.MarshalCustomDetails(raw.Details),
		QuoteAssets: history.QuoteAssets{
			QuoteAssets: quoteAssets,
		},
		BaseCurrentCap: int64(raw.CurrentCapInBase),
		BaseHardCap:    int64(raw.MaxAmountToBeSold),
		SaleType:       saleType,
	}, nil
}
