package horizon

import (
	"gitlab.com/distributed_lab/tokend/horizon/reap"
)

func initReaper(app *App) {
	app.reaper = reap.New(app.config.HistoryRetentionCount, app.HistoryRepo(nil))
}

func init() {
	appInit.Add("reaper", initReaper, "app-context", "log", "horizon-db")
}
