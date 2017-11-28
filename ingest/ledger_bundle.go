package ingest

import (
	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/swarmfund/horizon/db2/history"
	"strconv"
	"time"
)

// Load runs queries against `core` to fill in the records of the bundle.
func (lb *LedgerBundle) Load(coreQ core.QInterface, historyQ history.QInterface) error {

	// Load Header
	err := coreQ.LedgerHeaderBySequence(&lb.Header, lb.Sequence)
	if err != nil {
		return err
	}

	// Load transactions
	err = coreQ.TransactionsByLedger(&lb.Transactions, lb.Sequence)
	if err != nil {
		return err
	}

	err = coreQ.TransactionFeesByLedger(&lb.TransactionFees, lb.Sequence)
	if err != nil {
		return err
	}

	assetPairs, err := getAssetPairs(coreQ, historyQ)
	if err != nil {
		return err
	}

	lb.HistoryPriceProvide = new(PriceHistoryProvider)
	lb.HistoryPriceProvide.Init(assetPairs, time.Unix(lb.Header.CloseTime, 0).UTC())

	return nil
}

func getAssetPairs(coreQ core.QInterface, historyQ history.QInterface) (assetPairs []core.AssetPair, err error) {
	assetPairs, err = coreQ.AssetPairs()
	if err != nil {
		return assetPairs, err
	}

	for key, ap := range assetPairs {
		assetPairs[key].CurrentPrice, err = getPriceFromHistory(historyQ, ap)
		if err != nil {
			return assetPairs, err
		}
	}

	return assetPairs, nil
}

func getPriceFromHistory(historyQ history.QInterface, assetPair core.AssetPair) (int64, error) {
	price, err := historyQ.LastPrice(assetPair.BaseAsset, assetPair.QuoteAsset)
	if err != nil {
		return 0, err
	}

	// if the price is not in history db
	// set default value - One
	if price == nil {
		return amount.One, nil
	}

	priceStr := strconv.FormatFloat(price.Price, 'f', 10, 64)

	xPrice, err := amount.Parse(priceStr)
	if err != nil {
		return 0, err
	}

	return int64(xPrice), nil
}
