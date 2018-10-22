package ingest

import (
	"strconv"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/horizon/db2/core"
	"gitlab.com/tokend/horizon/db2/history"
)

// Load runs queries against `core` to fill in the records of the bundle.
func (lb *LedgerBundle) Load(coreQ core.QInterface) error {

	// Load Header
	err := coreQ.LedgerHeaderBySequence(&lb.Header, lb.Sequence)
	if err != nil {
		return errors.Wrap(err, "failed to get ledger header")
	}

	// Load transactions
	err = coreQ.TransactionsByLedger(&lb.Transactions, lb.Sequence)
	if err != nil {
		return errors.Wrap(err, "failed to get transactions for ledger")
	}

	err = coreQ.TransactionFeesByLedger(&lb.TransactionFees, lb.Sequence)
	if err != nil {
		return errors.Wrap(err, "failed to get transaction fees for ledger")
	}

	return nil
}

func getAssetPairs(coreQ core.QInterface, historyQ history.QInterface) (assetPairs []core.AssetPair, err error) {
	assetPairs, err = coreQ.AssetPairs().Select()
	if err != nil {
		return assetPairs, errors.Wrap(err, "failed to select asset pairs")
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
		return 0, errors.Wrap(err, "failed to get last price")
	}

	// if the price is not in history db
	// set default value - One
	if price == nil {
		return amount.One, nil
	}

	priceStr := strconv.FormatFloat(price.Price, 'f', 10, 64)

	xPrice, err := amount.Parse(priceStr)
	if err != nil {
		return 0, errors.Wrap(err, "failed to parse amount")
	}

	return int64(xPrice), nil
}
