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

func GetCreateIssuanceRequests(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetCreateIssuanceRequests(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	historyRepo := ctx.HistoryRepo(r)
	coreRepo := ctx.CoreRepo(r)
	handler := getCreateIssuanceRequestsHandler{
		R:         request,
		RequestsQ: history2.NewReviewableRequestsQ(historyRepo),
		AssetsQ:   core2.NewAssetsQ(coreRepo),
		Log:       ctx.Log(r),
	}

	constraints := []string{
		request.GetRequestsBase.Filters.Requestor,
		request.GetRequestsBase.Filters.Reviewer,
	}

	// receiving balance owner should be able to see issuance requests
	if request.Filters.Receiver != "" {
		balance, err := core2.NewBalancesQ(coreRepo).GetByAddress(request.Filters.Receiver)
		if err != nil {
			ctx.Log(r).
				WithError(err).
				WithFields(logan.F{"receiver": request.Filters.Receiver}).
				Error("failed to get receiver balance")
			ape.RenderErr(w, problems.InternalError())
			return
		}
		if balance != nil {
			constraints = append(constraints, balance.AccountAddress)
		}
	}

	if !isAllowed(r, w, constraints...) {
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

type getCreateIssuanceRequestsHandler struct {
	R         requests.GetCreateIssuanceRequests
	Base      getRequestListBaseHandler
	RequestsQ history2.ReviewableRequestsQ
	AssetsQ   core2.AssetsQ
	Log       *logan.Entry
}

func (h *getCreateIssuanceRequestsHandler) MakeAll(w http.ResponseWriter, request requests.GetCreateIssuanceRequests) error {
	q := h.RequestsQ.FilterByRequestType(uint64(xdr.ReviewableRequestTypeCreateIssuance))

	if request.ShouldFilter(requests.FilterTypeCreateIssuanceRequestsAsset) {
		q = q.FilterByCreateIssuanceAsset(request.Filters.Asset)
	}

	if request.ShouldFilter(requests.FilterTypeCreateIssuanceRequestsReceiver) {
		q = q.FilterByCreateIssuanceReceiver(request.Filters.Receiver)
	}

	return h.Base.SelectAndRender(w, *request.GetRequestsBase, q, h.RenderRecord)
}

func (h *getCreateIssuanceRequestsHandler) RenderRecord(included *regources.Included, record history2.ReviewableRequest) (regources.ReviewableRequest, error) {
	resource := h.Base.PopulateResource(*h.R.GetRequestsBase, included, record)

	if h.R.ShouldInclude(requests.IncludeTypeCreateIssuanceRequestsAsset) {
		asset, err := h.AssetsQ.GetByCode(record.Details.CreateIssuance.Asset)
		if err != nil {
			return regources.ReviewableRequest{}, errors.Wrap(err, "failed to get asset")
		}

		if asset == nil {
			return regources.ReviewableRequest{}, errors.New("asset not found")
		}
		resource := resources.NewAsset(*asset)
		included.Add(&resource)
	}

	return resource, nil
}
