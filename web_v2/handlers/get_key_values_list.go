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

func GetKeyValueList(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetKeyValueList(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	keyValueQ := core2.NewKeyValueQ(ctx.CoreRepo(r))

	records, err := keyValueQ.Page(request.PageParams).Select()
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get key values")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	response := &rgenerated.KeyValueEntrysResponse{
		Data:  make([]rgenerated.KeyValueEntry, 0, len(records)),
		Links: request.GetOffsetLinks(*request.PageParams),
	}

	for _, record := range records {
		resource := resources.NewKeyValue(record)
		response.Data = append(response.Data, resource)
	}

	ape.Render(w, response)
}
