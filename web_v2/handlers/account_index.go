package handlers

import (
	"github.com/go-chi/chi"
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

func (a *AccountIndex) Render(w http.ResponseWriter, r *http.Request) {
	a.collection = &resource.AccountCollection{}

	err := a.PrepareCollection(r, a.collection)
	if err != nil {
		a.RenderErr(w, err)
		return
	}

	a.filters.accountType = chi.URLParam(r, "account_type")
	a.pagingParams.Limit = chi.URLParam(r, "limit")
	a.pagingParams.Page = chi.URLParam(r, "page")

	err = a.RenderCollection(w, r, a.pagingParams, a.collection)
	if err != nil {
		a.RenderErr(w, err)
		return
	}
}
