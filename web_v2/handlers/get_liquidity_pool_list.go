package handlers

import (
	"net/http"

	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/requests"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	regources "gitlab.com/tokend/regources/generated"
)

type getLiquidityPoolListHandler struct {
	liquidityPoolsBaseHandler
}

func GetLiquidityPoolList(w http.ResponseWriter, r *http.Request) {
	historyRepo := ctx.HistoryRepo(r)

	handler := getLiquidityPoolListHandler{
		liquidityPoolsBaseHandler{
			LiquidityPoolQ: history2.NewLiquidityPoolQ(historyRepo),
			Log:            ctx.Log(r),
		},
	}

	request, err := requests.NewGetLiquidityPoolList(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.GetLiquidityPoolHandler(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error(logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

func (h *getLiquidityPoolListHandler) GetLiquidityPoolHandler(request *requests.GetLiquidityPoolList) (*regources.LiquidityPoolListResponse, error) {
	q := applyLiquidityPoolFilters(request.LiquidityPoolsBase, h.LiquidityPoolQ)

	q = applyInclude(request.LiquidityPoolsBase, q)

	historyPools, err := q.Page(*request.PageParams).Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get liquidity pools list")
	}

	response := &regources.LiquidityPoolListResponse{
		Data:  make([]regources.LiquidityPool, 0, len(historyPools)),
		Links: request.GetOffsetLinks(*request.PageParams),
	}

	err = h.populateResponse(historyPools, request.LiquidityPoolsBase, response)
	if err != nil {
		return nil, errors.Wrap(err, "failed to populate response")
	}

	return response, nil
}
