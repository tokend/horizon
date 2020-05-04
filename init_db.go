package horizon

import (
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/tokend/horizon/db2/core"
	"gitlab.com/tokend/horizon/db2/history"
	"gitlab.com/tokend/horizon/log"
)

func initHorizonDb(app *App) {
	repo, err := pgdb.Open(pgdb.Opts{
		URL:                app.config.DatabaseURL,
		MaxOpenConnections: 12,
		MaxIdleConnections: 4,
	})

	if err != nil {
		log.Panic(err)
	}

	app.historyQ = &history.Q{DB: repo}
}

func initCoreDb(app *App) {
	repo, err := pgdb.Open(pgdb.Opts{
		URL:                app.config.StellarCoreDatabaseURL,
		MaxOpenConnections: 12,
		MaxIdleConnections: 4,
	})

	if err != nil {
		log.Panic(err)
	}

	app.coreQ = core.NewQ(repo)
}

func init() {
	appInit.Add("horizon-db", initHorizonDb, "app-context", "log")
	appInit.Add("core-db", initCoreDb, "app-context", "log")
}
