package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/requests"
	"gitlab.com/tokend/horizon/web_v2/resources"
	"gitlab.com/tokend/regources/rgenerated"
)

func GetKeyValue(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetKeyValue(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	keyValueQ := core2.NewKeyValueQ(ctx.CoreRepo(r))

	record, err := keyValueQ.ByKey(request.Key)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get key value")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if record == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	resource := resources.NewKeyValue(*record)
	response := &rgenerated.KeyValueEntryResponse{
		Data: resource,
	}

	ape.Render(w, response)
}
