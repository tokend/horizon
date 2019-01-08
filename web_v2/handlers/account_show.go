package handlers

import (
	"github.com/go-chi/chi"
	"gitlab.com/tokend/horizon/web_v2/resource"
	"net/http"
)

type AccountShow struct {
	Base
	resource *resource.Account
}

func (a *AccountShow) Render(w http.ResponseWriter, r *http.Request) {
	accountId := chi.URLParam(r, "id")

	account, err := resource.NewAccount(accountId)
	if err != nil {
		a.RenderErr(w, err)
	}
	a.resource = account

	err = a.PrepareResource(r, a.resource)
	if err != nil {
		a.RenderErr(w, err)
		return
	}

	err = a.Base.RenderResource(w, r, accountId, a.resource)
	if err != nil {
		a.RenderErr(w, err)
		return
	}
}
