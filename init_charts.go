package horizon

import (
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon/charts"
	"gitlab.com/swarmfund/horizon/db2"
	"gitlab.com/swarmfund/horizon/db2/history"
)

func initCharts(app *App) {

	var histogram charts.Histogram

	hourDuration := time.Hour * 24

	histogram = *charts.NewHistogram(hourDuration, 24)

	historyStorage := charts.TxHistoryStorage{
		TxHistory: app.HistoryQ().Transactions(),
	}

	curs := db2.PageQuery{
		Cursor: "",
		Order:  "desc",
		Limit:  1,
	}

	var records []history.Transaction

	for {
		err := historyStorage.TxHistory.Page(curs).Select(&records)
		if err != nil {
			logrus.WithError(err).Error("Unable to select")
			return
		}
		for _, record := range records {
			curs.Cursor = record.PagingToken()

			assetChanges, err := process(record.TxMeta)
			if err != nil {
				logrus.WithError(err).Error("Unable to parse ledgerEntry")
				return
			}

			for _, issued := range assetChanges {
				histogram.Run(uint64(issued.Issued), record.LedgerCloseTime)
			}
		}

		if len(records) == 0 {
			break
		}
	}
}

func process(txMeta string) ([]xdr.AssetEntry, error) {

	var assetEntry []xdr.AssetEntry

	var transactionMeta xdr.TransactionMeta
	err := xdr.SafeUnmarshalBase64(txMeta, &transactionMeta)
	if err != nil {
		return assetEntry, errors.Wrap(err, "errUnmarshalBase64")
	}

	operations := transactionMeta.MustOperations()

	for _, changes := range operations {
		for _, change := range changes.Changes {
			switch change.Type {
			case xdr.LedgerEntryChangeTypeUpdated:
				switch change.Updated.Data.Type {
				case xdr.LedgerEntryTypeAsset:
					assetEntry = append(assetEntry, *change.Updated.Data.Asset)
				}
			}
		}
	}

	return assetEntry, nil
}

func init() {
	appInit.Add("swarm_horizon", initCharts, "horizon-db")
}
