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

	return response, nil
}
