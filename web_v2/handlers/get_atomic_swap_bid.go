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
	"gitlab.com/tokend/regources/generated"
)

// GetAtomicSwapBid - processes request to get atomic swap bid and it's details by id
func GetAtomicSwapBid(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)
	handler := getAtomicSwapBidHandler{
		AssetsQ:               core2.NewAssetsQ(coreRepo),
		AtomicSwapBidQ:        core2.NewAtomicSwapBidQ(coreRepo),
		AtomicSwapQuoteAssetQ: core2.NewAtomicSwapQuoteAssetQ(coreRepo),
		BalanceQ:              core2.NewBalancesQ(coreRepo),
		Log:                   ctx.Log(r),
	}

	request, err := requests.NewGetAtomicSwapBid(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.GetAtomicSwapBid(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get asset", logan.F{
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

type getAtomicSwapBidHandler struct {
	AssetsQ               core2.AssetsQ
	BalanceQ              core2.BalancesQ
	AtomicSwapBidQ        core2.AtomicSwapBidQ
	AtomicSwapQuoteAssetQ core2.AtomicSwapQuoteAssetQ
	Log                   *logan.Entry
}

// GetAtomicSwapBid returns atomic swap bid with related resources
func (h *getAtomicSwapBidHandler) GetAtomicSwapBid(request *requests.GetAtomicSwapBid,
) (*regources.AtomicSwapBidResponse, error) {
	bid, err := h.AtomicSwapBidQ.GetByID(request.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get atomic swap bid by id")
	}
	if bid == nil {
		return nil, nil
	}

	response := &regources.AtomicSwapBidResponse{
		Data: resources.NewAtomicSwapBid(*bid),
	}

	bidOwner := resources.NewAccountKey(bid.OwnerID)
	response.Data.Relationships.Owner = bidOwner.AsRelation()

	if request.ShouldInclude(requests.IncludeTypeBidOwner) {
		response.Included.Add(&bidOwner)
	}

	baseBalanceKey := resources.NewBalanceKey(bid.BaseBalanceID)
	response.Data.Relationships.BaseBalance = baseBalanceKey.AsRelation()

	if request.ShouldInclude(requests.IncludeTypeBidBaseBalance) {
		baseBalance := regources.Balance{
			Key: baseBalanceKey,
			Relationships: &regources.BalanceRelationships{
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

		if request.ShouldInclude(requests.IncludeTypeBidQuoteAssets) {
			response.Included.Add(&quoteAsset)
		}
	}
	response.Data.Relationships.QuoteAssets = quoteAssets

	return response, nil
}
