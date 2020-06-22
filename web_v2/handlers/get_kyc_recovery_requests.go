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
	regources "gitlab.com/tokend/regources/generated"
)

func GetKYCRecoveryRequests(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetKYCRecoveryRequests(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	historyRepo := ctx.HistoryRepo(r)
	coreRepo := ctx.CoreRepo(r)
	handler := getKYCRecoveryRequestsHandler{
		R:         request,
		RequestsQ: history2.NewReviewableRequestsQ(historyRepo),
		AccountsQ: core2.NewAccountsQ(coreRepo),
		Log:       ctx.Log(r),
	}

	if !isAllowed(r, w,
		request.GetRequestsBase.Filters.Requestor[0],
		request.GetRequestsBase.Filters.Reviewer[0],
		request.Filters.Account[0]) {
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

type getKYCRecoveryRequestsHandler struct {
	R         requests.GetKYCRecoveryRequests
	Base      getRequestListBaseHandler
	RequestsQ history2.ReviewableRequestsQ
	AccountsQ core2.AccountsQ
	Log       *logan.Entry
}

func (h *getKYCRecoveryRequestsHandler) MakeAll(w http.ResponseWriter, request requests.GetKYCRecoveryRequests) error {
	q := h.RequestsQ.FilterByRequestType(uint64(xdr.ReviewableRequestTypeKycRecovery))

	if request.ShouldFilter(requests.FilterTypeKYCRecoveryRequestsAccount) {
		q = q.FilterByKYCRecoveryTargetAccount(request.Filters.Account[0])
	}

	return h.Base.SelectAndRender(w, *request.GetRequestsBase, q, h.RenderRecord)
}

func (h *getKYCRecoveryRequestsHandler) RenderRecord(included *regources.Included, record history2.ReviewableRequest) (regources.ReviewableRequest, error) {
	resource := h.Base.PopulateResource(*h.R.GetRequestsBase, included, record)

	if h.R.ShouldInclude(requests.IncludeTypeKYCRecoveryRequestsAccount) {
		account, err := h.AccountsQ.GetByAddress(record.Details.KYCRecovery.TargetAccount)
		if err != nil {
			return regources.ReviewableRequest{}, errors.Wrap(err, "failed to get account")
		}

		if account == nil {
			return regources.ReviewableRequest{}, errors.New("account not found")
		}
		resource := resources.NewAccount(*account, nil)
		included.Add(&resource)
	}

	return resource, nil
}
