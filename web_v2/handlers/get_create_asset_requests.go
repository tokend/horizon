package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/web_v2/resources"

	"gitlab.com/tokend/horizon/db2/core2"

	"gitlab.com/tokend/go/xdr"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/requests"
	"gitlab.com/tokend/regources/rgenerated"
)

func GetCreateAssetRequests(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetCreateAssetRequests(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	historyRepo := ctx.HistoryRepo(r)
	coreRepo := ctx.CoreRepo(r)
	handler := getCreateAssetRequestsHandler{
		R:         request,
		RequestsQ: history2.NewReviewableRequestsQ(historyRepo),
		AssetsQ:   core2.NewAssetsQ(coreRepo),
		Log:       ctx.Log(r),
	}

	if !isAllowed(r, w, request.GetRequestsBase.Filters.Requestor, request.GetRequestsBase.Filters.Reviewer) {
		return
	}

	err = handler.MakeAll(w, request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get request list", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}
}

type getCreateAssetRequestsHandler struct {
	R         requests.GetCreateAssetRequests
	Base      getRequestListBaseHandler
	RequestsQ history2.ReviewableRequestsQ
	AssetsQ   core2.AssetsQ
	Log       *logan.Entry
}

func (h *getCreateAssetRequestsHandler) MakeAll(w http.ResponseWriter, request requests.GetCreateAssetRequests) error {
	q := h.RequestsQ.FilterByRequestType(uint64(xdr.ReviewableRequestTypeCreateAsset))

	if request.ShouldFilter(requests.FilterTypeCreateAssetRequestsAsset) {
		q = q.FilterByAssetCreateAsset(request.Filters.Asset)
	}

	return h.Base.SelectAndRender(w, *request.GetRequestsBase, q, h.RenderRecord)
}

func (h *getCreateAssetRequestsHandler) RenderRecord(included *rgenerated.Included, record history2.ReviewableRequest) (rgenerated.ReviewableRequest, error) {
	resource := h.Base.PopulateResource(*h.R.GetRequestsBase, included, record)

	if h.R.ShouldInclude(requests.IncludeTypeCreateAssetRequestsAsset) {
		asset, err := h.AssetsQ.GetByCode(record.Details.CreateAsset.Asset)
		if err != nil {
			return rgenerated.ReviewableRequest{}, errors.Wrap(err, "failed to get asset")
		}

		if asset == nil {
			return rgenerated.ReviewableRequest{}, errors.New("asset not found")
		}
		resource := resources.NewAsset(*asset)
		included.Add(&resource)
	}
	return resource, nil
}
