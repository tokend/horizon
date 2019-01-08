package handlers

import (
	"github.com/go-chi/chi"
	"gitlab.com/tokend/horizon/web_v2/resource"
	"net/http"
)

type AccountShow struct {
	Base
}

func (a *AccountShow) Render(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	account, err := resource.NewAccount(id)
	if err != nil {
		a.RenderErr(w, err)
	}

	err = a.RenderResource(w, r, account)
	if err != nil {
		a.RenderErr(w, err)
		return
	}
}
