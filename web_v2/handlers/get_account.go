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
	"gitlab.com/tokend/horizon/web_v2/resources"
	"gitlab.com/tokend/regources/v2"
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
	if request.ShouldIncludeAny(
		requests.IncludeTypeAccountBalancesState,
		requests.IncludeTypeAccountAccountReferrer,
		requests.IncludeTypeAccountState,
		requests.IncludeTypeAccountRole,
		requests.IncludeTypeAccountRoleRules,
	) {
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
	AccountsQ core2.AccountsQ
	BalancesQ core2.BalancesQ
	Log       *logan.Entry
}

//GetAccount - returns Account resources
func (h *getAccountHandler) GetAccount(request *requests.GetAccount) (*regources.Account, error) {
	account, err := h.AccountsQ.GetByAddress(request.Address)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get account by address")
	}

	if account == nil {
		return nil, nil
	}

	response := resources.NewAccount(account)
	if request.ShouldInclude(requests.IncludeTypeAccountState) {
		response.State = resources.NewAccountState(account)
	}

	if request.ShouldIncludeAny(
		requests.IncludeTypeAccountRole,
		requests.IncludeTypeAccountRoleRules,
	) {
		response.Role = h.getRole(request)
	}

	if request.ShouldInclude(requests.IncludeTypeAccountBalances) {
		response.Balances, err = h.getBalances(request)
		if err != nil {
			return nil, errors.Wrap(err, "failed to include balances")
		}
	}

	if request.ShouldInclude(requests.IncludeTypeAccountAccountReferrer) {
		response.Referrer, err = h.getReferrer(account)
		if err != nil {
			return nil, errors.Wrap(err, "failed to include referrer")
		}
	}

	return response, nil
}

func (h *getAccountHandler) getReferrer(account *core2.Account) (*regources.Account, error) {
	if account.Referrer == "" {
		return nil, nil
	}

	referrer, err := h.AccountsQ.GetByAddress(account.Referrer)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load referrer", logan.F{
			"referrer_address": account.Referrer,
		})
	}

	if referrer == nil {
		return nil, errors.Wrap(err, "expected referrer to exist", logan.F{
			"account_address":  account.Address,
			"referrer_address": account.Referrer,
		})
	}

	return resources.NewAccount(referrer), nil
}

func (h *getAccountHandler) getBalances(request *requests.GetAccount) ([]*regources.Balance, error) {

	balancesQ := h.BalancesQ.FilterByAccount(request.Address)
	if request.ShouldInclude(requests.IncludeTypeAccountBalancesAsset) {
		balancesQ = balancesQ.WithAsset()
	}

	balances, err := balancesQ.Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to select balances for account")
	}

	result := make([]*regources.Balance, len(balances))
	for i := range balances {
		responseBalance := resources.NewBalance(&balances[i])
		if request.ShouldInclude(requests.IncludeTypeAccountBalancesState) {
			responseBalance.State = resources.NewBalanceState(&balances[i])
		}

		result[i] = responseBalance
	}

	return result, nil
}

func (h *getAccountHandler) getRole(request *requests.GetAccount) *regources.Role {
	result := regources.Role{
		ID: "mocked_role",
		Details: map[string]interface{}{
			"name": "Name of the Mocked Role",
		},
	}

	if !request.ShouldInclude(requests.IncludeTypeAccountRoleRules) {
		return &result
	}

	result.Rules = []*regources.Rule{
		{
			ID:       "mocked_rule_id",
			Resource: "NOTE: format will be changed",
			Action:   "view",
			Details: map[string]interface{}{
				"name": "Name of the mocked Rule",
			},
		},
	}

	return &result
}
