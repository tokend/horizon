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
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/log"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/tokend/go/signcontrol"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/regources"
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
	templateProxy := httputil.NewSingleHostReverseProxy(app.config.TemplateBackend)
	investReadyProxy := httputil.NewSingleHostReverseProxy(app.config.InvestReady)

	var telegramAirdropProxy *httputil.ReverseProxy
	if app.config.TelegramAirdrop != nil {
		telegramAirdropProxy = httputil.NewSingleHostReverseProxy(app.config.TelegramAirdrop)
	} else {
		telegramAirdropProxy = nil
	}

	operationTypesPayment := []xdr.OperationType{
		xdr.OperationTypePayment,
		xdr.OperationTypeCreateIssuanceRequest,
		xdr.OperationTypeCreateWithdrawalRequest,
		xdr.OperationTypeManageOffer,
		xdr.OperationTypeCheckSaleState,
		xdr.OperationTypeManageKeyValue,
		xdr.OperationTypePaymentV2,
	}

	r := app.web.router
	r.Get("/", &RootAction{})
	// system summary variables too verbose to be included into /
	r.Get("/statistics", &StatisticsAction{})
	r.Get("/statistics/balances", &BalancesReportAction{})
	r.Get("/metrics", &MetricsAction{})

	// ledger actions
	r.Get("/ledgers", &LedgerIndexAction{})
	r.Get("/ledgers/:id", &LedgerShowAction{})
	r.Get("/ledgers/:ledger_id/transactions", &TransactionIndexAction{})
	r.Get("/ledger_changes", &LedgerChangesAction{})

	// account actions
	r.Get("/accounts/:id", &AccountShowAction{})
	r.Get("/accounts/:id/limits", &LimitsV2AccountShowAction{})
	r.Get("/accounts/:id/signers", &SignersIndexAction{})
	r.Get("/accounts/:id/account_kyc", &AccountKYCAction{})
	r.Get("/accounts/:id/balances", &AccountBalancesAction{})
	r.Get("/accounts/:id/balances/details", &AccountDetailedBalancesAction{})
	r.Get("/accounts/:id/fees", &AccountFeesAction{})

	r.Get("/accounts/:account_id/signers/:id", &SignerShowAction{})
	r.Get("/accounts/:account_id/operations", &OperationIndexAction{}, 1)
	r.Get("/accounts/:account_id/payments", &OperationIndexAction{
		Types: operationTypesPayment,
	})
	r.Get("/accounts/:account_id/offers_history", &OperationIndexAction{
		Types: []xdr.OperationType{xdr.OperationTypeManageOffer},
	})
	r.Get("/accounts/:account_id/references", &CoreReferencesAction{})

	//keyValue actions
	r.Get("/key_value", &KeyValueShowAllAction{})
	r.Get("/key_value/:key", &KeyValueShowAction{})

	// offers
	r.Get("/accounts/:account_id/offers", &OffersAction{})

	// order book
	r.Get("/order_book", &OrderBookAction{})
	r.Get("/trades", &TradesAction{})

	r.Get("/trusts/:balance_id", &BalanceTrustsAction{})

	r.Get("/limits", &LimitsV2ShowAction{})

	// transaction history actions
	r.Get("/transactions", &TransactionIndexAction{})
	r.Get("/transactions/:id", &TransactionShowAction{})
	r.Get("/transactions/:tx_id/operations", &OperationIndexAction{})
	r.Get("/transactions/:tx_id/payments", &OperationIndexAction{
		Types: operationTypesPayment,
	})

	// transactions with ledger changes history action
	r.Get("/v2/transactions", &TransactionV2IndexAction{})

	// operation actions
	r.Get("/public/operations", &HistoryOperationIndexAction{})
	r.Get("/public/operations/:id", &HistoryOperationShowAction{})
	r.Get("/public/ledgers", &LedgerOperationsIndexAction{})

	r.Get("/operations", &OperationIndexAction{})
	r.Get("/payments", &OperationIndexAction{
		Types: operationTypesPayment,
	})
	r.Get("/operations/:id", &OperationShowAction{})

	//get fees action
	r.Get("/fees", &FeesAllAction{})
	r.Get("/fees_overview", &FeesAllAction{
		IsOverview: true,
	})
	r.Get("/fees/:fee_type", &FeesShowAction{})

	// assets
	r.Get("/charts/:code", &ChartsAction{})
	r.Get("/prices/history", &PricesHistoryAction{})
	r.Get("/assets", &AssetsIndexAction{})
	r.Get("/assets/:code", &AssetsShowAction{})
	r.Get("/assets/:code/holders", &AssetHoldersShowAction{})
	r.Get("/asset_pairs", &AssetPairsAction{})
	r.Get("/asset_pairs/convert", &AssetPairsConverterAction{})

	// balances
	r.Get("/balances", &BalanceIndexAction{})
	r.Get("/balances/:balance_id/asset", &BalanceAssetAction{})
	r.Get("/balances/:balance_id/account", &BalanceAccountAction{})

	// Reviewable Request actions
	r.Get("/requests/:id", &ReviewableRequestShowAction{})
	r.Get("/requests", &ReviewableRequestIndexAction{
		RequestTypes: xdr.ReviewableRequestTypeAll,
	})
	r.Get("/request/assets", &ReviewableRequestIndexAction{
		CustomFilter: func(action *ReviewableRequestIndexAction) {
			asset := action.GetString("asset")
			action.Page.Filters["asset"] = asset
			if asset != "" {
				action.q = action.q.AssetManagementByAsset(asset)
			}
		},
		RequestTypes: []xdr.ReviewableRequestType{xdr.ReviewableRequestTypeAssetCreate, xdr.ReviewableRequestTypeAssetUpdate},
	})
	r.Get("/request/preissuances", &ReviewableRequestIndexAction{
		CustomFilter: func(action *ReviewableRequestIndexAction) {
			asset := action.GetString("asset")
			action.Page.Filters["asset"] = asset
			if asset != "" {
				action.q = action.q.PreIssuanceByAsset(asset)
			}
		},
		RequestTypes: []xdr.ReviewableRequestType{xdr.ReviewableRequestTypePreIssuanceCreate},
	})
	r.Get("/request/issuances", &ReviewableRequestIndexAction{
		CustomFilter: func(action *ReviewableRequestIndexAction) {
			asset := action.GetString("asset")
			action.Page.Filters["asset"] = asset
			approvedCountingQ := action.HistoryQ().ReviewableRequests().CountQuery().ForTypes(action.RequestTypes)
			pendingCountingQ := action.HistoryQ().ReviewableRequests().CountQuery().ForTypes(action.RequestTypes)
			if asset != "" {
				action.q = action.q.IssuanceByAsset(asset)
				approvedCountingQ = approvedCountingQ.IssuanceByAsset(asset)
				pendingCountingQ = pendingCountingQ.IssuanceByAsset(asset)
			}

			action.Page.Embedded.Meta = &hal.PageMeta{
				Count: &regources.RequestsCount{},
			}

			var err error
			action.Page.Embedded.Meta.Count.Approved, err = approvedCountingQ.ForState(int64(history.ReviewableRequestStateApproved)).Count()
			if err != nil {
				action.Log.WithError(err).Error("failed to load count of approved reviewable requests")
				action.Err = &problem.ServerError
				return
			}

			action.Page.Embedded.Meta.Count.Pending, err = pendingCountingQ.ForState(int64(history.ReviewableRequestStatePending)).Count()
			if err != nil {
				action.Log.WithError(err).Error("failed to load count of pending reviewable requests")
				action.Err = &problem.ServerError
				return
			}
		},
		RequestTypes: []xdr.ReviewableRequestType{xdr.ReviewableRequestTypeIssuanceCreate},
	})
	r.Get("/request/withdrawals", &ReviewableRequestIndexAction{
		CustomFilter: func(action *ReviewableRequestIndexAction) {
			asset := action.GetString("dest_asset_code")
			action.Page.Filters["dest_asset_code"] = asset
			if asset != "" {
				action.q = action.q.WithdrawalByDestAsset(asset)
			}
		},
		RequestTypes: []xdr.ReviewableRequestType{xdr.ReviewableRequestTypeWithdraw, xdr.ReviewableRequestTypeTwoStepWithdrawal},
	})
	r.Get("/request/sales", &ReviewableRequestIndexAction{
		CustomFilter: func(action *ReviewableRequestIndexAction) {
			asset := action.GetString("base_asset")
			action.Page.Filters["base_asset"] = asset
			if asset != "" {
				action.q = action.q.SalesByBaseAsset(asset)
			}
		},
		RequestTypes: []xdr.ReviewableRequestType{xdr.ReviewableRequestTypeSale},
	})

	r.Get("/request/aml_alerts", &ReviewableRequestIndexAction{
		RequestTypes: []xdr.ReviewableRequestType{xdr.ReviewableRequestTypeAmlAlert},
	})

	r.Get("/request/limits_updates", &ReviewableRequestIndexAction{
		CustomFilter: func(action *ReviewableRequestIndexAction) {
			hash := action.GetString("document_hash")
			action.Page.Filters["document_hash"] = hash
			if hash != "" {
				action.q = action.q.LimitsByDocHash(hash)
			}
		},
		RequestTypes: []xdr.ReviewableRequestType{xdr.ReviewableRequestTypeLimitsUpdate},
	})
	r.Get("/request/update_kyc", &ReviewableRequestIndexAction{
		CustomCheckAllowed: func(action *ReviewableRequestIndexAction) {
			action.IsAllowed(action.GetString("account_to_update_kyc"), action.Reviewer, action.Requestor)
		},
		CustomFilter: func(action *ReviewableRequestIndexAction) {
			account := action.GetString("account_to_update_kyc")
			maskSet := action.GetInt64("mask_set")
			maskSetPartialEq := action.GetBool("mask_set_part_eq")
			maskNotSet := action.GetOptionalInt64("mask_not_set")
			accountTypeToSet := action.GetOptionalInt64("account_type_to_set")
			if action.Err != nil {
				return
			}
			action.Page.Filters["account_to_update_kyc"] = account
			action.Page.Filters["mask_set"] = action.GetString("mask_set")
			action.Page.Filters["mask_set_part_eq"] = action.GetString("mask_set_part_eq")
			action.Page.Filters["mask_not_set"] = action.GetString("mask_not_set")
			action.Page.Filters["account_type_to_set"] = action.GetString("account_type_to_set")

			if account != "" {
				action.q = action.q.KYCByAccountToUpdateKYC(account)
			}

			action.q = action.q.KYCByMaskSet(maskSet, maskSetPartialEq)
			if maskNotSet != nil {
				action.q = action.q.KYCByMaskNotSet(*maskNotSet)
			}

			if accountTypeToSet != nil {
				action.q = action.q.KYCByAccountTypeToSet(xdr.AccountType(*accountTypeToSet))
			}
		},
		RequestTypes: []xdr.ReviewableRequestType{xdr.ReviewableRequestTypeUpdateKyc},
	})
	r.Get("/request/update_sale_details", &ReviewableRequestIndexAction{
		RequestTypes: []xdr.ReviewableRequestType{xdr.ReviewableRequestTypeUpdateSaleDetails},
	})

	r.Get("/request/invoices", &ReviewableRequestIndexAction{
		CustomFilter: func(action *ReviewableRequestIndexAction) {
			contractID := action.GetOptionalInt64("contract_id")
			if contractID != nil {
				action.q = action.q.InvoicesByContract(*contractID)
			}
		},
		RequestTypes: []xdr.ReviewableRequestType{xdr.ReviewableRequestTypeInvoice},
	})

	r.Get("/request/contracts", &ReviewableRequestIndexAction{
		CustomFilter: func(action *ReviewableRequestIndexAction) {
			contractNumber := action.GetString("contract_number")
			if contractNumber != "" {
				action.q = action.q.ContractsByContractNumber(contractNumber)
			}

			// TODO: FIX ME!!!!!!!!!!!!!!!!!
			if action.Requestor != "" {
				action.q = action.q.ForCounterparty(action.Requestor)
			}

			action.Requestor = ""

			if action.Reviewer != "" {
				action.q = action.q.ForCounterparty(action.Reviewer)
			}

			action.Reviewer = ""
			action.Page.Filters["contract_number"] = contractNumber
		},
		RequestTypes: []xdr.ReviewableRequestType{xdr.ReviewableRequestTypeContract},
	})

	r.Get("/request/update_sale_end_time", &ReviewableRequestIndexAction{
		RequestTypes: []xdr.ReviewableRequestType{xdr.ReviewableRequestTypeUpdateSaleEndTime},
	})

	// Sales actions
	r.Get("/sales/:id", &SaleShowAction{})
	r.Get("/sales", &SaleIndexAction{})
	r.Get("/core_sales", &CoreSalesAction{})

	// Sale antes actions
	r.Get("/sale_antes", &SaleAnteAction{})

	// Contracts actions
	r.Get("/contracts", &ContractIndexAction{})
	r.Get("/contracts/:id", &ContractShowAction{})

	r.Post("/transactions", web.HandlerFunc(func(c web.C, w http.ResponseWriter, r *http.Request) {
		// DISCLAIMER: while following is true, it does not currently applies
		// API does not accept transactions make sure DisableAPISubmit is set to true
		//
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
			TransactionCreateAction{APIUrl: app.config.APIBackend}.ServeHTTPC(c, w, r)
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

	r.Handle(regexp.MustCompile(`^/templates/.*`), func() func(web.C, http.ResponseWriter, *http.Request) {
		return func(c web.C, w http.ResponseWriter, r *http.Request) {
			templateProxy.ServeHTTP(w, r)
		}
	}())

	r.Handle(regexp.MustCompile(`^/integrations/invest-ready`), func() func(web.C, http.ResponseWriter, *http.Request) {
		return func(c web.C, w http.ResponseWriter, r *http.Request) {
			investReadyProxy.ServeHTTP(w, r)
		}
	}())

	if telegramAirdropProxy != nil {
		r.Handle(regexp.MustCompile(`^/integrations/telegram-airdrop`), func() func(web.C, http.ResponseWriter, *http.Request) {
			return func(c web.C, w http.ResponseWriter, r *http.Request) {
				telegramAirdropProxy.ServeHTTP(w, r)
			}
		}())
	}

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
