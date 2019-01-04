package web_v2

import (
	"github.com/go-chi/chi"
	"gitlab.com/tokend/horizon/web_v2/handlers"
	"net/http"
)

type Mux struct {
	router chi.Router
}

func NewMux(r chi.Router) *Mux {
	return &Mux{
		router: r,
	}
}

func (m *Mux) Get(pattern string, handler handlers.Handler) {
	m.router.Get(pattern, func(w http.ResponseWriter, r *http.Request) {
		handler.Prepare(w, r)
		handler.Serve(w, r)
	})
}

func (m *Mux) Use(middlewares ...func(http.Handler) http.Handler) {
	m.router.Use(middlewares...)
}
