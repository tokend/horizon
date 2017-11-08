package horizon

import (
	"fmt"

	"time"

	"gitlab.com/tokend/horizon/log"
	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/rcrowley/go-metrics"
)

func initMetrics(app *App) {
	app.metrics = metrics.NewRegistry()
}

func initDbMetrics(app *App) {
	app.historyLatestLedgerGauge = metrics.NewGauge()
	app.historyElderLedgerGauge = metrics.NewGauge()
	app.coreLatestLedgerGauge = metrics.NewGauge()
	app.coreElderLedgerGauge = metrics.NewGauge()

	app.horizonConnGauge = metrics.NewGauge()
	app.coreConnGauge = metrics.NewGauge()
	app.goroutineGauge = metrics.NewGauge()
	app.metrics.Register("history.latest_ledger", app.historyLatestLedgerGauge)
	app.metrics.Register("history.elder_ledger", app.historyElderLedgerGauge)
	app.metrics.Register("stellar_core.latest_ledger", app.coreLatestLedgerGauge)
	app.metrics.Register("stellar_core.elder_ledger", app.coreElderLedgerGauge)
	app.metrics.Register("history.open_connections", app.horizonConnGauge)
	app.metrics.Register("stellar_core.open_connections", app.coreConnGauge)
	app.metrics.Register("goroutines", app.goroutineGauge)
}

func initIngesterMetrics(app *App) {
	if app.ingester == nil {
		return
	}
	app.metrics.Register("ingester.ingest_ledger",
		app.ingester.Metrics.IngestLedgerTimer)
	app.metrics.Register("ingester.clear_ledger",
		app.ingester.Metrics.ClearLedgerTimer)
}

func initLogMetrics(app *App) {
	for level, meter := range *log.DefaultMetrics {
		key := fmt.Sprintf("logging.%s", level)
		app.metrics.Register(key, meter)
	}
}

func initTxSubMetrics(app *App) {
	app.submitter.Init()
	app.metrics.Register("txstub.buffered", app.submitter.Metrics.BufferedSubmissionsGauge)
	app.metrics.Register("txsub.open", app.submitter.Metrics.OpenSubmissionsGauge)
	app.metrics.Register("txsub.succeeded", app.submitter.Metrics.SuccessfulSubmissionsMeter)
	app.metrics.Register("txsub.failed", app.submitter.Metrics.FailedSubmissionsMeter)
	app.metrics.Register("txsub.total", app.submitter.Metrics.SubmissionTimer)
}

// initWebMetrics registers the metrics for the web server into the provided
// app's metrics registry.
func initWebMetrics(app *App) {
	app.metrics.Register("requests.total", app.web.requestTimer)
	app.metrics.Register("requests.succeeded", app.web.successMeter)
	app.metrics.Register("requests.failed", app.web.failureMeter)
}

func initInflux(app *App) {
	client, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     "http://localhost:8186",
		Username: "root",
		Password: "root",
	})
	if err != nil {
		log.WithField("service", "metrics").WithError(err).Error("failed to enable influx")
		return
	}
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		entry := log.WithField("service", "influx.metrics")
		lastSubmit := time.Now().UTC()
		for {
			<-ticker.C
			now := time.Now().UTC()
			// get current metric values
			failedRequests := app.web.failureCounter.Count()
			totalRequests := app.web.requestCounter.Count()
			timerSnapshot := app.web.requestTimer.Snapshot()
			// TODO revert it on error
			app.web.requestTimer = metrics.NewTimer()

			// create a new point batch
			points, err := influx.NewBatchPoints(influx.BatchPointsConfig{
				Database: "telegraf",
			})
			if err != nil {
				entry.WithError(err).Error("failed to create points batch")
				return
			}

			tags := map[string]string{"type": "horizon"}
			fields := map[string]interface{}{
				"rps":         float64(totalRequests) / now.Sub(lastSubmit).Seconds(),
				"failed_rps":  float64(failedRequests) / now.Sub(lastSubmit).Seconds(),
				"response_95": timerSnapshot.Percentile(95) / 1000000000,
				"response_20": timerSnapshot.Percentile(20) / 1000000000,
			}

			point, err := influx.NewPoint("stats", tags, fields, now)
			if err != nil {
				entry.WithError(err).Error("failed to create point")
				return
			}

			points.AddPoint(point)

			if err := client.Write(points); err != nil {
				entry.WithError(err).Error("failed to submit points")
				return
			}

			lastSubmit = now

			// decrease counter values
			app.web.failureCounter.Dec(failedRequests)
			app.web.requestCounter.Dec(totalRequests)
		}
	}()

}

func init() {
	appInit.Add("metrics", initMetrics)
	appInit.Add("log.metrics", initLogMetrics, "metrics")
	appInit.Add("db-metrics", initDbMetrics, "metrics", "horizon-db", "core-db")
	appInit.Add("web.metrics", initWebMetrics, "web.init", "metrics")
	appInit.Add("txsub.metrics", initTxSubMetrics, "txsub", "metrics")
	appInit.Add("ingester.metrics", initIngesterMetrics, "ingester", "metrics")
	appInit.Add("influx.metrics", initInflux, "web.init")
}
