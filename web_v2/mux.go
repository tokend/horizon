package web_v2

import (
	"github.com/go-chi/chi"
	"net/http"
)

type Mux struct {
	router chi.Router
}

type Handler interface {
	Prepare(w http.ResponseWriter, r *http.Request)
	Render(w http.ResponseWriter, r *http.Request)
}

func NewMux(r chi.Router) *Mux {
	return &Mux{
		router: r,
	}
}

func (m *Mux) Get(pattern string, handler Handler) {
	m.router.Get(pattern, func(w http.ResponseWriter, r *http.Request) {
		handler.Prepare(w, r)
		handler.Render(w, r)
	})
}

func (m *Mux) Use(middlewares ...func(http.Handler) http.Handler) {
	m.router.Use(middlewares...)
}
