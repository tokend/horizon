package handlers

import (
	"github.com/go-chi/chi"
	"gitlab.com/tokend/horizon/web_v2/resource"
	"net/http"
)

type AccountIndex struct {
	Base
}

func (a *AccountIndex) Render(w http.ResponseWriter, r *http.Request) {
	collection := &resource.AccountCollection{}
	collection.Filters.AccountType = chi.URLParam(r, "account_type")

	err := a.RenderCollection(w, r, collection)
	if err != nil {
		a.RenderErr(w, err)
		return
	}
}
