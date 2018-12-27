package middleware

import (
	"github.com/go-chi/chi/middleware"
	"gitlab.com/tokend/horizon"
	"net/http"
	"time"
)

func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app := r.Context().Value("app").(*horizon.App)

		mw := middleware.NewWrapResponseWriter(w, 1)
		ts := time.Now()
		next.ServeHTTP(mw.(http.ResponseWriter), r)

		app.UpdateWebV2Metrics(time.Since(ts), mw.Status())
	})
}
