package handlers

import (
	"net/http"

	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/horizon/web_v2/resources"

	"gitlab.com/tokend/go/xdr"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/requests"
	"gitlab.com/tokend/regources/v2"
)

func GetCreateAtomicSwapRequests(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetCreateAtomicSwapRequests(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	historyRepo := ctx.HistoryRepo(r)
	coreRepo := ctx.CoreRepo(r)
	handler := getCreateAtomicSwapRequestsHandler{
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

type getCreateAtomicSwapRequestsHandler struct {
	R         requests.GetCreateAtomicSwapRequests
	Base      getRequestListBaseHandler
	RequestsQ history2.ReviewableRequestsQ
	AssetsQ   core2.AssetsQ
	Log       *logan.Entry
}

func (h *getCreateAtomicSwapRequestsHandler) MakeAll(w http.ResponseWriter, request requests.GetCreateAtomicSwapRequests) error {
	q := h.RequestsQ.FilterByRequestType(uint64(xdr.ReviewableRequestTypeCreateAtomicSwap))

	if request.ShouldFilter(requests.FilterTypeCreateAtomicSwapRequestsQuoteAsset) {
		q = q.FilterByAtomicSwapQuoteAsset(request.Filters.QuoteAsset)
	}

	return h.Base.SelectAndRender(w, *request.GetRequestsBase, q, h.RenderRecord)
}

func (h *getCreateAtomicSwapRequestsHandler) RenderRecord(included *regources.Included, record history2.ReviewableRequest) (regources.ReviewableRequest, error) {
	resource := h.Base.PopulateResource(*h.R.GetRequestsBase, included, record)

	if h.R.ShouldInclude(requests.IncludeTypeCreateAtomicSwapRequestsQuoteAsset) {
		asset, err := h.AssetsQ.FilterByCode(record.Details.CreateAtomicSwap.QuoteAsset).Get()
		if err != nil {
			return regources.ReviewableRequest{}, errors.Wrap(err, "failed to get quote asset")
		}
		if asset == nil {
			return regources.ReviewableRequest{}, errors.New("quote asset not found")
		}

		resource := resources.NewAsset(*asset)
		included.Add(&resource)
	}

	return resource, nil
}
