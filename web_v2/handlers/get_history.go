package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/requests"
	"gitlab.com/tokend/horizon/web_v2/resources"
	"gitlab.com/tokend/regources/v2"
)

// GetHistory - processes request to get the list of participant effects
func GetHistory(w http.ResponseWriter, r *http.Request) {
	historyRepo := ctx.HistoryRepo(r)
	handler := getHistory{
		EffectsQ:  history2.NewParticipantEffectsQ(historyRepo),
		AccountsQ: history2.NewAccountsQ(historyRepo),
		BalanceQ:  history2.NewBalancesQ(historyRepo),
		Log:       ctx.Log(r),
	}

	request, err := requests.NewGetHistory(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	//if !handler.ensureAllowed(w, r, request) {
	//	return
	//}

	result, err := handler.GetHistory(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get asset list", logan.F{})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

type getHistory struct {
	EffectsQ  history2.ParticipantEffectsQ
	AccountsQ history2.AccountsQ
	BalanceQ  history2.BalancesQ
	Log       *logan.Entry
}

// GetHistory returns the list of participant effects with related resources
func (h *getHistory) GetHistory(request *requests.GetHistory) (regources.ParticipantEffectsResponse, error) {
	result := regources.ParticipantEffectsResponse{
		Data: []regources.ParticipantEffect{},
	}
	q := h.EffectsQ.WithAccount().WithBalance().Page(*request.PageParams)
	if request.ShouldInclude(requests.IncludeTypeHistoryOperation) {
		q = q.WithOperation()
	}

	if request.Filters.Account != "" {
		account, err := h.AccountsQ.ByAddress(request.Filters.Account)
		if err != nil {
			return result, errors.Wrap(err, "failed to get account by id")
		}

		// this could happen due to invalid account address, partial history or account does not exists
		// in any case we do not care and return empty page
		if account == nil {
			return result, nil
		}

		q = q.ForAccount(account.ID)
	}

	if request.Filters.Balance != "" {
		balance, err := h.BalanceQ.GetByAddress(request.Filters.Balance)
		if err != nil {
			return result, errors.Wrap(err, "failed to load balance by id")
		}

		// this could happen due to invalid balance address, partial history or balance does not exists
		// in any case we do not care and return empty page
		if balance == nil {
			return result, nil
		}

		q = q.ForBalance(balance.ID)
	}

	effects, err := q.Select()
	if err != nil {
		return result, errors.Wrap(err, "failed to load participant effects")
	}

	if len(effects) == 0 {
		result.Links = request.GetCursorLinks(*request.PageParams, "")
		return result, nil
	}

	result.Data = make([]regources.ParticipantEffect, 0, len(effects))
	for i := range effects {
		effect := getEffect(effects[i])
		if request.ShouldInclude(requests.IncludeTypeHistoryOperation) {
			op := resources.NewOperation(*effects[i].Operation)

			if request.ShouldInclude(requests.IncludeTypeHistoryOperationDetails) {
				opDetails := resources.NewOperationDetails(*effects[i].Operation)
				op.Relationships.Details = opDetails.GetKey().AsRelation()
				result.Included.Add(opDetails)
			}

			result.Included.Add(&op)
		}

		if effects[i].Effect != nil {
			change := resources.NewEffect(effects[i].ID, *effects[i].Effect)
			effect.Relationships.Effect = change.GetKey().AsRelation()
			if request.ShouldInclude(requests.IncludeTypeHistoryEffect) {
				result.Included.Add(change)
			}
		}

		result.Data = append(result.Data, effect)
	}

	result.Links = request.GetCursorLinks(*request.PageParams, result.Data[len(result.Data)-1].ID)

	return result, nil
}

func getEffect(effect history2.ParticipantEffect) regources.ParticipantEffect {
	var balance *regources.Relation
	if effect.BalanceAddress != nil {
		balance = resources.NewBalanceKey(*effect.BalanceAddress).AsRelation()
	}

	var asset *regources.Relation
	if effect.AssetCode != nil {
		asset = resources.NewAssetKey(*effect.AssetCode).AsRelation()
	}

	return regources.ParticipantEffect{
		Key: resources.NewParticipantEffectKey(effect.ID),
		Relationships: regources.ParticipantEffectRelation{
			Account:   resources.NewAccountKey(effect.AccountAddress).AsRelation(),
			Balance:   balance,
			Asset:     asset,
			Operation: resources.NewOperationKey(effect.OperationID).AsRelation(),
		},
	}

}

// ensure allowed - checks it requester is allowed to access the data. If not it renders error and returns false.
func (h *getHistory) ensureAllowed(w http.ResponseWriter, httpRequest *http.Request, request *requests.GetHistory) bool {
	if request.Filters.Account != "" {
		return isAllowed(httpRequest, w, request.Filters.Account)
	}

	if request.Filters.Balance != "" {
		account, err := h.tryGetAccountForBalance(request.Filters.Balance)
		if err != nil {
			ctx.Log(httpRequest).WithError(err).Error("failed to load account for balance")
			ape.RenderErr(w, problems.InternalError())
			return false
		}

		// if we failed to find account by balance address - balance does not exists, so we can render empty page
		if account == nil {
			ape.Render(w, regources.ParticipantEffectsResponse{})
			return false
		}

		return isAllowed(httpRequest, w, account.Address)
	}

	return isAllowed(httpRequest, w, ctx.CoreInfo(httpRequest).MasterAccountID)
}

func (h *getHistory) tryGetAccountForBalance(balanceAddress string) (*history2.Account, error) {
	balance, err := h.BalanceQ.GetByAddress(balanceAddress)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load balance by address")
	}

	if balance == nil {
		return nil, nil
	}

	// TODO: fuck normalization
	account, err := h.AccountsQ.ByID(balance.AccountID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load account by ID")
	}

	if account == nil {
		return nil, errors.From(errors.New("found balance, but failed to find account for it"), logan.F{
			"balance": balance.ID,
		})
	}

	return account, nil
}
