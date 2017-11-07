package horizon

import (
	"net/http"

	"time"

	"bullioncoin.githost.io/development/api/log"
	"bullioncoin.githost.io/development/horizon/corer"
)

func initCoreConnector(app *App) {
	var err error
	app.CoreConnector, err = corer.NewConnector(&http.Client{
		Timeout: time.Duration(1) * time.Minute,
	}, app.config.StellarCoreURL)
	if err != nil {
		log.WithField("service", "initCoreConnector").WithError(err).Panic("Failed to create core connector")
	}
}

func init() {
	appInit.Add("core_connector", initCoreConnector, "log")
}
