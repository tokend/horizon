package handlers

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/requests"
	"gitlab.com/tokend/horizon/web_v2/resources"
	regources "gitlab.com/tokend/regources/generated"
)

type getHistory struct {
	// cached filter entries
	Account *history2.Account
	Balance *history2.Balance

	EffectsQ  history2.ParticipantEffectsQ
	AccountsQ history2.AccountsQ
	BalanceQ  history2.BalancesQ
	Log       *logan.Entry
}

func newHistoryHandler(r *http.Request) getHistory {
	historyRepo := ctx.HistoryRepo(r)
	handler := getHistory{
		EffectsQ:  history2.NewParticipantEffectsQ(historyRepo),
		AccountsQ: history2.NewAccountsQ(historyRepo),
		BalanceQ:  history2.NewBalancesQ(historyRepo),
		Log:       ctx.Log(r),
	}

	return handler

}

func (h *getHistory) prepare(w http.ResponseWriter, r *http.Request) (*requests.GetHistory, bool) {
	request, err := requests.NewGetHistory(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return nil, false
	}

	// TODO: need to refactor
	if request.Filters.Account != "" {
		h.Account, err = h.AccountsQ.ByAddress(request.Filters.Account)
		if err != nil {
			ctx.Log(r).WithError(err).Error("failed to account", logan.F{
				"account_address": request.Filters.Account,
			})
			ape.RenderErr(w, problems.InternalError())
			return nil, false
		}

		if h.Account == nil {
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"filter[account]": errors.New("not found"),
			})...)
			return nil, false
		}
	}

	if request.Filters.Balance != "" {
		h.Balance, err = h.BalanceQ.GetByAddress(request.Filters.Balance)
		if err != nil {
			ctx.Log(r).WithError(err).Error("failed to account", logan.F{
				"balance_address": request.Filters.Balance,
			})
			ape.RenderErr(w, problems.InternalError())
			return nil, false
		}

		if h.Balance == nil {
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"filter[balance]": errors.New("not found"),
			})...)
			return nil, false
		}
	}

	if !h.ensureAllowed(w, r) {
		return nil, false
	}

	return request, true
}

func (h *getHistory) ApplyFilters(request *requests.GetHistory,
	q history2.ParticipantEffectsQ) history2.ParticipantEffectsQ {
	q = q.WithAccount().WithBalance().Page(*request.PageParams)
	if request.ShouldInclude(requests.IncludeTypeHistoryOperation) {
		q = q.WithOperation()
	}

	if h.Account != nil {
		q = q.ForAccount(h.Account.ID)
	}

	if h.Balance != nil {
		q = q.ForBalance(h.Balance.ID)
	}

	return q
}

func (h *getHistory) SelectAndPopulate(
	request *requests.GetHistory,
	effectsQ history2.ParticipantEffectsQ,
) (regources.ParticipantsEffectsResponse, error) {

	result := regources.ParticipantsEffectsResponse{
		Data: []regources.ParticipantsEffect{},
	}

	effects, err := effectsQ.Select()
	if err != nil {
		return result, errors.Wrap(err, "failed to load participant effects")
	}

	if len(effects) == 0 {
		result.Links = request.GetCursorLinks(*request.PageParams, "")
		return result, nil
	}

	result.Data = make([]regources.ParticipantsEffect, 0, len(effects))
	for i := range effects {
		effect := getEffect(effects[i])
		if request.ShouldInclude(requests.IncludeTypeHistoryOperation) {
			op := resources.NewOperation(*effects[i].Operation)

			opDetails := resources.NewOperationDetails(*effects[i].Operation)
			op.Relationships.Details = opDetails.GetKey().AsRelation()

			if request.ShouldInclude(requests.IncludeTypeHistoryOperationDetails) {
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

func getEffect(effect history2.ParticipantEffect) regources.ParticipantsEffect {
	var balance *regources.Relation
	if effect.BalanceAddress != nil {
		balance = resources.NewBalanceKey(*effect.BalanceAddress).AsRelation()
	}

	var asset *regources.Relation
	if effect.AssetCode != nil {
		asset = resources.NewAssetKey(*effect.AssetCode).AsRelation()
	}

	return regources.ParticipantsEffect{
		Key: resources.NewParticipantEffectKey(effect.ID),
		Relationships: regources.ParticipantsEffectRelationships{
			Account:   resources.NewAccountKey(effect.AccountAddress).AsRelation(),
			Balance:   balance,
			Asset:     asset,
			Operation: resources.NewOperationKey(effect.OperationID).AsRelation(),
		},
	}

}

// ensure allowed - checks it requester is allowed to access the data. If not it renders error and returns false.
func (h *getHistory) ensureAllowed(w http.ResponseWriter, httpRequest *http.Request) bool {
	if h.Account != nil {
		return isAllowed(httpRequest, w, h.Account.Address)
	}

	if h.Balance != nil {
		account, err := h.tryGetAccountForBalance(h.Balance)
		if err != nil {
			ctx.Log(httpRequest).WithError(err).Error("failed to load account for balance")
			ape.RenderErr(w, problems.InternalError())
			return false
		}

		return isAllowed(httpRequest, w, account.Address)
	}

	return isAllowed(httpRequest, w, ctx.CoreInfo(httpRequest).AdminAccountID)
}

func (h *getHistory) tryGetAccountForBalance(balance *history2.Balance) (history2.Account, error) {
	// TODO: fuck normalization
	account, err := h.AccountsQ.ByID(balance.AccountID)
	if err != nil {
		return history2.Account{}, errors.Wrap(err, "failed to load account by ID")
	}

	if account == nil {
		return history2.Account{}, errors.From(errors.New("found balance, but failed to find account for it"), logan.F{
			"balance": balance.ID,
		})
	}

	return *account, nil
}