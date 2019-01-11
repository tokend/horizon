package horizon

import (
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ledger"
	"gitlab.com/tokend/horizon/log"
)

func initLedgerState(app *App) {
	logger := log.WithField("service", "ledger_state_updater")
	ledger.StartLedgerStateUpdater(app.ctx, logger, ledger.Config{
		CoreDB:    app.config.StellarCoreDatabaseURL,
		HistoryDB: app.config.DatabaseURL,
		Core:      core2.NewLedgerHeaderQ(app.CoreRepoLogged(&logger.Entry)),
		History:   app.HistoryQ(),
		History2:  history2.NewLedgerQ(app.HistoryRepoLogged(&logger.Entry)),
	})
}

func init() {
	appInit.Add("ledger-state", initLedgerState, "core_connector", "app-context", "log", "horizon-db", "core-db")
}
