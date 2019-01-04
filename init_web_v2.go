package horizon

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/horizon/v2"
	"gitlab.com/tokend/horizon/v2/handlers"
	v2middleware "gitlab.com/tokend/horizon/v2/middleware"
	"time"
)

type WebV2 struct {
	mux     *v2.Mux
	metrics *v2.WebMetrics
}

func initWebV2(app *App) {
	router := chi.NewRouter()
	mux := v2.NewMux(router)

	app.webV2.mux = mux
	app.webV2.metrics = v2.NewWebMetrics()
}

func initWebV2Middleware(app *App) {
	m := app.webV2.mux

	logger := logan.New()

	m.Use(
		middleware.StripSlashes,
		middleware.SetHeader(upstreamHeader, app.config.Hostname),
		middleware.RequestID,
		ape.LoganMiddleware(logger, 300*time.Millisecond),
		ape.RecoverMiddleware(logger),
		ape.CtxMiddleWare(
			v2middleware.CtxCoreQ(app.coreQ),
			v2middleware.CtxHistoryQ(app.historyQ),
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
}

func initWebV2Actions(app *App) {
	m := app.webV2.mux

	m.Get("/v2/accounts/{id}", &handlers.AccountShow{})
}

func init() {
	appInit.Add(
		"web.init",
		initWebV2,
		"app-context", "stellarCoreInfo", "memory_cache",
	)

	appInit.Add(
		"web.middleware",
		initWebV2Middleware,
		"web.init",
		"web.metrics",
	)

	appInit.Add(
		"web.actions",
		initWebV2Actions,
		"web.init",
	)
}
