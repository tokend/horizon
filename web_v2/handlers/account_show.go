package handlers

import (
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/horizon/web_v2/resource"
	"net/http"
)

func ShowAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	account, err := resource.NewAccount(id)
	if err != nil {
		Log(r).WithError(err).Error("Failed to create account resource")
		RenderErr(r, w, err)
		return
	}

	response, err := BuildResource(r, account)
	if err != nil {
		RenderErr(r, w, err)
		return
	}

	err = RenderResource(w, *response)
	if err != nil {
		Log(r).WithError(err).Error("Failed to render account resource")
		RenderErr(r, w, problems.InternalError())
		return
	}
}
