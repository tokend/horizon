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

// GetSwapList - processes request to get the list of sales
func GetSwapList(w http.ResponseWriter, r *http.Request) {
	historyRepo := ctx.HistoryRepo(r)

	handler := getSwapListHandler{
		SwapsQ:    history2.NewSwapsQ(historyRepo),
		AssetsQ:   core2.NewAssetsQ(ctx.CoreRepo(r)),
		BalancesQ: core2.NewBalancesQ(ctx.CoreRepo(r)),
		Log:       ctx.Log(r),
	}

	request, err := requests.NewGetSwapList(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.GetSwapList(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get swap list ", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

type getSwapListHandler struct {
	SwapsQ    history2.SwapsQ
	AssetsQ   core2.AssetsQ
	BalancesQ core2.BalancesQ
	Log       *logan.Entry
}

// GetSwapList returns the list of assets with related resources
func (h *getSwapListHandler) GetSwapList(request *requests.GetSwapList) (*regources.SwapListResponse, error) {
	q := h.SwapsQ

	if request.ShouldFilter(requests.FilterTypeSwapListAsset) {
		q = q.FilterByAsset(request.Filters.Asset)
	}

	if request.ShouldFilter(requests.FilterTypeSwapListDestination) {
		q = q.FilterByDestination(request.Filters.Destination)
	}

	if request.ShouldFilter(requests.FilterTypeSwapListSource) {
		q = q.FilterBySource(request.Filters.Source)
	}

	if request.ShouldFilter(requests.FilterTypeSwapListSourceBalance) {
		q = q.FilterBySourceBalance(request.Filters.SourceBalance)
	}

	if request.ShouldFilter(requests.FilterTypeSwapListDestinationBalance) {
		q = q.FilterByDestinationBalance(request.Filters.DestinationBalance)
	}

	if request.ShouldFilter(requests.FilterTypeSwapListState) {
		q = q.FilterByState(request.Filters.State)
	}

	historySwaps, err := q.Page(*request.PageParams).Select()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get swap list ")
	}

	response := &regources.SwapListResponse{
		Data: make([]regources.Swap, 0, len(historySwaps)),
	}

	for _, historySwap := range historySwaps {
		swap := resources.NewSwap(historySwap)
		response.Data = append(response.Data, swap)

		if request.ShouldInclude(requests.IncludeTypeSwapAsset) {
			histAsset, err := h.AssetsQ.GetByCode(historySwap.Asset)
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
			histBalance, err := h.BalancesQ.GetByAddress(historySwap.SourceBalance)
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
			histBalance, err := h.BalancesQ.GetByAddress(historySwap.SourceBalance)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get balance for swap")
			}
			if histBalance == nil {
				return nil, errors.New("Expected balance to exist")
			}
			balance := resources.NewBalance(histBalance)
			response.Included.Add(balance)
		}
	}
	h.PopulateLinks(response, request)

	return response, nil
}

func (h *getSwapListHandler) PopulateLinks(
	response *regources.SwapListResponse, request *requests.GetSwapList,
) {
	if len(response.Data) > 0 {
		response.Links = request.GetCursorLinks(*request.PageParams, response.Data[len(response.Data)-1].ID)
	} else {
		response.Links = request.GetCursorLinks(*request.PageParams, "")
	}
}
