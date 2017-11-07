package horizon

import (
	"bullioncoin.githost.io/development/horizon/db2"
	"bullioncoin.githost.io/development/horizon/db2/core"
	"bullioncoin.githost.io/development/horizon/db2/history"
	"bullioncoin.githost.io/development/horizon/log"
)

func initHorizonDb(app *App) {
	repo, err := db2.Open(app.config.DatabaseURL)

	if err != nil {
		log.Panic(err)
	}
	repo.DB.SetMaxIdleConns(4)
	repo.DB.SetMaxOpenConns(12)

	app.historyQ = &history.Q{Repo: repo}
}

func initCoreDb(app *App) {
	repo, err := db2.Open(app.config.StellarCoreDatabaseURL)

	if err != nil {
		log.Panic(err)
	}

	repo.DB.SetMaxIdleConns(4)
	repo.DB.SetMaxOpenConns(12)
	app.coreQ = core.NewQ(repo)
}

func init() {
	appInit.Add("horizon-db", initHorizonDb, "app-context", "log")
	appInit.Add("core-db", initCoreDb, "app-context", "log")
}
