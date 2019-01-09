package handlers

import (
	"github.com/go-chi/chi"
	"gitlab.com/tokend/horizon/web_v2/resource"
	"net/http"
)

func ShowAccountList(w http.ResponseWriter, r *http.Request) {
	collection := &resource.AccountCollection{}
	collection.Filters.AccountType = chi.URLParam(r, "account_type")

	err := RenderCollection(w, r, collection)
	if err != nil {
		RenderErr(w, err)
		return
	}
}
