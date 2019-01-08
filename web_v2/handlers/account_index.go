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

	pagingParams resource.PagingParams
}

func (a *AccountIndex) Prepare(w http.ResponseWriter, r *http.Request) {
	a.filters.accountType = chi.URLParam(r, "account_type")

	a.pagingParams.Limit = chi.URLParam(r, "limit")
	a.pagingParams.Page = chi.URLParam(r, "page")

	a.collection.Signer, _ = signcontrol.CheckSignature(r)

	err := a.CheckAllowed(a.collection)
	if err != nil {
		a.RenderErr(w, err)
		return
	}
}

func (a *AccountIndex) Render(w http.ResponseWriter, r *http.Request) {
	a.Prepare(w, r)
	a.Base.RenderCollection(w, r, a.pagingParams, a.collection)
}
