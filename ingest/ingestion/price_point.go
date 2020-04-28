package ingestion

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2/history"
)

type LedgerPricePoint struct {
	BaseAsset  string
	QuoteAsset string
	history.PricePoint
}

func (ingest *Ingestion) StorePricePoints(priceHistory []LedgerPricePoint) error {
	if len(priceHistory) == 0 {
		return nil
	}

	q := ingest.priceHistory
	for _, price := range priceHistory {
		q = q.Values(price.BaseAsset, price.QuoteAsset, price.Timestamp, price.Price)
	}

	err := ingest.DB.Exec(q)
	if err != nil {
		return errors.Wrap(err, "failed to insert price points")
	}

	return nil
}
