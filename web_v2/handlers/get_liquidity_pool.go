package handlers

import (
	"net/http"

	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/requests"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	regources "gitlab.com/tokend/regources/generated"
)

func GetLiquidityPool(w http.ResponseWriter, r *http.Request) {
	historyRepo := ctx.HistoryRepo(r)

	handler := getLiquidityPoolHandler{
		getLiquidityPoolBase{
			LiquidityPoolQ: history2.NewLiquidityPoolQ(historyRepo),
			Log:            ctx.Log(r),
		},
	}

	request, err := requests.NewGetLiquidityPool(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.GetLiquidityPool(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get liquidity pool", logan.F{
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

type getLiquidityPoolHandler struct {
	getLiquidityPoolBase
}

func (h *getLiquidityPoolHandler) GetLiquidityPool(request *requests.GetLiquidityPool,
) (*regources.LiquidityPoolResponse, error) {
	q := h.LiquidityPoolQ.FilterByID(request.ID)

	if request.ShouldInclude(requests.IncludeTypeLiquidityPoolAssets) {
		q = q.WithAssets()
	}

	return h.getAndPopulateResponse(q, request)
}
