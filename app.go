package horizon

import (
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"

	"gitlab.com/distributed_lab/tokend/horizon/cache"
	"gitlab.com/distributed_lab/tokend/horizon/config"
	"gitlab.com/distributed_lab/tokend/horizon/corer"
	"gitlab.com/distributed_lab/tokend/horizon/db2"
	"gitlab.com/distributed_lab/tokend/horizon/db2/core"
	"gitlab.com/distributed_lab/tokend/horizon/db2/history"
	"gitlab.com/distributed_lab/tokend/horizon/ingest"
	"gitlab.com/distributed_lab/tokend/horizon/ledger"
	"gitlab.com/distributed_lab/tokend/horizon/log"
	"gitlab.com/distributed_lab/tokend/horizon/notificator"
	"gitlab.com/distributed_lab/tokend/horizon/reap"
	"gitlab.com/distributed_lab/tokend/horizon/render/sse"
	"github.com/garyburd/redigo/redis"
	"github.com/pkg/errors"
	"github.com/rcrowley/go-metrics"
	"gitlab.com/distributed_lab/txsub"
	"golang.org/x/net/context"
	"golang.org/x/net/http2"
	"gopkg.in/tylerb/graceful.v1"
)

// You can override this variable using: gb build -ldflags "-X main.version aabbccdd"
var version = ""

// App represents the root of the state of a horizon instance.
type App struct {
	config         config.Config
	web            *Web
	historyQ       history.QInterface
	coreQ          core.QInterface
	ctx            context.Context
	cancel         func()
	redis          *redis.Pool
	submitter      *txsub.System
	ingester       *ingest.System
	reaper         *reap.System
	ticks          *time.Ticker
	CoreInfo       *corer.Info
	horizonVersion string
	cacheProvider  *cache.Provider
	notificator    *notificator.Connector
	CoreConnector  *corer.Connector

	// metrics
	metrics                  metrics.Registry
	historyLatestLedgerGauge metrics.Gauge
	historyElderLedgerGauge  metrics.Gauge
	horizonConnGauge         metrics.Gauge
	coreLatestLedgerGauge    metrics.Gauge
	coreElderLedgerGauge     metrics.Gauge
	coreConnGauge            metrics.Gauge
	goroutineGauge           metrics.Gauge
}

// SetVersion records the provided version string in the package level `version`
// var, which will be used for the reported horizon version.
func SetVersion(v string) {
	version = v
}

// NewApp constructs an new App instance from the provided config.
func NewApp(config config.Config) (*App, error) {
	result := &App{config: config}
	result.horizonVersion = version
	result.ticks = time.NewTicker(1 * time.Second)
	result.init()
	return result, nil
}

// Serve starts the horizon web server, binding it to a socket, setting up
// the shutdown signals.
func (a *App) Serve() {

	a.web.router.Compile()
	http.Handle("/", a.web.router)

	addr := fmt.Sprintf(":%d", a.config.Port)

	srv := &graceful.Server{
		Timeout: 10 * time.Second,

		Server: &http.Server{
			Addr:    addr,
			Handler: http.DefaultServeMux,
		},

		ShutdownInitiated: func() {
			log.Info("received signal, gracefully stopping")
			a.Close()
		},
	}

	http2.ConfigureServer(srv.Server, nil)

	log.Infof("Starting horizon on %s", addr)

	go a.run()

	var err error
	if a.config.TLSCert != "" {
		err = srv.ListenAndServeTLS(a.config.TLSCert, a.config.TLSKey)
	} else {
		err = srv.ListenAndServe()
	}

	if err != nil {
		log.Panic(err)
	}

	log.Info("stopped")
}

// Close cancels the app and forces the closure of db connections
func (a *App) Close() {
	a.cancel()
	a.ticks.Stop()

	a.historyQ.GetRepo().DB.Close()
	a.coreQ.GetRepo().DB.Close()

}

// CoreQ returns a helper object for performing sql queries against the
// stellar core database.
func (a *App) CoreQ() core.QInterface {
	return a.coreQ
}

// HistoryQ returns a helper object for performing sql queries against the
// history portion of horizon's database.
func (a *App) HistoryQ() history.QInterface {
	return a.historyQ
}

// CoreRepo returns a new repo that loads data from the stellar core
// database. The returned repo is bound to `ctx`.
func (a *App) CoreRepo(ctx context.Context) *db2.Repo {
	return &db2.Repo{DB: a.coreQ.GetRepo().DB, Ctx: ctx}
}

