package handlers

import (
	"github.com/go-chi/chi"
	"gitlab.com/tokend/go/signcontrol"
	"gitlab.com/tokend/horizon/web_v2/resource"
	"net/http"
)

type AccountShow struct {
	Base
	resource *resource.Account

	filters  struct {
		accountId string
	}
}

func (a *AccountShow) Prepare(w http.ResponseWriter, r *http.Request) {
	a.filters.accountId = chi.URLParam(r, "id")

	a.resource = &resource.Account{}
	a.resource.Signer, _ = signcontrol.CheckSignature(r)
	a.resource.W = w
	a.resource.R = r

	err := a.CheckAllowed(a.resource)
	if err != nil {
		a.RenderErr()
		return
	}
}

func (a *AccountShow) Render(w http.ResponseWriter, r *http.Request) {
	a.Prepare(w, r)
	a.Base.RenderResource(w, r, a.filters.accountId, a.resource)
}
