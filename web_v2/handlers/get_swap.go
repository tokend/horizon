package handlers

import (
	"net/http"

	"gitlab.com/tokend/horizon/db2/core2"

	regources "gitlab.com/tokend/regources/generated"

	"gitlab.com/tokend/horizon/db2/history2"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/requests"
	"gitlab.com/tokend/horizon/web_v2/resources"
)

func GetSwap(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetSwap(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	handler := getSwapHandler{
		SwapsQ:    history2.NewSwapsQ(ctx.HistoryRepo(r)),
		AssetsQ:   core2.NewAssetsQ(ctx.CoreRepo(r)),
		BalancesQ: core2.NewBalancesQ(ctx.CoreRepo(r)),
		Log:       ctx.Log(r),
	}

	result, err := handler.getSwap(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get swap", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if result == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if request.ShouldIncludeAny(
		requests.IncludeTypeSwapSourceBalance,
		requests.IncludeTypeSwapDestinationBalance,
	) {
		if !isAllowed(r, w, result.Data.Relationships.Source.Data.ID, result.Data.Relationships.Destination.Data.ID) {
			return
		}
	}

	ape.Render(w, result)
}

type getSwapHandler struct {
	SwapsQ    history2.SwapsQ
	AssetsQ   core2.AssetsQ
	BalancesQ core2.BalancesQ
	Log       *logan.Entry
}

// GetSale returns sale with related resources
func (h *getSwapHandler) getSwap(request *requests.GetSwap) (*regources.SwapResponse, error) {

	record, err := h.SwapsQ.GetByID(request.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get swap")
	}

	if record == nil {
		return nil, nil
	}

	resource := resources.NewSwap(*record)
	response := &regources.SwapResponse{
		Data: resource,
	}
	if request.ShouldInclude(requests.IncludeTypeSwapAsset) {
		histAsset, err := h.AssetsQ.GetByCode(record.Asset)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get asset for swap")
		}
		if histAsset == nil {
			return nil, errors.New("Expected asset to exist")
		}
		asset := resources.NewAsset(*histAsset)
		response.Included.Add(&asset)
	}

	if request.ShouldInclude(requests.IncludeTypeSwapSourceBalance) {
		histBalance, err := h.BalancesQ.GetByAddress(record.SourceBalance)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get balance for swap")
		}
		if histBalance == nil {
			return nil, errors.New("Expected balance to exist")
		}
		balance := resources.NewBalance(histBalance)
		response.Included.Add(balance)
	}

	if request.ShouldInclude(requests.IncludeTypeSwapDestinationBalance) {
		histBalance, err := h.BalancesQ.GetByAddress(record.SourceBalance)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get balance for swap")
		}
		if histBalance == nil {
			return nil, errors.New("Expected balance to exist")
		}
		balance := resources.NewBalance(histBalance)
		response.Included.Add(balance)
	}

	return response, nil
}
