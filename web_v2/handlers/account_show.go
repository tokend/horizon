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
		return
	}

	response, err := BuildResource(r, account)
	if err != nil {
		RenderErr(w, err)
		return
	}

	err = RenderResource(w, *response)
	if err != nil {
		RenderErr(w, err)
		return
	}
}
