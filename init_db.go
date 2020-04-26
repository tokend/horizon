package horizon

import (
	"gitlab.com/tokend/horizon/bridge"
	"gitlab.com/tokend/horizon/db2/core"
	"gitlab.com/tokend/horizon/db2/history"
	"gitlab.com/tokend/horizon/log"
)

func initHorizonDb(app *App) {
	repo, err := bridge.Open(app.config.DatabaseURL)

	if err != nil {
		log.Panic(err)
	}

	app.historyQ = &history.Q{Mediator: repo}
}

func initCoreDb(app *App) {
	repo, err := bridge.Open(app.config.StellarCoreDatabaseURL)

	if err != nil {
		log.Panic(err)
	}

	app.coreQ = core.NewQ(repo)
}

func init() {
	appInit.Add("horizon-db", initHorizonDb, "app-context", "log")
	appInit.Add("core-db", initCoreDb, "app-context", "log")
}
