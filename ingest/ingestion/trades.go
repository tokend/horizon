package ingestion

import (
	"gitlab.com/distributed_lab/logan/v2/errors"
	"gitlab.com/swarmfund/go/xdr"
	"time"
)

func (ingest *Ingestion) StoreTrades(source xdr.AccountId, result xdr.ManageOfferResult, ledgerCloseTime int64) error {
	if result.Success == nil || len(result.Success.OffersClaimed) == 0 {
		return nil
	}

	q := ingest.trades
	for i := range result.Success.OffersClaimed {
		claimed := result.Success.OffersClaimed[i]
		q = q.Values(string(result.Success.BaseAsset),
			string(result.Success.QuoteAsset), int64(claimed.BaseAmount),
			int64(claimed.QuoteAmount), int64(claimed.CurrentPrice), time.Unix(ledgerCloseTime, 0).UTC())
	}

	_, err := ingest.DB.Exec(q)
	if err != nil {
		return errors.Wrap(err, "failed to store trades")
	}

	return nil
}
