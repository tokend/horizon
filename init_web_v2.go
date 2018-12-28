package horizon

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/horizon/v2"
	v2middleware "gitlab.com/tokend/horizon/v2/middleware"
	"net/http"
	"time"
)

type WebV2 struct {
	router  chi.Router
	metrics *v2.WebMetrics
}

func initWebV2(app *App) {
	router := chi.NewRouter()

	app.webV2.router = router
	app.webV2.metrics = v2.NewWebMetrics()
}

func initWebV2Middleware(app *App) {
	r := app.webV2.router

	logger := logan.New()

	r.Use(
		middleware.StripSlashes,
		middleware.SetHeader(upstreamHeader, app.config.Hostname),
		middleware.RequestID,
		ape.LoganMiddleware(logger, 300*time.Millisecond),
		ape.RecoverMiddleware(logger),
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

		r.Use(c.Handler)
	}

	signatureValidator := &SignatureValidator{app.config.SkipCheck}

	r.Use(signatureValidator.MiddlewareV2)
}

func initWebV2Actions(app *App) {
	r := app.webV2.router

	r.Get("/v2/accounts", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
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
