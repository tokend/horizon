package horizon

import (
	"net/http"

	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/horizon/db2/history2"

	"time"

	"gitlab.com/distributed_lab/corer"
	"gitlab.com/tokend/horizon/log"
	"gitlab.com/tokend/horizon/txsub/v2"
)

func initSubmissionV2(app *App) {
	logger := &log.WithField("service", "initSubmissionV2").Entry
	cq := core2.NewTransactionQ(app.CoreRepoLogged(logger))
	hq := history2.NewTransactionQ(app.HistoryRepoLogged(logger))
	coreConnector, err := corer.NewConnector(&http.Client{
		Timeout: time.Duration(1 * time.Minute),
	}, app.config.StellarCoreURL)
	if err != nil {
		logger.WithError(err).Panic("Failed to create core connector")
	}
	app.submitterV2 = &txsub.System{
		TickPeriod: time.Second,
		Pending:    txsub.NewDefaultSubmissionList(3 * time.Second),
		Submitter:  txsub.NewDefaultSubmitter(coreConnector),
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
