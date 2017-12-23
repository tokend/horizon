package horizon

import (
	"database/sql"
	"net/http"
	"net/http/httputil"
	"regexp"

	"github.com/rcrowley/go-metrics"
	"github.com/rs/cors"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
	"gitlab.com/swarmfund/go/signcontrol"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon/log"
	"gitlab.com/swarmfund/horizon/render/problem"
)

// Web contains the http server related fields for horizon: the router,
// rate limiter, etc.
type Web struct {
	router *RateLimitedMux

	requestCounter metrics.Counter
	failureCounter metrics.Counter

	requestTimer metrics.Timer
	failureMeter metrics.Meter
	successMeter metrics.Meter
}

// initWeb installed a new Web instance onto the provided app object.
func initWeb(app *App) {
	mux, err := NewRateLimitedMux(app)
	if err != nil {
		log.WithField("service", "web").WithError(err).Fatal("failed to init mux")
	}
	app.web = &Web{
		router: mux,

		requestCounter: metrics.NewCounter(),
		failureCounter: metrics.NewCounter(),

		requestTimer: metrics.NewTimer(),
		failureMeter: metrics.NewMeter(),
		successMeter: metrics.NewMeter(),
	}

	// register problems
	problem.RegisterError(sql.ErrNoRows, problem.NotFound)
}

// initWebMiddleware installs the middleware stack used for horizon onto the
// provided app.
func initWebMiddleware(app *App) {
	r := app.web.router

	r.Use(stripTrailingSlashMiddleware())
	r.Use(middleware.EnvInit)
	r.Use(app.Middleware)
	r.Use(UpstreamMiddleware)
	r.Use(middleware.RequestID)
	r.Use(contextMiddleware(app.ctx))
	r.Use(LoggerMiddleware)
	r.Use(requestMetricsMiddleware)
	r.Use(RecoverMiddleware)

	if app.config.CORSEnabled {
		r.Use(middleware.AutomaticOptions)

		c := cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedHeaders:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "DELETE", "PATCH", "PUT"},
			AllowCredentials: true,
		})

		r.Use(c.Handler)
	}

	signatureValidator := &SignatureValidator{app.config.SkipCheck}

	r.Use(signatureValidator.Middleware)
}

const (
	levelLow int = 1 << (1 * iota)
	levelMid
	levelHigh
	levelCritical
)

const (
	IsAdminHeader      = "x-is-admin"
	IsAdminHeaderValue = "1"
	IsSignedHeader     = "x-is-signed"
	IsSignedValue      = "1"
)

