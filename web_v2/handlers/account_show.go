package handlers

import (
	"github.com/go-chi/chi"
	"gitlab.com/tokend/horizon/web_v2/resource"
	"net/http"
)

type AccountShow struct {
	Base
	resource *resource.Account

	filters struct {
		accountId string
	}
}

func (a *AccountShow) Render(w http.ResponseWriter, r *http.Request) {
	a.resource = &resource.Account{}

	err := a.PrepareResource(r, a.resource)
	if err != nil {
		a.RenderErr(w, err)
		return
	}

	a.filters.accountId = chi.URLParam(r, "id")

	err = a.Base.RenderResource(w, r, a.filters.accountId, a.resource)
	if err != nil {
		a.RenderErr(w, err)
		return
	}
}
