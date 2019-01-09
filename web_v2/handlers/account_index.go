package handlers

import (
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/resource"
	"net/http"
)

func ShowAccountCollection(w http.ResponseWriter, r *http.Request) {
	collection, err := resource.NewAccountCollection()
	if err != nil {
		ctx.Log(r).WithError(err).Error("Failed to create account resource")
		RenderErr(r, w, err)
		return
	}

	collection.Filters.AccountType = chi.URLParam(r, "account_type")

	response, err := BuildCollection(r, collection)
	if err != nil {
		RenderErr(r, w, err)
		return
	}

	err = RenderCollection(w, response)
	if err != nil {
		ctx.Log(r).WithError(err).Error("Failed to render account resource")
		RenderErr(r, w, problems.InternalError())
		return
	}

}
