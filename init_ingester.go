package horizon

import (
	"gitlab.com/swarmfund/horizon/log"
	"gitlab.com/swarmfund/horizon/ingest"
)

func initIngester(app *App) {
	if !app.config.Ingest {
		return
	}

	ingester := ingest.New(app.CoreConnector,
		app.CoreInfo,
		app.CoreRepo(nil),
		app.HistoryRepo(nil),
	)

	if err := ingester.IntegrityCheck(); err != nil {
		log.WithField("service", "ingester").
			WithError(err).
			Fatal("integrity check failed")
	}

	app.ingester = ingester
}

func init() {
	appInit.Add("ingester", initIngester, "core_connector", "app-context", "log", "horizon-db", "core-db", "stellarCoreInfo")
}
