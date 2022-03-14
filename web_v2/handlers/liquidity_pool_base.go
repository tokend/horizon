package handlers

import (
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/web_v2/requests"
	"gitlab.com/tokend/horizon/web_v2/resources"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	regources "gitlab.com/tokend/regources/generated"
)

type getLiquidityPoolBase struct {
	LiquidityPoolQ history2.LiquidityPoolQ
	Log            *logan.Entry
}

func (h *getLiquidityPoolBase) getAndPopulateResponse(q history2.LiquidityPoolQ, request *requests.GetLiquidityPool,
) (*regources.LiquidityPoolResponse, error) {
	historyLiquidityPool, err := q.Get()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get liquidity pool")
	}
	if historyLiquidityPool == nil {
		return nil, nil
	}

	response := &regources.LiquidityPoolResponse{
		Data: resources.NewLiquidityPool(*historyLiquidityPool),
	}

	if request.ShouldInclude(requests.IncludeTypeLiquidityPoolAssets) {
		firstAsset := resources.NewAssetV2(*historyLiquidityPool.FirstAsset)
		secondAsset := resources.NewAssetV2(*historyLiquidityPool.SecondAsset)
		lpTokensAsset := resources.NewAssetV2(*historyLiquidityPool.LPTokenAsset)

		response.Included.Add(&firstAsset)
		response.Included.Add(&secondAsset)
		response.Included.Add(&lpTokensAsset)
	}

	return response, nil
}

type liquidityPoolsBaseHandler struct {
	LiquidityPoolQ history2.LiquidityPoolQ
	Log            *logan.Entry
}

func applyLiquidityPoolFilters(r requests.LiquidityPoolsBase, q history2.LiquidityPoolQ) history2.LiquidityPoolQ {
	if r.ShouldFilter(requests.FilterTypeLiquidityPoolListAsset) {
		q = q.FilterByPairAsset(r.Filters.Asset)
	}

	if r.ShouldFilter(requests.FilterTypeLiquidityPoolListLPToken) {
		q = q.FilterByLPAsset(r.Filters.LPToken)
	}

	return q
}

func applyInclude(r requests.LiquidityPoolsBase, q history2.LiquidityPoolQ) history2.LiquidityPoolQ {
	if r.ShouldInclude(requests.IncludeTypeLiquidityPoolListAssets) {
		q = q.WithAssets()
	}

	return q
}

func (h *liquidityPoolsBaseHandler) populateResponse(historyLiquidityPools []history2.LiquidityPool,
	request requests.LiquidityPoolsBase,
	response *regources.LiquidityPoolListResponse) error {

	for _, historyLP := range historyLiquidityPools {
		lp := resources.NewLiquidityPool(historyLP)

		if request.ShouldInclude(requests.IncludeTypeLiquidityPoolListAssets) {
			firstAsset := resources.NewAssetV2(*historyLP.FirstAsset)
			secondAsset := resources.NewAssetV2(*historyLP.SecondAsset)
			lpTokensAsset := resources.NewAssetV2(*historyLP.LPTokenAsset)

			response.Included.Add(&firstAsset)
			response.Included.Add(&secondAsset)
			response.Included.Add(&lpTokensAsset)
		}

		response.Data = append(response.Data, lp)
	}

	return nil
}
