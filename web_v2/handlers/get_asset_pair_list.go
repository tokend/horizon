package handlers

import (
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/requests"
	"gitlab.com/tokend/horizon/web_v2/resources"
	"gitlab.com/tokend/regources/v2"
	"net/http"
)

// GetAssetPairList - processes request to get the list of asset pairs and their details
func GetAssetPairList(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)
	handler := getAssetPairListHandler{
		AssetPairsQ: core2.NewAssetPairsQ(coreRepo),
		AssetsQ:     core2.NewAssetsQ(coreRepo),
		Log:         ctx.Log(r),
	}

	request, err := requests.NewGetAssetPairList(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.GetAssetPairList(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get asset list", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

type getAssetPairListHandler struct {
	AssetPairsQ core2.AssetPairsQ
	AssetsQ     core2.AssetsQ
	Log         *logan.Entry
}

// GetAssetPairList returns asset pair list with related resources
func (h *getAssetPairListHandler) GetAssetPairList(request *requests.GetAssetPairList) (*regources.AssetPairsResponse, error) {
	q := h.AssetPairsQ.Page(request.PageParams.Limit(), request.PageParams.Offset())

	if request.ShouldFilter(requests.FilterTypeAssetPairListAsset) {
		q = q.FilterByAsset(request.Filters.Asset)
	}

	if request.ShouldFilter(requests.FilterTypeAssetPairListBaseAsset) {
		q = q.FilterByBaseAsset(request.Filters.BaseAsset)
	}

	if request.ShouldFilter(requests.FilterTypeAssetPairListQuoteAsset) {
		q = q.FilterByQuoteAsset(request.Filters.QuoteAsset)
	}

	if request.ShouldFilter(requests.FilterTypeAssetPairListPolicy) {
		q = q.FilterByPolicy(request.Filters.Policy)
	}

	coreAssetPairs, err := q.Select()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get asset pair list")
	}

	response := &regources.AssetPairsResponse{
		Data:  make([]regources.AssetPair, 0, len(coreAssetPairs)),
		Links: request.PageParams.Links(request.URL()),
	}

	for i := range coreAssetPairs {
		assetPair := resources.NewAssetPair(coreAssetPairs[i])

		baseAssetKey := resources.NewAssetKey(coreAssetPairs[i].BaseAsset)
		quoteAssetKey := resources.NewAssetKey(coreAssetPairs[i].QuoteAsset)
		assetPair.Relationships.BaseAsset = baseAssetKey.AsRelation()
		assetPair.Relationships.QuoteAsset = quoteAssetKey.AsRelation()

		if request.ShouldInclude(requests.IncludeTypeAssetPairListBaseAssets) {
			baseAsset, err := h.getAsset(coreAssetPairs[i].BaseAsset)
			if err != nil {
				return nil, errors.Wrap(err, "Failed to get base asset")
			}

			response.Included.Add(baseAsset)
		}
		if request.ShouldInclude(requests.IncludeTypeAssetPairListQuoteAssets) {
			quoteAsset, err := h.getAsset(coreAssetPairs[i].BaseAsset)
			if err != nil {
				return nil, errors.Wrap(err, "Failed to get quote asset")
			}

			response.Included.Add(quoteAsset)
		}

		response.Data = append(response.Data, assetPair)
	}

	return response, nil
}

func (h *getAssetPairListHandler) getAsset(code string) (*regources.Asset, error) {
	coreAsset, err := h.AssetsQ.GetByCode(code)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get asset by code")
	}

	asset := resources.NewAsset(*coreAsset)

	return &asset, nil
}
