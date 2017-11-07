package horizon

import (
	"net/http"

	"time"

	"bullioncoin.githost.io/development/api/log"
	"bullioncoin.githost.io/development/horizon/db2/core"
	"bullioncoin.githost.io/development/horizon/db2/history"
	hTxsub "bullioncoin.githost.io/development/horizon/txsub"
	"gitlab.com/distributed_lab/corer"
	"gitlab.com/distributed_lab/txsub"
)

func initSubmissionSystem(app *App) {
	cq := &core.Q{Repo: app.CoreRepo(nil)}
	hq := &history.Q{Repo: app.HistoryRepo(nil)}
	coreConnector, err := corer.NewConnector(&http.Client{
		Timeout: time.Duration(1 * time.Minute),
	}, app.config.StellarCoreURL)
	if err != nil {
		log.WithField("service", initSubmissionSystem).WithError(err).Panic("Failed to create core connector")
	}
	app.submitter = &txsub.System{
		Pending:   txsub.NewDefaultSubmissionList(),
		Submitter: txsub.NewDefaultSubmitter(coreConnector),
		Results: &hTxsub.ResultsProvider{
			Core:    cq,
			History: hq,
		},
		NetworkPassphrase: app.CoreInfo.NetworkPassphrase,
	}
}

func init() {
	appInit.Add("txsub", initSubmissionSystem, "app-context", "log", "horizon-db", "core-db", "stellarCoreInfo")
}
