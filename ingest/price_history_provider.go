package ingest

import (
	"strconv"
	"time"

	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/log"
	"gitlab.com/swarmfund/horizon/ingest/ingestion"
)

type assetPairPriceKey struct {
	BaseAsset  string
	QuoteAsset string
}

func newAssetPairPriceKey(baseAsset, quoteAsset string) assetPairPriceKey {
	return assetPairPriceKey{
		BaseAsset:  baseAsset,
		QuoteAsset: quoteAsset,
	}
}

type PriceHistoryProvider struct {
	prices          map[assetPairPriceKey]int64
	ledgerCloseTime time.Time

	log *log.Entry
}

func (h *PriceHistoryProvider) Init(assetPairs []core.AssetPair, ledgerCloseTime time.Time) {
	h.ledgerCloseTime = ledgerCloseTime
	h.log = log.WithField("service", "price_history_provider")
	h.prices = make(map[assetPairPriceKey]int64)

	for _, assetPair := range assetPairs {
		if assetPair.BaseAsset == assetPair.QuoteAsset {
			continue
		}

		h.Put(assetPair.BaseAsset, assetPair.QuoteAsset, assetPair.CurrentPrice)
	}
}

func (h *PriceHistoryProvider) Put(base, quote string, price int64) {
	h.prices[newAssetPairPriceKey(base, quote)] = price
}

func (h *PriceHistoryProvider) ToPricePoints() ([]ingestion.LedgerPricePoint, error) {
	result := make([]ingestion.LedgerPricePoint, 0, len(h.prices))

	for assetPair, pricePoint := range h.prices {
		price, err := strconv.ParseFloat(amount.String(pricePoint), 64)
		if err != nil {
			h.log.WithError(err).Error("Failed to get price history")
			return nil, err
		}

		result = append(result, ingestion.LedgerPricePoint{
			BaseAsset:  assetPair.BaseAsset,
			QuoteAsset: assetPair.QuoteAsset,
			PricePoint: history.PricePoint{
				Price:     price,
				Timestamp: h.ledgerCloseTime,
			},
		})
	}

	return result, nil
}
