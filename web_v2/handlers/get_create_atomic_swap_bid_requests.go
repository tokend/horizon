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
	"gitlab.com/tokend/regources/v2/generated"
)

func GetCreateAtomicSwapBidRequests(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetCreateAtomicSwapBidRequests(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	historyRepo := ctx.HistoryRepo(r)
	coreRepo := ctx.CoreRepo(r)
	handler := getCreateAtomicSwapBidRequestsHandler{
		R:         request,
		RequestsQ: history2.NewReviewableRequestsQ(historyRepo),
		BalancesQ: core2.NewBalancesQ(coreRepo),
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

type getCreateAtomicSwapBidRequestsHandler struct {
	R         requests.GetCreateAtomicSwapBidRequests
	Base      getRequestListBaseHandler
	RequestsQ history2.ReviewableRequestsQ
	BalancesQ core2.BalancesQ
	AssetsQ   core2.AssetsQ
	Log       *logan.Entry
}

func (h *getCreateAtomicSwapBidRequestsHandler) MakeAll(w http.ResponseWriter, request requests.GetCreateAtomicSwapBidRequests) error {
	q := h.RequestsQ.FilterByRequestType(uint64(xdr.ReviewableRequestTypeCreateAtomicSwapBid))

	if request.ShouldFilter(requests.FilterTypeCreateAtomicSwapBidRequestsBalance) {
		q = q.FilterByCreateAtomicSwapBidBalance(request.Filters.BaseBalance)
	}

	return h.Base.SelectAndRender(w, *request.GetRequestsBase, q, h.RenderRecord)
}

func (h *getCreateAtomicSwapBidRequestsHandler) RenderRecord(included *regources.Included, record history2.ReviewableRequest) (regources.ReviewableRequest, error) {
	resource := h.Base.PopulateResource(*h.R.GetRequestsBase, included, record)

	if h.R.ShouldInclude(requests.IncludeTypeCreateAtomicSwapBidRequestsBalance) {
		balance, err := h.BalancesQ.GetByAddress(record.Details.CreateAtomicSwapBid.BaseBalance)
		if err != nil {
			return regources.ReviewableRequest{}, errors.Wrap(err, "failed to get balance")
		}

		if balance == nil {
			return regources.ReviewableRequest{}, errors.New("balance not found")
		}
		resource := resources.NewBalance(balance)
		included.Add(resource)
	}

	if h.R.ShouldInclude(requests.IncludeTypeCreateAtomicSwapBidRequestsQuoteAssets) {
		codes := make([]string, 0, len(record.Details.CreateAtomicSwapBid.QuoteAssets))
		for _, v := range record.Details.CreateAtomicSwapBid.QuoteAssets {
			codes = append(codes, v.Asset)
		}
		assets, err := h.AssetsQ.FilterByCodes(codes).Select()
		if err != nil {
			return regources.ReviewableRequest{}, errors.Wrap(err, "failed to get assets")
		}

		if assets == nil {
			return regources.ReviewableRequest{}, errors.New("assets not found")
		}
		for _, record := range assets {
			asset := resources.NewAsset(record)
			included.Add(&asset)
		}
	}

	return resource, nil
}
