package horizon

import (
	"time"

	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/go/doorman"
	"gitlab.com/tokend/horizon/db2/core2"
	hdoorman "gitlab.com/tokend/horizon/doorman"
	"gitlab.com/tokend/horizon/log"
	"gitlab.com/tokend/horizon/web_v2"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/handlers"
	v2middleware "gitlab.com/tokend/horizon/web_v2/middleware"
)

type WebV2 struct {
	mux     *chi.Mux
	metrics *web_v2.WebMetrics
}

func initWebV2(app *App) {
	mux := chi.NewMux()
	metrics := web_v2.NewWebMetrics()

	app.webV2 = &WebV2{
		mux:     mux,
		metrics: metrics,
	}
}

func initWebV2Middleware(app *App) {
	m := app.webV2.mux

	logger := &log.DefaultLogger.Entry

	signersProvider := hdoorman.NewSignersQ(core2.NewSignerQ(app.CoreRepoLogged(nil)))

	m.Use(
		middleware.StripSlashes,
		middleware.SetHeader(upstreamHeader, app.config.Hostname),
		middleware.RequestID,
		ape.LoganMiddleware(logger, time.Second, ape.LoggerSetter(ctx.SetLog),
			ape.RequestIDProvider(middleware.GetReqID)),
		ape.RecoverMiddleware(logger),
		ape.CtxMiddleWare(
			// log will be set by logger setter on handler call
			ctx.SetCoreRepo(app.CoreRepoLogged(nil)),
			// log will be set by logger setter on handler call
			ctx.SetHistoryRepo(app.HistoryRepoLogged(nil)),
			ctx.SetDoorman(doorman.New(app.config.SkipCheck, signersProvider)),
			ctx.SetCoreInfo(*app.CoreInfo),
		),
		v2middleware.WebMetrics(app),
	)

	if app.config.CORSEnabled {
		// TODO: chi doesn't provide an analogue, should write own implementation?
		//r.Use(middleware.AutomaticOptions)

		c := cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedHeaders:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "DELETE", "PATCH", "PUT"},
			AllowCredentials: true,
		})

		m.Use(c.Handler)
	}

	signatureValidator := &SignatureValidator{app.config.SkipCheck}

	m.Use(signatureValidator.MiddlewareV2)
	m.NotFound(func(w http.ResponseWriter, r *http.Request) {
		err := problems.NotFound()
		err.Detail = "Unknown path"
		err.Meta = &map[string]interface{}{
			"url": r.URL.String(),
		}
		ape.RenderErr(w, err)
	})
}

