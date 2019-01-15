package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/requests"
	"gitlab.com/tokend/horizon/web_v2/resource"
)

// GetAccount - processes request to get account and it's details by address
func GetAccount(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)
	handler := getAccountHandler{
		AccountsQ: core2.NewAccountsQ(coreRepo),
		BalancesQ: core2.NewBalancesQ(coreRepo),
		Log:       ctx.Log(r),
	}

	request, err := requests.NewGetAccount(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	// TODO: MUST be changes to role based access
	if request.NeedBalanceState() {
		if !isAllowed(r, w, request.Address) {
			return
		}
	}

	result, err := handler.GetAccount(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get account", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if result == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	ape.Render(w, result)

}

type getAccountHandler struct {
	AccountsQ *core2.AccountsQ
	BalancesQ *core2.BalancesQ
	Log       *logan.Entry
}

//GetAccount - returns Account resource
func (h *getAccountHandler) GetAccount(request *requests.GetAccount) (*resource.Account, error) {
	account, err := h.AccountsQ.GetByAddress(request.Address)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get account by address")
	}

	if account == nil {
		return nil, nil
	}

	response := resource.NewAccount(account)
	if request.NeedBalance() {
		balancesQ := h.BalancesQ.FilterByAccount(request.Address)
		if request.NeedBalanceWithAsset() {
			balancesQ = balancesQ.WithAsset()
		}

		var balances []core2.Balance
		balances, err = balancesQ.Select()
		if err != nil {
			return nil, errors.Wrap(err, "failed to select balances for account")
		}

		for i := range balances {
			responseBalance := resource.NewBalance(&balances[i])
			if request.NeedBalanceState() {
				responseBalance.State = resource.NewBalanceState(&balances[i])
			}

			response.Balances = append(response.Balances, responseBalance)
		}
	}

	return &response, nil
}
