package handlers

import (
	"github.com/go-chi/chi"
	"gitlab.com/tokend/go/signcontrol"
	"gitlab.com/tokend/horizon/web_v2/resource"
	"net/http"
)

type AccountIndex struct {
	Base

	collection *resource.AccountCollection
	filters    struct {
		accountType string
		signerType  int32
		isBlocked   bool
	}
}

func (a *AccountIndex) Prepare(w http.ResponseWriter, r *http.Request) {
	a.filters.accountType = chi.URLParam(r, "account_type")

	a.collection.Signer, _ = signcontrol.CheckSignature(r)
	a.collection.W = w
	a.collection.R = r

	err := a.CheckAllowed(r, a.collection)
	if err != nil {
		a.RenderErr()
		return
	}
}

func (a *AccountIndex) Render(w http.ResponseWriter, r *http.Request) {
	a.Base.RenderCollection(w, r, a.collection)
}
