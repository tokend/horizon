package handlers

import (
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/requests"
	"gitlab.com/tokend/horizon/web_v2/resources"
	regources "gitlab.com/tokend/regources/generated"
	"net/http"
)

func GetParticipantEffect(w http.ResponseWriter, r *http.Request) {
	handler := newParticipantEffectHandler(r)

	request, err := requests.NewGetParticipantEffect(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	result, err := handler.GetParticipantEffect(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get participant effect list", logan.F{})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

type getParticipantEffectHandler struct {
	AssetsQ   core2.AssetsQ
	EffectsQ  history2.ParticipantEffectsQ
	AccountsQ history2.AccountsQ
	BalanceQ  history2.BalancesQ
	Log       *logan.Entry
}

func newParticipantEffectHandler(r *http.Request) getParticipantEffectHandler {
	historyRepo := ctx.HistoryRepo(r)
	handler := getParticipantEffectHandler{
		AssetsQ:   core2.NewAssetsQ(ctx.CoreRepo(r)),
		EffectsQ:  history2.NewParticipantEffectsQ(historyRepo),
		AccountsQ: history2.NewAccountsQ(historyRepo),
		BalanceQ:  history2.NewBalancesQ(historyRepo),
		Log:       ctx.Log(r),
	}

	return handler
}

// GetParicipantEffect returns the participant effect with related resources by id
func (h *getParticipantEffectHandler) GetParticipantEffect(request *requests.GetParticipantEffect) (regources.ParticipantsEffectResponse, error) {

	result := regources.ParticipantsEffectResponse{}

	effect, err := h.EffectsQ.FilterByID(request.ID).Get()
	if err != nil {
		return result, errors.Wrap(err, "failed to load participant effects")
	}

	resEffect := getEffect(*effect)

	if request.ShouldInclude(requests.IncludeTypeHistoryOperation) {
		op := resources.NewOperation(*effect.Operation)

		opDetails := resources.NewOperationDetails(*effect.Operation)
		op.Relationships.Details = opDetails.GetKey().AsRelation()

		if request.ShouldInclude(requests.IncludeTypeHistoryOperationDetails) {
			result.Included.Add(opDetails)
		}

		result.Included.Add(&op)
	}

	if effect.Effect != nil {
		change := resources.NewEffect(effect.ID, *effect.Effect)
		resEffect.Relationships.Effect = change.GetKey().AsRelation()
		if request.ShouldInclude(requests.IncludeTypeHistoryEffect) {
			result.Included.Add(change)
		}
	}

	if effect.AssetCode != nil {
		if request.ShouldInclude(requests.IncludeTypeHistoryAsset) {
			rawAsset, err := h.AssetsQ.GetByCode(*effect.AssetCode)
			if err != nil {
				return result, errors.Wrap(err, "failed to load asset")
			}
			asset := resources.NewAsset(*rawAsset)
			result.Included.Add(&asset)
		}
	}

	return result, nil
}