// HistoryRepo returns a new repo that loads data from the horizon database. The
// returned repo is bound to `ctx`.
func (a *App) HistoryRepo(ctx context.Context) *db2.Repo {
	return &db2.Repo{DB: a.historyQ.GetRepo().DB, Ctx: ctx}
}

func (action *Action) Notificator() *notificator.Connector {
	return action.App.notificator
}

// IsHistoryStale returns true if the latest history ledger is more than
// `StaleThreshold` ledgers behind the latest core ledger
func (a *App) IsHistoryStale() bool {
	if a.config.StaleThreshold == 0 {
		return false
	}

	ls := ledger.CurrentState()
	return (ls.CoreLatest - ls.HistoryLatest) > int32(a.config.StaleThreshold)
}

// UpdateLedgerState triggers a refresh of several metrics gauges, such as open
// db connections and ledger state
func (a *App) UpdateLedgerState() {
	var err error
	var next ledger.State

	err = a.CoreQ().LatestLedger(&next.CoreLatest)
	if err != nil {
		goto Failed
	}

	err = a.CoreQ().ElderLedger(&next.CoreElder)
	if err != nil {
		goto Failed
	}

	err = a.HistoryQ().LatestLedger(&next.HistoryLatest)
	if err != nil {
		goto Failed
	}

	err = a.HistoryQ().ElderLedger(&next.HistoryElder)
	if err != nil {
		goto Failed
	}

	err = ledger.SetState(next)
	if err != nil {
		log.WithField("err", err.Error()).Error("core is hanging")
	}

	return

Failed:
	log.WithStack(err).
		WithField("err", err.Error()).
		Error("failed to load ledger state")

}

// UpdateStellarCoreInfo updates the value of coreVersion and networkPassphrase
// from the Stellar core API.
func (a *App) UpdateStellarCoreInfo() {
	if a.config.StellarCoreURL == "" {
		return
	}

	var err error
	a.CoreInfo, err = a.CoreConnector.GetCoreInfo()
	if err != nil {
		log.WithField("service", "core-info").WithError(err).Error("could not load stellar-core info")
		return
	}
}

// UpdateMetrics triggers a refresh of several metrics gauges, such as open
// db connections and ledger state
func (a *App) UpdateMetrics() {
	a.goroutineGauge.Update(int64(runtime.NumGoroutine()))
	ls := ledger.CurrentState()
	a.historyLatestLedgerGauge.Update(int64(ls.HistoryLatest))
	a.historyElderLedgerGauge.Update(int64(ls.HistoryElder))
	a.coreLatestLedgerGauge.Update(int64(ls.CoreLatest))
	a.coreElderLedgerGauge.Update(int64(ls.CoreElder))

	//a.horizonConnGauge.Update(int64(a.historyQ.Repo.DB.Stats().OpenConnections))
	//a.coreConnGauge.Update(int64(a.coreQ.Repo.DB.Stats().OpenConnections))
}

// DeleteUnretainedHistory forwards to the app's reaper.  See
// `reap.DeleteUnretainedHistory` for details
func (a *App) DeleteUnretainedHistory() error {
	return a.reaper.DeleteUnretainedHistory()
}

// Tick triggers horizon to update all of it's background processes such as
// transaction submission, metrics, ingestion and reaping.
func (a *App) Tick() {
	var wg sync.WaitGroup
	log.Debug("ticking app")
	// update ledger state and stellar-core info in parallel
	wg.Add(2)
	go func() { a.UpdateLedgerState(); wg.Done() }()
	go func() { a.UpdateStellarCoreInfo(); wg.Done() }()
	wg.Wait()

	if a.ingester != nil {
		go a.ingester.Tick()
	}

	wg.Add(2)
	go func() { a.reaper.Tick(); wg.Done() }()
	go func() { a.submitter.Tick(a.ctx); wg.Done() }()
	wg.Wait()

	sse.Tick()

	// finally, update metrics
	a.UpdateMetrics()
	log.Debug("finished ticking app")
}

// Init initializes app, using the config to populate db connections and
// whatnot.
func (a *App) init() {
	appInit.Run(a)
}

// run is the function that runs in the background that triggers Tick each
// second
func (a *App) run() {
	for {
		select {
		case <-a.ticks.C:
			a.Tick()
		case <-a.ctx.Done():
			log.Info("finished background ticker")
			return
		}
	}
}

func (a *App) obtainAvailableEmissions() (map[string]int64, error) {
	result := make(map[string]int64)

	emissions, err := a.CoreQ().AvailableEmissions(a.CoreInfo.MasterAccountID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get core emissions")
	}

	for _, emission := range emissions {
		result[emission.Asset] = emission.Amount
	}

	return result, nil
}
