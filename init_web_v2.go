package horizon

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rcrowley/go-metrics"
	"github.com/rs/cors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/logan/v3"
	"time"
)

type WebV2 struct {
	router chi.Router

	requestCounter metrics.Counter
	failureCounter metrics.Counter

	requestTimer metrics.Timer
	failureMeter metrics.Meter
	successMeter metrics.Meter
}

func initWebV2(app *App) {
	router := chi.NewRouter()

	app.webV2.router = router
}

func initWebV2Middleware(app *App) {
	r := app.webV2.router

	logger := logan.New()

	r.Use(
		middleware.StripSlashes,
		ape.CtxMiddleWare(
			app.CtxExtender(),
		),
		middleware.SetHeader(upstreamHeader, app.config.Hostname),
		middleware.RequestID,
		ape.LoganMiddleware(logger, 300*time.Millisecond),
		ape.RecoverMiddleware(logger),
		requestMetricsMiddlewareV2,
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

		r.Use(c.Handler)
	}

	signatureValidator := &SignatureValidator{app.config.SkipCheck}

	r.Use(signatureValidator.MiddlewareV2)
}

func initWebV2Actions(app *App) {
	// TODO
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