// initWebActions installs the routing configuration of horizon onto the
// provided app.  All route registration should be implemented here.
func initWebActions(app *App) {
	apiProxy := httputil.NewSingleHostReverseProxy(app.config.APIBackend)
	keychainProxy := httputil.NewSingleHostReverseProxy(app.config.KeychainBackend)

	operationTypesPayment := []xdr.OperationType{
		xdr.OperationTypePayment,
		xdr.OperationTypeCreateIssuanceRequest,
		xdr.OperationTypeCreateWithdrawalRequest,
		xdr.OperationTypeManageOffer,
		xdr.OperationTypeManageInvoice,
	}

	r := app.web.router
	r.Get("/", &RootAction{})
	r.Get("/metrics", &MetricsAction{})

	// ledger actions
	r.Get("/ledgers", &LedgerIndexAction{})
	r.Get("/ledgers/:id", &LedgerShowAction{})
	r.Get("/ledgers/:ledger_id/transactions", &TransactionIndexAction{})
	r.Get("/ledger_changes", &LedgerChangesAction{})

	// account actions
	r.Get("/accounts/:id", &AccountShowAction{})
	r.Get("/accounts/:id/signers", &SignersIndexAction{})
	r.Get("/accounts/:id/summary", &AccountSummaryAction{})
	r.Get("/accounts/:id/balances", &AccountBalancesAction{})
	r.Get("/accounts/:account_id/signers/:id", &SignerShowAction{})
	r.Get("/accounts/:account_id/operations", &OperationIndexAction{}, 1)
	r.Get("/accounts/:account_id/payments", &OperationIndexAction{
		Types: operationTypesPayment,
	})

	// offers
	r.Get("/accounts/:account_id/offers", &OffersAction{})

	// order book
	r.Get("/order_book", &OrderBookAction{})
	r.Get("/trades", &TradesAction{})

	r.Get("/trusts/:balance_id", &BalanceTrustsAction{})

	r.Get("/default_limits", &AccountTypeLimitsAllAction{})
	r.Get("/default_limits/:account_type", &AccountTypeLimitsShowAction{})

	// transaction history actions
	r.Get("/transactions", &TransactionIndexAction{})
	r.Get("/transactions/:id", &TransactionShowAction{})
	r.Get("/transactions/:tx_id/operations", &OperationIndexAction{})
	r.Get("/transactions/:tx_id/payments", &OperationIndexAction{
		Types: operationTypesPayment,
	})

	// operation actions
	r.Get("/public/operations", &HistoryOperationIndexAction{})
	r.Get("/public/operations/:id", &HistoryOperationShowAction{})

	r.Get("/operations", &OperationIndexAction{})
	r.Get("/payments", &OperationIndexAction{
		Types: operationTypesPayment,
	})
	r.Get("/operations/:id", &OperationShowAction{})

	r.Get("/payment_requests", &PaymentRequestIndexAction{})
	r.Get("/forfeit_requests", &PaymentRequestIndexAction{
		OnlyForfeits: true,
	})

	r.Get("/payment_requests/:id", &PaymentRequestShowAction{})

	//get fees action
	r.Get("/fees", &FeesAllAction{})
	r.Get("/fees_overview", &FeesAllAction{
		IsOverview: true,
	})
	r.Get("/fees/:fee_type", &FeesShowAction{})

	// assets
	r.Get("/prices/history", &PricesHistoryAction{})
	r.Get("/assets", &AssetsIndexAction{})
	r.Get("/assets/:code", &AssetsShowAction{})
	r.Get("/asset_pairs", &AssetPairsAction{})
	r.Get("/asset_pairs/convert", &AssetPairsConverterAction{})

	// balances
	r.Get("/balances", &BalanceIndexAction{})
	r.Get("/balances/:balance_id/asset", &BalanceAssetAction{})
	r.Get("/balances/:balance_id/account", &BalanceAccountAction{})

	// Reviewable Request actions
	r.Get("/requests/:id", &ReviewableRequestShowAction{})
	r.Get("/requests", &ReviewableRequestIndexAction{})
	r.Get("/request/withdrawals", &WithdrawalIndexAction{})

	// Sales actions
	r.Get("/sales/:id", &SaleShowAction{})
	r.Get("/sales", &SaleIndexAction{})

	r.Post("/transactions", web.HandlerFunc(func(c web.C, w http.ResponseWriter, r *http.Request) {
		// legacy constraints:
		// * not signed POST /transactions should trigger TFA flow if needed
		// * not signed POST /transactions should eventually make network submission
		//
		// signed request submission flow:
		// * horizon accepts only admin signed request and submits them directly
		//   omitting pending transaction flow
		//
		// not signed request submission flow:
		// * horizon accepts not signed request and proxy it to api
		// * api handles request and make it's thing about TFA and stuff
		// * api eventually submits transaction with admin signature to horizon
		// * api returns response from horizon as-is
		// * api handles pending transactions silently, cleaning up after time bounds
		//  expired
		//
		// that solves most of horizon/api abstractions issues but also leads to
		// not-so-obvious flow for some transactions, basically rules are:
		// * clients should sign their submission requests if, and only if, their
		//   intent is to by-pass implicit pending transactions flow
		// * requests with transactions of user account type sources must not be
		//   signed by client or they will be rejected

		// checking if request is signed and deciding on proper handler
		// (we rely on SignatureValidator middleware here)
		signer := r.Header.Get(signcontrol.PublicKeyHeader)
		if signer != "" || app.config.DisableAPISubmit {
			TransactionCreateAction{}.ServeHTTPC(c, w, r)
		} else {
			apiProxy.ServeHTTP(w, r)
		}
	}))

	r.Get("/accounts/:account_id/transactions", web.HandlerFunc(func(c web.C, w http.ResponseWriter, r *http.Request) {
		// while current implementation is clearly lame, more important is to make
		// public API clear and intuitive since it's impossible to change it later
		query := r.URL.Query()
		if query.Get("pending") != "" {
			apiProxy.ServeHTTP(w, r)
		} else {
			TransactionIndexAction{}.ServeHTTPC(c, w, r)
		}
	}))

	r.Handle(regexp.MustCompile(`^/users/\w+/keys`), func() func(web.C, http.ResponseWriter, *http.Request) {
		return func(c web.C, w http.ResponseWriter, r *http.Request) {
			keychainProxy.ServeHTTP(w, r)
		}
	}())

	// proxy pass every request horizon could not handle to API
	r.Handle(regexp.MustCompile(`^.*`), func() func(web.C, http.ResponseWriter, *http.Request) {
		return func(c web.C, w http.ResponseWriter, r *http.Request) {
			apiProxy.ServeHTTP(w, r)
		}
	}())
}

func init() {
	appInit.Add(
		"web.init",
		initWeb,
		"app-context", "stellarCoreInfo", "memory_cache",
	)

	appInit.Add(
		"web.middleware",
		initWebMiddleware,
		"web.init",
		"web.metrics",
	)

	appInit.Add(
		"web.actions",
		initWebActions,
		"web.init",
	)
}
