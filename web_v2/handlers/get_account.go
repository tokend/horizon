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
		AccountsQ:    core2.NewAccountsQ(coreRepo),
		BalancesQ:    core2.NewBalancesQ(coreRepo),
		AccountRoleQ: core2.NewAccountRoleQ(coreRepo),
		AccountRuleQ: core2.NewAccountRuleQ(coreRepo),
		Log:          ctx.Log(r),
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
	AccountsQ    core2.AccountsQ
	BalancesQ    core2.BalancesQ
	AccountRoleQ core2.AccountRoleQ
	AccountRuleQ core2.AccountRuleQ
	Log          *logan.Entry
}

//GetAccount - returns Account resources
func (h *getAccountHandler) GetAccount(request *requests.GetAccount) (*regources.AccountResponse, error) {
	account, err := h.AccountsQ.GetByAddress(request.Address)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get account by address")
	}

	if account == nil {
		return nil, nil
	}

	response := regources.AccountResponse{
		Data: resources.NewAccount(*account),
	}

	response.Data.Relationships.Role, err = h.getRole(request, &response.Included, *account)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get role")
	}

	response.Data.Relationships.Balances, err = h.getBalances(request, &response.Included)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get balances")
	}

	response.Data.Relationships.Referrer, err = h.getReferrer(account, request, &response.Included)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get referrer")
	}

	return &response, nil
}

func (h *getAccountHandler) getReferrer(account *core2.Account, request *requests.GetAccount, includes *regources.Included) (*regources.Relation, error) {
	if account.Referrer == nil {
		return nil, nil
	}

	if !request.ShouldInclude(requests.IncludeTypeAccountAccountReferrer) {
		result := resources.NewAccountKey(*account.Referrer)

		return result.AsRelation(), nil
	}

	referrer, err := h.AccountsQ.GetByAddress(*account.Referrer)
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

	result := resources.NewAccount(*referrer)
	includes.Add(&result)

	return result.AsRelation(), nil
}

func (h *getAccountHandler) getBalances(request *requests.GetAccount, includes *regources.Included) (*regources.RelationCollection, error) {
	balancesQ := h.BalancesQ.FilterByAccount(request.Address)

	if request.ShouldInclude(requests.IncludeTypeAccountBalancesAsset) {
		balancesQ = balancesQ.WithAsset()
	}

	coreBalances, err := balancesQ.Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to select balances for account")
	}

	result := regources.RelationCollection{
		Data: make([]regources.Key, 0, len(coreBalances)),
	}

	for i, coreBalance := range coreBalances {
		balance := resources.NewBalance(&coreBalances[i])
		result.Data = append(result.Data, balance.Key)

		if request.ShouldInclude(requests.IncludeTypeAccountBalances) {
			balance.Relationships.State = resources.NewBalanceStateKey(coreBalance.BalanceAddress).AsRelation()
			balance.Relationships.Asset = resources.NewAssetKey(coreBalance.AssetCode).AsRelation()

			includes.Add(balance)
		}

		if request.ShouldInclude(requests.IncludeTypeAccountBalancesState) {
			state := resources.NewBalanceState(&coreBalances[i])
			includes.Add(state)
		}

		if request.ShouldInclude(requests.IncludeTypeAccountBalancesAsset) {
			asset := resources.NewAsset(*coreBalances[i].Asset)
			includes.Add(&asset)
		}
	}

	return &result, nil
}

func (h *getAccountHandler) getRole(request *requests.GetAccount,
	includes *regources.Included, account core2.Account,
) (*regources.Relation, error) {
	if !request.ShouldInclude(requests.IncludeTypeAccountRole) {
		role := resources.NewAccountRoleKey(account.RoleID)
		return role.AsRelation(), nil
	}

	roleRaw, err := h.AccountRoleQ.GetByID(account.RoleID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get account role by id")
	}

	if roleRaw == nil {
		return nil, errors.New("role not found")
	}

	role := resources.NewAccountRole(*roleRaw)
	ruleKeys := []regources.Key(nil)

	if request.ShouldInclude(requests.IncludeTypeAccountRoleRules) {
		rules, err := h.AccountRuleQ.FilterByIDs(roleRaw.RuleIDs...).Select()
		if err != nil {
			return nil, errors.Wrap(err, "failed to select account rules")
		}

		for _, ruleRaw := range rules {
			rule := resources.NewAccountRule(ruleRaw)
			ruleKeys = append(ruleKeys, rule.Key)
			includes.Add(&rule)
		}
	} else {
		for _, ruleID := range roleRaw.RuleIDs {
			ruleKeys = append(ruleKeys, resources.NewAccountRuleKey(ruleID))
		}
	}

	role.Relationships.Rules = &regources.RelationCollection{
		Data: ruleKeys,
	}
	includes.Add(&role)

	return role.AsRelation(), nil
}
