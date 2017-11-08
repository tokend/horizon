package horizon

import (
	"time"

	"gitlab.com/tokend/horizon/log"
	"github.com/multiplay/go-slack/chat"
	"github.com/multiplay/go-slack/lrhook"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

// initLog initialized the logging subsystem, attaching app.log and
// app.logMetrics.  It also configured the logger's level using Config.LogLevel.
func initLog(app *App) {
	log.DefaultLogger.Logger.Level = app.config.LogLevel
	log.DefaultLogger = log.DefaultLogger.WithField("host", app.config.Hostname)

	if app.config.LogToJSON {
		log.DefaultLogger.Entry.Logger.Formatter = &logrus.JSONFormatter{}
	}

	if app.config.SlackWebhook != nil {
		cfg := lrhook.Config{
			MinLevel: logrus.WarnLevel,
			Limit:    rate.Every(1 * time.Second),
			Message: chat.Message{
				Channel:   app.config.SlackChannel,
				IconEmoji: ":glitch_crab:",
			},
		}
		h := lrhook.New(cfg, app.config.SlackWebhook.String())
		log.DefaultLogger.Logger.Hooks.Add(h)
	}
}

func init() {
	appInit.Add("log", initLog)
}
