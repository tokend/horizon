package handlers

import (
	"net/http"

	"gitlab.com/tokend/horizon/web_v2/resources"

	"gitlab.com/tokend/horizon/db2/history2"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/requests"
	regources "gitlab.com/tokend/regources/generated"
)

func GetDataList(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetDataList(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	hRepo := ctx.HistoryRepo(r)
	handler := getDataListHandler{
		DataQ: history2.NewDataQ(hRepo),
		Log:   ctx.Log(r),
	}

	dataOwners := []*string{}

	if request.ShouldFilter(requests.FilterTypeDataListOwner) {
		dataOwners = append(dataOwners, &request.Filters.Owner)
	}

	if !isAllowed(r, w, dataOwners...) {
		return
	}

	response, err := handler.GetDataList(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get data")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, response)
}

type getDataListHandler struct {
	DataQ    history2.DataQ
	AccountQ core2.AccountsQ
	Log      *logan.Entry
}

func (h *getDataListHandler) GetDataList(request *requests.GetDataList) (*regources.DataListResponse, error) {
	q := h.DataQ

	if request.ShouldFilter(requests.FilterTypeDataListOwner) {
		q = q.FilterByOwner(request.Filters.Owner)
	}

	if request.ShouldFilter(requests.FilterTypeDataListType) {
		q = q.FilterByType(request.Filters.Type)
	}

	q = q.Page(request.PageParams)

	dataSet, err := q.Select()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get data list")
	}

	response := regources.DataListResponse{
		Data: make([]regources.Data, 0, len(dataSet)),
	}

	for _, dataEntry := range dataSet {
		response.Data = append(response.Data, resources.NewData(dataEntry))

		if request.ShouldInclude(requests.IncludeTypeDataOwner) {
			owner, err := h.AccountQ.GetByAddress(dataEntry.Owner)
			if err != nil {
				return nil, err
			}
			if owner == nil {
				return nil, errors.New("owner not found")
			}
			ownerAccount := resources.NewAccount(*owner, nil)
			response.Included.Add(&ownerAccount)
		}
	}

	if len(response.Data) > 0 {
		response.Links = request.GetCursorLinks(*request.PageParams, response.Data[len(response.Data)-1].ID)
	} else {
		response.Links = request.GetCursorLinks(*request.PageParams, "")
	}

	return &response, nil
}
