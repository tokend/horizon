package horizon

import (
	"gitlab.com/tokend/horizon/log"
	"gitlab.com/tokend/horizon/notificator"
)

func initNotificator(app *App) {
	var err error
	app.notificator = notificator.NewConnector(&app.config.Notificator)
	if err != nil {
		log.WithField("service", "notificator").Fatal("failed to create notificator")
	}
}

func init() {
	appInit.Add("notificator", initNotificator)
}
