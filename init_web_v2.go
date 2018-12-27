package horizon

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rcrowley/go-metrics"
	"github.com/rs/cors"
	"net/http"
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

	r.Use(middleware.StripSlashes)
	// FIXME pls: use ctxMiddleware from ape here
	r.Use(
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				extender := app.CtxExtender()
				extender(ctx)
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		},
	)
	r.Use(middleware.SetHeader(upstreamHeader, app.config.Hostname))
	r.Use(middleware.RequestID)
	// FIXME pls: use loganMiddleware from ape
	r.Use(middleware.Logger)
	r.Use(requestMetricsMiddlewareV2)
	// FIXME pls: use recoverMiddleware from ape
	r.Use(middleware.Recoverer)
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
