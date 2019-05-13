package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/requests"
	"gitlab.com/tokend/regources/generated"
)

// GetSaleParticipation - processes request to get list of sale participations
func GetSaleParticipation(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetSaleParticipation(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	handler := getSaleParticipationHandler{
		AssetsQ:        core2.NewAssetsQ(ctx.CoreRepo(r)),
		SalesQ:         history2.NewSalesQ(ctx.HistoryRepo(r)),
		ParticipationQ: history2.NewSaleParticipationQ(ctx.HistoryRepo(r)),
		OffersQ:        core2.NewOffersQ(ctx.CoreRepo(r)),
		Log:            ctx.Log(r),
	}

	sale, err := handler.getSale(request.SaleID)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get sale", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if sale == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if !isAllowed(r, w, sale.OwnerAddress) {
		return
	}

	result, err := handler.GetSaleParticipations(sale, request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get sale participations", logan.F{
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

type getSaleParticipationHandler struct {
	SalesQ         history2.SalesQ
	OffersQ        core2.OffersQ
	AssetsQ        core2.AssetsQ
	ParticipationQ history2.SaleParticipationQ
	Log            *logan.Entry
}

// GetSaleParticipations returns sale with related resources
func (h *getSaleParticipationHandler) GetSaleParticipations(sale *history2.Sale, request *requests.GetSaleParticipation) (*regources.SaleParticipationsResponse, error) {
	switch sale.State {
	case regources.SaleStateOpen:
		return h.GetPendingSaleParticipations(request)
	case regources.SaleStateCanceled:
		return nil, nil
	case regources.SaleStateClosed:
		return h.GetClosedSaleParticipations(request)
	default:
		return nil, errors.From(errors.New("unexpected sale state"), logan.F{
			"sale_state": sale.State,
		})
	}
}

func (h *getSaleParticipationHandler) getSale(id uint64) (*history2.Sale, error) {
	sale, err := h.SalesQ.GetByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load sale from db", logan.F{
			"id": id,
		})
	}

	if sale == nil {
		return nil, nil
	}

	return sale, nil
}
