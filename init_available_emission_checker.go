package horizon

import (
	"time"

	"gitlab.com/swarmfund/horizon/errors"
	"gitlab.com/swarmfund/horizon/log"
)

func notifyForAsset(app *App, asset string) error {
	err := app.notificator.SendLowAvailableEmissions(app.config.Notificator.AdminNotification.EmissionNotificationReceiver, asset)
	if err != nil {
		return err
	}
	return err
}

func loadAssetsAndCheckPreEmissions(app *App) {
	defer func() {
		if rec := recover(); rec != nil {
			err := errors.FromPanic(rec)
			log.WithStack(err).WithError(err).Errorf("admin emission notification for asset failed")
		}
	}()

	availableEmissions, err := app.obtainAvailableEmissions()
	if err != nil {
		log.Error(err)
		time.Sleep(10 * time.Second)
		return
	}
	for asset, amount := range availableEmissions {
		if amount < int64(app.config.Notificator.AdminNotification.EmissionThreshold) {
			err := notifyForAsset(app, asset)
			if err != nil {
				log.WithField("error", err.Error()).Warn("Failed to send available emission notification")
				time.Sleep(60 * time.Second)
			}
		}
	}
	time.Sleep(60 * time.Second)

}

func checkAvailableEmissions(app *App) {
	for {
		loadAssetsAndCheckPreEmissions(app)
	}
}

func initAvailableEmissionChecker(app *App) {
	go checkAvailableEmissions(app)
}

func init() {
	appInit.Add("available_emission_checker", initAvailableEmissionChecker, "core-db", "notificator")
}
