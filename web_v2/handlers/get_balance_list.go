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
	regources "gitlab.com/tokend/regources/generated"
)

// GetBalanceList - processes request to get the list of balances
func GetBalanceList(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)
	handler := getBalanceListHandler{
		AssetsQ:   core2.NewAssetsQ(coreRepo),
		BalancesQ: core2.NewBalancesQ(coreRepo),
		Log:       ctx.Log(r),
	}

	request, err := requests.NewGetBalanceList(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	assetOwner, err := handler.getAssetOwner(request.Filters.Asset)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get asset owner", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if !isAllowed(r, w, assetOwner) {
		return
	}

	result, err := handler.GetBalanceList(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get asset list", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

type getBalanceListHandler struct {
	BalancesQ core2.BalancesQ
	AssetsQ   core2.AssetsQ
	Log       *logan.Entry
}

func (h *getBalanceListHandler) getAssetOwner(assetCode string) (string, error) {
	if assetCode == "" {
		return "", nil
	}

	coreAsset, err := h.AssetsQ.GetByCode(assetCode)
	if err != nil {
		return "", errors.Wrap(err, "Failed to get asset")
	}

	return coreAsset.Owner, nil
}

// GetBalanceList returns list of balances with related resources
func (h *getBalanceListHandler) GetBalanceList(request *requests.GetBalanceList) (*regources.BalancesResponse, error) {
	q := h.BalancesQ.Page(*request.PageParams)
	if request.ShouldFilter(requests.FilterTypeBalanceListAsset) {
		q = q.FilterByAsset(request.Filters.Asset)
	}

	coreBalances, err := q.Select()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get balance list")
	}

	response := &regources.BalancesResponse{
		Data:  make([]regources.Balance, 0, len(coreBalances)),
		Links: request.GetOffsetLinks(*request.PageParams),
	}

	for _, coreBalance := range coreBalances {
		balance := resources.NewBalance(&coreBalance)
		balance.Relationships.Asset = resources.NewAssetKey(coreBalance.AssetCode).AsRelation()
		balance.Relationships.State = resources.NewBalanceStateKey(coreBalance.BalanceAddress).AsRelation()

		if request.ShouldInclude(requests.IncludeTypeBalanceListState) {
			response.Included.Add(resources.NewBalanceState(&coreBalance))
		}

		response.Data = append(response.Data, *balance)
	}

	return response, nil
}
