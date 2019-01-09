package handlers

import (
	"github.com/go-chi/chi"
	"gitlab.com/tokend/horizon/web_v2/resource"
	"net/http"
)

func ShowAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	account, err := resource.NewAccount(id)
	if err != nil {
		RenderErr(w, err)
	}

	err = RenderResource(w, r, account)
	if err != nil {
		RenderErr(w, err)
		return
	}
}
