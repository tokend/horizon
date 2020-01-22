package handlers

import (
	"net/http"

	"gitlab.com/tokend/go/amount"

	"gitlab.com/tokend/horizon/db2/history2"

	"gitlab.com/distributed_lab/logan/v3/errors"
	regources "gitlab.com/tokend/regources/generated"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/requests"
)

func GetBalanceStatistic(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)

	converter, err := newBalanceStateConverterForHandler(coreRepo)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed failed to create balance state converted")
		ape.Render(w, problems.InternalError())
		return
	}

	handler := getBalancesStatisticHandler{
		balanceStateConverter:  converter,
		AssetsQ:                core2.NewAssetsQ(coreRepo),
		BalancesQ:              core2.NewBalancesQ(coreRepo),
		saleConvertedBalancesQ: history2.NewSaleConvertedBalancesQ(ctx.HistoryRepo(r)),
		offersQ:                core2.NewOffersQ(ctx.CoreRepo(r)),
		Log:                    ctx.Log(r),
	}

	request, err := requests.NewGetBalancesStatistic(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if !isAllowed(r, w, request.AccountAddress) {
		return
	}

	result, err := handler.GetBalancesStatistic(request)
	if err != nil {
		ctx.Log(r).WithError(err).WithField("request", request).Error("failed to get balances statistic")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if result == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	ape.Render(w, result)
}

type getBalancesStatisticHandler struct {
	AssetsQ                core2.AssetsQ
	saleConvertedBalancesQ history2.SaleConvertedBalancesQ
	offersQ                core2.OffersQ
	BalancesQ              core2.BalancesQ
	Log                    *logan.Entry
	balanceStateConverter  *balanceStateConverter
}

func (h *getBalancesStatisticHandler) GetBalancesStatistic(request *requests.GetBalancesStatistic) (*regources.BalancesStatisticResponse, error) {
	coreAsset, err := h.AssetsQ.GetByCode(request.AssetCode)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get asset by code")
	}
	if coreAsset == nil {
		return nil, nil
	}

	coreBalances, err := h.BalancesQ.FilterByAccount(request.AccountAddress).Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get balances by account address")
	}

	closedSalesParticipations, err := h.saleConvertedBalancesQ.FilterByParticipant(request.AccountAddress).Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load closed sale participations")
	}

	pendingSaleParticipations, err := h.offersQ.FilterByOwnerID(request.AccountAddress).Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load opened sale participations")
	}

	var closedSaleResult int64
	for _, participation := range closedSalesParticipations {
		baseAmount := amount.MustParseU(participation.BaseAmount)
		converted, err := h.balanceStateConverter.converter.TryToConvertWithOneHop(int64(baseAmount), participation.BaseAsset, request.AssetCode)
		if err != nil {
			return nil, errors.Wrap(err, "fialed to convert sale amount")
		}
		if converted == nil {
			continue
		}
		closedSaleResult += *converted
	}

	var pendingSaleResult int64
	for _, participation := range pendingSaleParticipations {
		converted, err := h.balanceStateConverter.converter.TryToConvertWithOneHop(int64(participation.BaseAmount), participation.BaseAssetCode, request.AssetCode)
		if err != nil {
			return nil, errors.Wrap(err, "fialed to convert sale amount")
		}
		if converted == nil {
			continue
		}
		pendingSaleResult += *converted
	}
	var fullBalanceResult int64
	for _, coreBalance := range coreBalances {
		converted, err := h.balanceStateConverter.converter.TryToConvertWithOneHop(int64(coreBalance.Amount), coreBalance.AssetCode, request.AssetCode)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get converted balance state")
		}
		if converted == nil {
			continue
		}
		fullBalanceResult = *converted
	}

	var response regources.BalancesStatisticResponse
	response.Data.Key = regources.Key{
		ID:   request.AccountAddress,
		Type: regources.BALANCES_STATISTIC,
	}
	response.Data.Attributes = regources.BalancesStatisticAttributes{
		Asset:              request.AssetCode,
		ClosedSalesAmount:  regources.Amount(closedSaleResult),
		PendingSalesAmount: regources.Amount(pendingSaleResult),
		FullAmount:         regources.Amount(fullBalanceResult),
	}

	return &response, nil
}
