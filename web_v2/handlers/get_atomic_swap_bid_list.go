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

// GetAssetList - processes request to get the list of assets
func GetAtomicSwapBidList(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)
	handler := getAtomicSwapBidListHandler{
		AtomicSwapBidQ:        core2.NewAtomicSwapBidQ(coreRepo),
		AtomicSwapQuoteAssetQ: core2.NewAtomicSwapQuoteAssetQ(coreRepo),
		BalanceQ:              core2.NewBalancesQ(coreRepo),
		AssetsQ:               core2.NewAssetsQ(coreRepo),
		Log:                   ctx.Log(r),
	}

	request, err := requests.NewGetAtomicSwapBidList(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.GetAtomicSwapBidList(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get asset list", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

type getAtomicSwapBidListHandler struct {
	AssetsQ               core2.AssetsQ
	BalanceQ              core2.BalancesQ
	AtomicSwapBidQ        core2.AtomicSwapBidQ
	AtomicSwapQuoteAssetQ core2.AtomicSwapQuoteAssetQ
	Log                   *logan.Entry
}

// GetAssetList returns the list of assets with related resources
func (h *getAtomicSwapBidListHandler) GetAtomicSwapBidList(request *requests.GetAtomicSwapBidList,
) (*regources.AtomicSwapBidsResponse, error) {
	q := h.AtomicSwapBidQ.Page(*request.PageParams)
	if request.ShouldFilter(requests.FilterTypeBidListOwner) {
		q = q.FilterByOwner(request.Filters.Owner)
	}
	if request.ShouldFilter(requests.FilterTypeBidListBaseAsset) {
		q = q.FilterByBaseAssets([]string{request.Filters.BaseAsset})
	}
	if request.ShouldFilter(requests.FilterTypeBidListQuoteAssets) {
		q = q.FilterByQuoteAssets(request.Filters.QuoteAssets)
	}
	bids, err := q.Select()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get asset list")
	}

	response := &regources.AtomicSwapBidsResponse{
		Data:  make([]regources.AtomicSwapBid, 0, len(bids)),
		Links: request.GetOffsetLinks(*request.PageParams),
	}

	for _, bid := range bids {
		data := resources.NewAtomicSwapBid(bid)
		bidOwner := resources.NewAccountKey(bid.OwnerID)
		data.Relationships.Owner = bidOwner.AsRelation()

		if request.ShouldInclude(requests.IncludeTypeBidOwner) {
			response.Included.Add(&bidOwner)
		}

		baseBalanceKey := resources.NewBalanceKey(bid.BaseBalanceID)
		data.Relationships.BaseBalance = baseBalanceKey.AsRelation()

		if request.ShouldInclude(requests.IncludeTypeBidBaseBalance) {
			baseBalance := regources.Balance{
				Key: baseBalanceKey,
				Relationships: regources.BalanceRelation{
					Asset: resources.NewAssetKey(bid.BaseAsset).AsRelation(),
				},
			}
			response.Included.Add(&baseBalance)
		}

		if request.ShouldInclude(requests.IncludeTypeBidBaseAsset) {
			assetRaw, err := h.AssetsQ.GetByCode(bid.BaseAsset)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get base asset")
			}
			if assetRaw == nil {
				return nil, errors.From(errors.New("base asset not found"), logan.F{
					"asset code": bid.BaseAsset,
				})
			}

			asset := resources.NewAsset(*assetRaw)
			response.Included.Add(&asset)
		}

		quoteAssetsRaw, err := h.AtomicSwapQuoteAssetQ.FilterByIDs([]int64{bid.BidID}).Select()
		if err != nil {
			return nil, errors.Wrap(err, "failed to select bid quote assets")
		}
		if quoteAssetsRaw == nil {
			return nil, errors.From(errors.New("expected bid quote assets to exists"), logan.F{
				"bid_id": bid.BidID,
			})
		}

		quoteAssets := &regources.RelationCollection{
			Data: make([]regources.Key, 0, len(quoteAssetsRaw)),
		}

		for _, quoteAssetRaw := range quoteAssetsRaw {
			quoteAsset := resources.NewAtomicSwapBidQuoteAsset(quoteAssetRaw)
			quoteAssets.Data = append(quoteAssets.Data, quoteAsset.Key)

			if request.ShouldInclude(requests.IncludeTypeBidListQuoteAssets) {
				response.Included.Add(&quoteAsset)
			}
		}
		data.Relationships.QuoteAssets = quoteAssets

		response.Data = append(response.Data, data)
	}

	return response, nil
}
