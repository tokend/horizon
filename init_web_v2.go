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

	m.Route("/v3", func(r chi.Router) {
		registerV3Routes(r)
	})
	m.Route("/v4", func (r chi.Router) {
		registerV3Routes(r)
		registerV4Routes(r)
	})

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