func initWebV2Actions(app *App) {
	m := app.webV2.mux

	m.Get("/v3/accounts/{id}", handlers.GetAccount)
	m.Get("/v3/accounts/{id}/signers", handlers.GetAccountSigners)
	m.Get("/v3/accounts/{id}/calculated_fees", handlers.GetCalculatedFees)
	m.Get("/v3/assets/{code}", handlers.GetAsset)
	m.Get("/v3/assets", handlers.GetAssetList)
	m.Get("/v3/balances", handlers.GetBalanceList)
	m.Get("/v3/fees", handlers.GetFeeList)
	m.Get("/v3/history", handlers.GetHistory)
	m.Get("/v3/asset_pairs/{id}", handlers.GetAssetPair)
	m.Get("/v3/asset_pairs", handlers.GetAssetPairList)
	m.Get("/v3/offers/{id}", handlers.GetOffer)
	m.Get("/v3/offers", handlers.GetOfferList)
	m.Get("/v3/public_keys/{id}", handlers.GetPublicKey)

	// reviewable requests
	m.Get("/v3/requests", handlers.GetRequests)
	m.Get("/v3/requests/{id}", handlers.GetRequests)
	m.Get("/v3/create_asset_requests", handlers.GetCreateAssetRequests)
	m.Get("/v3/create_asset_requests/{id}", handlers.GetCreateAssetRequests)
	m.Get("/v3/create_sale_requests", handlers.GetCreateSaleRequests)
	m.Get("/v3/create_sale_requests/{id}", handlers.GetCreateSaleRequests)
	m.Get("/v3/update_asset_requests", handlers.GetUpdateAssetRequests)
	m.Get("/v3/update_asset_requests/{id}", handlers.GetUpdateAssetRequests)
	m.Get("/v3/create_pre_issuance_requests", handlers.GetCreatePreIssuanceRequests)
	m.Get("/v3/create_pre_issuance_requests/{id}", handlers.GetCreatePreIssuanceRequests)
	m.Get("/v3/create_issuance_requests", handlers.GetCreateIssuanceRequests)
	m.Get("/v3/create_issuance_requests/{id}", handlers.GetCreateIssuanceRequests)
	m.Get("/v3/create_withdraw_requests", handlers.GetCreateWithdrawRequests)
	m.Get("/v3/create_withdraw_requests/{id}", handlers.GetCreateWithdrawRequests)
	m.Get("/v3/update_limits_requests", handlers.GetUpdateLimitsRequests)
	m.Get("/v3/update_limits_requests/{id}", handlers.GetUpdateLimitsRequests)
	m.Get("/v3/create_aml_alert_requests", handlers.GetCreateAmlAlertRequests)
	m.Get("/v3/create_aml_alert_requests/{id}", handlers.GetCreateAmlAlertRequests)
	m.Get("/v3/change_role_requests", handlers.GetChangeRoleRequests)
	m.Get("/v3/change_role_requests/{id}", handlers.GetChangeRoleRequests)
	m.Get("/v3/update_sale_details_requests", handlers.GetUpdateSaleDetailsRequests)
	m.Get("/v3/update_sale_details_requests/{id}", handlers.GetUpdateSaleDetailsRequests)
	m.Get("/v3/create_atomic_swap_bid_requests", handlers.GetCreateAtomicSwapBidRequests)
	m.Get("/v3/create_atomic_swap_bid_requests/{id}", handlers.GetCreateAtomicSwapBidRequests)
	m.Get("/v3/create_atomic_swap_requests", handlers.GetCreateAtomicSwapRequests)
	m.Get("/v3/create_atomic_swap_requests{id}", handlers.GetCreateAtomicSwapRequests)

	m.Get("/v3/key_values", handlers.GetKeyValueList)
	m.Get("/v3/key_values/{key}", handlers.GetKeyValue)

	m.Get("/v3/sales", handlers.GetSaleList)
	m.Get("/v3/sales/{id}", handlers.GetSale)

	m.Get("/v3/order_book/{id}", handlers.GetOrderBook)

	m.Get("/v3/account_roles/{id}", handlers.GetAccountRole)
	m.Get("/v3/account_roles", handlers.GetAccountRoleList)
	m.Get("/v3/account_rules/{id}", handlers.GetAccountRule)
	m.Get("/v3/account_rules", handlers.GetAccountRuleList)

	m.Get("/v3/signer_roles/{id}", handlers.GetSignerRole)
	m.Get("/v3/signer_roles", handlers.GetSignerRoleList)
	m.Get("/v3/signer_rules/{id}", handlers.GetSignerRule)
	m.Get("/v3/signer_rules", handlers.GetSignerRuleList)

	janus := app.config.Janus()
	if err := janus.RegisterChi(m); err != nil {
		panic(errors.Wrap(err, "failed to register service"))
	}
}

func init() {
	appInit.Add(
		"web2.init",
		initWebV2,
		"app-context", "core-info", "memory_cache", "ledger-state",
	)

	appInit.Add(
		"web2.middleware",
		initWebV2Middleware,
		"web2.init",
		"web.metrics",
	)

	appInit.Add(
		"web2.actions",
		initWebV2Actions,
		"web2.middleware", "web2.init",
	)
}
