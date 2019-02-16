package handlers

import (
	"net/http"

	"gitlab.com/tokend/horizon/fees/fees2"

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
		AccountsQ:          core2.NewAccountsQ(coreRepo),
		BalancesQ:          core2.NewBalancesQ(coreRepo),
		FeesQ:     core2.NewFeesQ(coreRepo),
		LimitsV2Q:          core2.NewLimitsQ(coreRepo),
		ExternalSystemIDsQ: core2.NewExternalSystemIDsQ(coreRepo),
		Log:                ctx.Log(r),
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
		requests.IncludeTypeAccountFees,
		requests.IncludeTypeAccountLimits,
		requests.IncludeTypeAccountExternalSystemIDs,
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
	AccountsQ          core2.AccountsQ
	BalancesQ          core2.BalancesQ
	LimitsV2Q          core2.LimitsQ
	FeesQ     core2.FeesQ
	ExternalSystemIDsQ core2.ExternalSystemIDsQ
	Log                *logan.Entry
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

	response.Data.Relationships.Role, err = h.getRole(request, &response.Included)
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

	response.Data.Relationships.Fees, err = h.getFees(request, &response.Included, account)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get fees for account")
	}
	
	response.Data.Relationships.Limits, err = h.getLimits(request, &response.Included)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get limits")
	}

	response.Data.Relationships.ExternalSystemIDs, err = h.getExternalSystemIDs(request, &response.Included)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get external system IDs")
	}
	return &response, nil
}

func (h *getAccountHandler) getLimits(request *requests.GetAccount, includes *regources.Included) (*regources.RelationCollection, error) {
	limitsQ := h.LimitsV2Q.FilterByAccountID(request.Address)

	coreLimits, err := limitsQ.Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to select limits for account")
	}

	result := &regources.RelationCollection{
		Data: make([]regources.Key, 0, len(coreLimits)),
	}

	for _, coreLimitsUnit := range coreLimits {
		limitsUnit := resources.NewLimits(coreLimitsUnit)
		result.Data = append(result.Data, limitsUnit.Key)

		if request.ShouldInclude(requests.IncludeTypeAccountLimits) {
			includes.Add(limitsUnit)
		}
	}

	return result, nil
}

func (h *getAccountHandler) getExternalSystemIDs(request *requests.GetAccount, includes *regources.Included) (*regources.RelationCollection, error) {
	externalSystemIDsQ := h.ExternalSystemIDsQ.FilterByAccount(request.Address)

	coreExternalSystemIDs, err := externalSystemIDsQ.Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to select external system IDs for account")
	}

	result := regources.RelationCollection{
		Data: make([]regources.Key, 0, len(coreExternalSystemIDs)),
	}

	for _, coreExtSysIDUnit := range coreExternalSystemIDs {
		externalSystemID := resources.NewExternalSystemID(coreExtSysIDUnit)

		result.Data = append(result.Data, externalSystemID.Key)

		if request.ShouldInclude(requests.IncludeTypeAccountExternalSystemIDs) {
			includes.Add(externalSystemID)
		}
	}
	return &result, nil
}

func (h *getAccountHandler) getReferrer(account *core2.Account, request *requests.GetAccount, includes *regources.Included) (*regources.Relation, error) {
	if account.Referrer == "" {
		return nil, nil
	}

	if !request.ShouldInclude(requests.IncludeTypeAccountAccountReferrer) {
		result := resources.NewAccountKey(account.Referrer)

		return result.AsRelation(), nil
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

func (h *getAccountHandler) getRole(request *requests.GetAccount, includes *regources.Included) (*regources.Relation, error) {
	if !request.ShouldInclude(requests.IncludeTypeAccountRole) {
		role := resources.NewRoleKey(request.Address)
		return role.AsRelation(), nil
	}

	role := resources.NewRole(request.Address)

	if request.ShouldInclude(requests.IncludeTypeAccountRoleRules) {
		rules := []regources.Rule{
			resources.NewRule(),
		}

		role.Relationships.Rules = &regources.RelationCollection{
			Data: make([]regources.Key, 0, len(rules)),
		}

		for _, rule := range rules {
			role.Relationships.Rules.Data = append(role.Relationships.Rules.Data, rule.Key)
			includes.Add(&rule)
		}
	} else {
		rulesKeys := []regources.Key{
			resources.NewRuleKey(),
		}
		role.Relationships.Rules = &regources.RelationCollection{
			Data: rulesKeys,
		}
	}

	includes.Add(&role)
	return role.AsRelation(), nil
}

func (h *getAccountHandler) getFees(request *requests.GetAccount, includes *regources.Included, account *core2.Account) (*regources.RelationCollection, error) {

	coreBalances, err := h.BalancesQ.FilterByAccount(request.Address).Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to select balances for account")
	}

	assets := make([]string, 0, len(coreBalances))
	for _, coreBalance := range coreBalances {
		assets = append(assets, coreBalance.AssetCode)
	}

	forAccount, err := h.FeesQ.FilterByAddress(request.Address).Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get fees for account")
	}

	forAccountRole, err := h.FeesQ.FilterByAccountType(account.AccountType).Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get fees for account role")
	}

	generalFees, err := h.FeesQ.FilterByAccountType(core2.GlobalAccountRole).Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get global fees")
	}

	sft := fees2.NewSmartFeeTable(forAccount)
	sft.Update(forAccountRole)
	sft.Update(generalFees)
	sft.AddZeroFees(assets)

	result := regources.RelationCollection{
		Data: make([]regources.Key, 0),
	}
	for _, v := range sft {
		for _, feeRecord := range v {
			hash := resources.CalculateFeeHash(feeRecord.Fee)
			feeKey := resources.NewFeeKey(hash)
			result.Data = append(result.Data, feeKey)

			if request.ShouldInclude(requests.IncludeTypeAccountFees) {
				fee := resources.NewFee(feeRecord.Fee)
				includes.Add(&fee)
			}
		}
	}

	return &result, nil
}
