package handlers

import (
	"net/http"

	"gitlab.com/tokend/horizon/web_v2/resources"
	regources "gitlab.com/tokend/regources/generated"

	"gitlab.com/distributed_lab/logan/v3/errors"

	"gitlab.com/tokend/go/xdr"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/requests"
)

func GetUpdateSaleDetailsRequests(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetUpdateSaleDetailsRequests(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	historyRepo := ctx.HistoryRepo(r)
	handler := getUpdateSaleDetailsRequestsHandler{
		R:         request,
		RequestsQ: history2.NewReviewableRequestsQ(historyRepo),
		SalesQ:    history2.NewSalesQ(historyRepo),
		Log:       ctx.Log(r),
	}

	if !isAllowed(r, w, request.GetRequestsBase.Filters.Requestor[0], request.GetRequestsBase.Filters.Reviewer[0]) {
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

type getUpdateSaleDetailsRequestsHandler struct {
	R         requests.GetUpdateSaleDetailsRequests
	Base      getRequestListBaseHandler
	RequestsQ history2.ReviewableRequestsQ
	SalesQ    history2.SalesQ
	Log       *logan.Entry
}

func (h *getUpdateSaleDetailsRequestsHandler) MakeAll(w http.ResponseWriter, request requests.GetUpdateSaleDetailsRequests) error {
	q := h.RequestsQ.FilterByRequestType(uint64(xdr.ReviewableRequestTypeUpdateSaleDetails))

	return h.Base.SelectAndRender(w, *request.GetRequestsBase, q, h.RenderRecord)
}

func (h *getUpdateSaleDetailsRequestsHandler) RenderRecord(included *regources.Included, record history2.ReviewableRequest) (regources.ReviewableRequest, error) {
	resource := h.Base.PopulateResource(*h.R.GetRequestsBase, included, record)

	if h.R.ShouldInclude(requests.IncludeTypeUpdateSaleDetailsRequestsSale) {
		record, err := h.SalesQ.GetByID(record.Details.UpdateSaleDetails.SaleID)
		if err != nil {
			return regources.ReviewableRequest{}, errors.Wrap(err, "failed to get sale")
		}
		if record == nil {
			return regources.ReviewableRequest{}, errors.New("sale not found")
		}

		sale := resources.NewSale(*record)
		included.Add(&sale)
	}
	return resource, nil
}
