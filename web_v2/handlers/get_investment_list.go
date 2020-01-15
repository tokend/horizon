package handlers

import (
	"net/http"

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

func GetInvestmentList(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)

	converter, err := newBalanceStateConverterForHandler(coreRepo)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed failed to create balance state converted")
		ape.Render(w, problems.InternalError())
		return
	}

	handler := getInvestmentHandler{
		balanceStateConverter: converter,
		AssetsQ:               core2.NewAssetsQ(coreRepo),
		saleParticipationQ:    history2.NewSaleParticipationQ(ctx.HistoryRepo(r)),
		offersQ:               core2.NewOffersQ(ctx.CoreRepo(r)),
		Log:                   ctx.Log(r),
	}

	request, err := requests.NewGetInvestments(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if !isAllowed(r, w, request.AccountAddress) {
		return
	}

	result, err := handler.GetInvestments(request)
	if err != nil {
		ctx.Log(r).WithError(err).WithField("request", request).Error("failed to get converted balances")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if result == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	ape.Render(w, result)
}

type getInvestmentHandler struct {
	AssetsQ               core2.AssetsQ
	saleParticipationQ    history2.SaleParticipationQ
	offersQ               core2.OffersQ
	Log                   *logan.Entry
	balanceStateConverter *balanceStateConverter
}

func (h *getInvestmentHandler) GetInvestments(request *requests.GetInvestments) (*regources.InvestmentResponse, error) {
	coreAsset, err := h.AssetsQ.GetByCode(request.AssetCode)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get asset by code")
	}
	if coreAsset == nil {
		return nil, nil
	}

	closedSalesParticipations, err := h.saleParticipationQ.FilterByParticipant(request.AccountAddress).Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load closed sale participations")
	}

	pendingSaleParticipations, err := h.offersQ.FilterByOwnerID(request.AccountAddress).Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load opened sale participations")
	}

	var closedSaleResult int64
	for _, participation := range closedSalesParticipations {
		converted, err := h.balanceStateConverter.converter.TryToConvertWithOneHop(participation.BaseAmount, participation.BaseAsset, request.AssetCode)
		if err != nil {
			return nil, errors.Wrap(err, "fialed to convert sale amount")
		}

		closedSaleResult += *converted
	}

	var pendingSaleResult int64
	for _, participation := range pendingSaleParticipations {
		converted, err := h.balanceStateConverter.converter.TryToConvertWithOneHop(participation.BaseAmount, participation.BaseAssetCode, request.AssetCode)
		if err != nil {
			return nil, errors.Wrap(err, "fialed to convert sale amount")
		}

		pendingSaleResult += *converted
	}

	var response regources.InvestmentResponse
	response.Data.Attributes = regources.InvestmentAttributes{
		Asset:            request.AssetCode,
		ClosedSaleAmount: closedSaleResult,
		OpenedSaleAmount: pendingSaleResult,
		// TODO add amount not from sales
	}

	return &response, nil
}
