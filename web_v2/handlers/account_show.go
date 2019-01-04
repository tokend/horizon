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

	err := a.CheckAllowed(r, a.resource)
	if err != nil {
		a.RenderErr()
		return
	}
}

func (a *AccountShow) Serve(w http.ResponseWriter, r *http.Request) {
	err := a.resource.PopulateModel()
	if err != nil {
		a.RenderErr()
		return
	}

	js, err := a.resource.MarshalModel()
	if err != nil {
		a.RenderErr()
		return
	}

	_, err = w.Write(js)
	if err != nil {
		a.RenderErr()
		return
	}
}

func (a *AccountShow) RenderErr() {
	// TODO
}
