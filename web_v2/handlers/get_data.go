package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/horizon/web_v2/resources"
	regources "gitlab.com/tokend/regources/generated"

	"gitlab.com/tokend/horizon/db2/history2"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/horizon/web_v2/ctx"

	"gitlab.com/tokend/horizon/web_v2/requests"
)

func GetData(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetData(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	hrepo := ctx.HistoryRepo(r)
	h := getDataHandler{
		DataQ:    history2.NewDataQ(hrepo),
		AccountQ: core2.NewAccountsQ(ctx.CoreRepo(r)),
		log:      ctx.Log(r),
	}

	result, err := h.DataQ.GetByID(request.DataID)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get data")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if result == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if !isAllowed(r, w, &result.Owner) {
		return
	}

	response := regources.DataResponse{
		Data: resources.NewData(*result),
	}

	if request.ShouldInclude(requests.IncludeTypeDataOwner) {

		ownerAccount := resources.NewAccountKey(result.Owner)
		response.Included.Add(&ownerAccount)
	}

	ape.Render(w, response)
}

type getDataHandler struct {
	DataQ    history2.DataQ
	AccountQ core2.AccountsQ

	log *logan.Entry
}
