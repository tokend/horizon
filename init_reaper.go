package horizon

import (
	"bullioncoin.githost.io/development/horizon/reap"
)

func initReaper(app *App) {
	app.reaper = reap.New(app.config.HistoryRetentionCount, app.HistoryRepo(nil))
}

func init() {
	appInit.Add("reaper", initReaper, "app-context", "log", "horizon-db")
}
