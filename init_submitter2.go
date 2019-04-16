package horizon

import (
	"time"

	"gitlab.com/tokend/horizon/ingest2/storage"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"

	"github.com/lib/pq"
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/horizon/db2/history2"

	"gitlab.com/tokend/horizon/log"
	"gitlab.com/tokend/horizon/txsub/v2"
)

func initSubmissionV2(app *App) {
	logger := log.WithField("service", "submitter_v2")
	cq := core2.NewTransactionQ(app.CoreRepoLogged(&logger.Entry))
	hq := history2.NewTransactionsQ(app.HistoryRepoLogged(&logger.Entry))
	coreConnector := app.CoreConnector

	listener := pq.NewListener(
		app.config.StellarCoreDatabaseURL,
		time.Second,
		5*time.Second,
		log.PQEvent(logger),
	)
	err := listener.Listen(storage.ChanNewLedgerSeq)
	if err != nil {
		panic(errors.Wrap(err, "failed to init history listener", logan.F{
			"channel": storage.ChanNewLedgerSeq,
		}))
	}
	app.submitterV2 = &txsub.System{
		Log:               log.WithField("service", "txsub2"),
		SubmissionTimeout: time.Minute,
		List:              txsub.NewDefaultSubmissionList(10 * time.Second),
		Submitter:         txsub.NewDefaultSubmitter(*coreConnector),
		Listener:          listener,
		Results: &txsub.ResultsProvider{
			Core:    cq,
			History: hq,
		},
		NetworkPassphrase: app.CoreInfo.NetworkPassphrase,
	}

	app.submitterV2.Start(app.ctx)
}

func init() {
	appInit.Add("txsub2", initSubmissionV2, "app-context", "log", "horizon-db", "core-db", "core-info")
}
