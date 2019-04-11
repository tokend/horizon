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
	"gitlab.com/tokend/regources/generated"
	"net/http"
)

func GetBalance(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)

	handler := getBalanceHandler{
		BalancesQ: core2.NewBalancesQ(coreRepo),
		Log:       ctx.Log(r),
	}

	request, err := requests.NewGetBalance(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if !isAllowed(r, w) {
		ape.RenderErr(w, problems.NotAllowed())
		return
	}

	result, err := handler.GetBalance(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get balance", logan.F{
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

type getBalanceHandler struct {
	BalancesQ core2.BalancesQ
	AssetsQ   core2.AssetsQ
	Log       *logan.Entry
}

func (h *getBalanceHandler) GetBalance(request *requests.GetBalance) (*regources.BalanceResponse, error) {
	balance, err := h.BalancesQ.GetByAddress(request.ID)
	if err != nil {
		return nil, errors.Wrap(err, "cannot filter balances by ID", logan.F{
			"id": request.ID,
		})
	}
	if balance == nil {
		return nil, nil
	}

	response := regources.BalanceResponse{
		Data: *resources.NewBalance(balance),
	}

	response.Data.Relationships = &regources.BalanceRelationships{
		Account: resources.NewAccountKey(balance.AccountAddress).AsRelation(),
		Asset:   resources.NewAssetKey(balance.AssetCode).AsRelation(),
		State:   resources.NewBalanceState(balance).AsRelation(),
	}

	return &response, nil
}
