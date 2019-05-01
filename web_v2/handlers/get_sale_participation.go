package handlers

import (
	"net/http"

	"gitlab.com/tokend/horizon/web_v2/resources"

	"gitlab.com/tokend/horizon/db2/core2"

	regources "gitlab.com/tokend/regources/generated"

	"gitlab.com/tokend/horizon/db2/history2"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/requests"
)

func GetSaleParticipation(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetSaleParticipation(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	handler := getSaleParticipationHandler{
		SalesQ:         history2.NewSalesQ(ctx.HistoryRepo(r)),
		ParticipationQ: history2.NewSaleParticipationQ(ctx.HistoryRepo(r)),
		OffersQ:        core2.NewOffersQ(ctx.CoreRepo(r)),
		Log:            ctx.Log(r),
	}

	sale, err := handler.SalesQ.GetByID(request.SaleID)
	if err != nil {
		ctx.Log(r).WithError(err).WithFields(logan.F{
			"sale_id": request.SaleID,
		}).Error("failed to get sale by ID")
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

	result, err := handler.getSaleParticipation(sale, request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get whitelist", logan.F{
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
	ParticipationQ history2.SaleParticipationQ
	Log            *logan.Entry
}

// GetSale returns sale with related resources
func (h *getSaleParticipationHandler) getSaleParticipation(sale *history2.Sale, request *requests.GetSaleParticipation) (*regources.SaleParticipationsResponse, error) {
	response := &regources.SaleParticipationsResponse{}

	switch sale.State {
	case regources.SaleStateOpen:
		offers, err := h.OffersQ.
			FilterByIsBuy(true).
			FilterByOrderBookID(int64(request.SaleID)).
			CursorPage(*request.PageParams).
			Select()
		if err != nil {
			return nil, errors.Wrap(err, "failed to load sale participations from core")
		}

		response.Data = resources.SaleParticipationsFromCore(offers)

	case regources.SaleStateCanceled:
		return nil, nil
	case regources.SaleStateClosed:
		participations, err := h.ParticipationQ.FilterBySale(request.SaleID).Page(*request.PageParams).Select()
		if err != nil {
			return nil, errors.Wrap(err, "failed to load participations from history")
		}

		response.Data = resources.SaleParticipationsFromHistory(participations)
	default:
		return nil, errors.From(errors.New("unexpected sale state"), logan.F{
			"sale_state": sale.State,
		})
	}
	return response, nil
}
