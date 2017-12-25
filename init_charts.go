package horizon

import (
	"gitlab.com/swarmfund/horizon/charts"
	"gitlab.com/swarmfund/horizon/db2"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon/db2/history"
	"fmt"
	"github.com/pkg/errors"
)

func initCharts(app *App) {
	historyStorage := charts.TxHistoryStorage{
		TxHistory: app.HistoryQ().Transactions(),
	}

	curs := db2.PageQuery{
		Cursor: "",
		Order:  "desc",
		Limit:  1,
	}

	var (
		records     []history.Transaction
		ledgerEntry []xdr.LedgerEntry
	)

	for {
		err := historyStorage.TxHistory.Page(curs).Select(&records)
		if err != nil {
			errors.Wrap(err, "Unable to select")
		}
		for _, record := range records {
			//record.LedgerCloseTime ---!!! достает время закрыя леджера,нужно для point::time
			curs.Cursor = record.PagingToken()

			assetChanges, err := process(record.TxMeta)
			if err != nil {
				errors.Wrap(err, "Unable to parse ledgerEntry")
			}

			for _, issued := range assetChanges{
				historyStorage.Run(uint64(issued.Issued), record.LedgerCloseTime)
			}
		}

		if len(records) == 0 {
			break
		}
	}

	for _, ledger := range ledgerEntry {
		fmt.Println(ledger.Data.ReviewableRequest.CreatedAt)
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
