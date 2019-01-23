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

// GetAssetPair - processes request to get asset pair and it's details by it's ID (`base:quote`)
func GetAssetPair(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)
	handler := getAssetPairHandler{
		AssetPairsQ: core2.NewAssetPairsQ(coreRepo),
		AssetsQ:     core2.NewAssetsQ(coreRepo),
		Log:         ctx.Log(r),
	}

	request, err := requests.NewGetAssetPair(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.GetAssetPair(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get asset pair", logan.F{
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

type getAssetPairHandler struct {
	AssetPairsQ core2.AssetPairsQ
	AssetsQ     core2.AssetsQ
	Log         *logan.Entry
}

// GetAssetPair returns asset pair with related resources
func (h *getAssetPairHandler) GetAssetPair(request *requests.GetAssetPair) (*regources.AssetPairResponse, error) {
	coreAssetPair, err := h.AssetPairsQ.GetByBaseAndQuote(request.BaseAsset, request.QuoteAsset)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get asset pair by base and quote")
	}
	if coreAssetPair == nil {
		return nil, nil
	}

	response := &regources.AssetPairResponse{
		Data: resources.NewAssetPair(*coreAssetPair),
	}

	baseAssetKey := resources.NewAssetKey(coreAssetPair.Base)
	quoteAssetKey := resources.NewAssetKey(coreAssetPair.Quote)
	response.Data.Relationships.BaseAsset = baseAssetKey.AsRelation()
	response.Data.Relationships.QuoteAsset = quoteAssetKey.AsRelation()

	if request.ShouldInclude(requests.IncludeTypeAssetPairBaseAsset) {
		baseAsset, err := h.getAsset(coreAssetPair.Base)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to get base asset")
		}

		response.Included.Add(baseAsset)
	}

	if request.ShouldInclude(requests.IncludeTypeAssetPairQuoteAsset) {
		quoteAsset, err := h.getAsset(coreAssetPair.Quote)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to get quote asset")
		}

		response.Included.Add(quoteAsset)
	}

	return response, nil
}

func (h *getAssetPairHandler) getAsset(code string) (*regources.Asset, error) {
	coreAsset, err := h.AssetsQ.GetByCode(code)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get asset by code")
	}

	asset := resources.NewAsset(*coreAsset)

	return &asset, nil
}
