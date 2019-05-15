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
		AssetsQ:        core2.NewAssetsQ(ctx.CoreRepo(r)),
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

	ape.Render(w, result)
}

type getSaleParticipationHandler struct {
	SalesQ         history2.SalesQ
	OffersQ        core2.OffersQ
	AssetsQ        core2.AssetsQ
	ParticipationQ history2.SaleParticipationQ
	Log            *logan.Entry
}

// GetSale returns sale with related resources
func (h *getSaleParticipationHandler) getSaleParticipation(sale *history2.Sale, request *requests.GetSaleParticipation) (*regources.SaleParticipationListResponse, error) {
	response := &regources.SaleParticipationListResponse{
		Data: make([]regources.SaleParticipation, 0),
	}

	switch sale.State {
	case regources.SaleStateOpen:
		q := populteOfferFilters(h.OffersQ, request)
		offers, err := q.
			FilterByIsBuy(true).
			FilterByOrderBookID(int64(request.SaleID)).
			CursorPage(*request.PageParams).
			Select()
		if err != nil {
			return nil, errors.Wrap(err, "failed to load sale participations from core")
		}

		err = h.fromCore(offers, request, response)
		if err != nil {
			return nil, errors.Wrap(err, "failed to populate response")
		}
	case regources.SaleStateCanceled:
	case regources.SaleStateClosed:
		q := populateSaleParticipationFilters(h.ParticipationQ, request)
		participations, err := q.
			FilterBySale(request.SaleID).
			Page(*request.PageParams).
			Select()
		if err != nil {
			return nil, errors.Wrap(err, "failed to load participations from history")
		}

		err = h.fromHistory(participations, request, response)
		if err != nil {
			return nil, errors.Wrap(err, "failed to populate response")
		}
	default:
		return nil, errors.From(errors.New("unexpected sale state"), logan.F{
			"sale_state": sale.State,
		})
	}

	h.populateLinks(response, request)

	return response, nil
}

func (h *getSaleParticipationHandler) populateLinks(
	response *regources.SaleParticipationListResponse, request *requests.GetSaleParticipation,
) {
	if len(response.Data) > 0 {
		response.Links = request.GetCursorLinks(*request.PageParams, response.Data[len(response.Data)-1].ID)
	} else {
		response.Links = request.GetCursorLinks(*request.PageParams, "")
	}
}

func populteOfferFilters(q core2.OffersQ, request *requests.GetSaleParticipation) core2.OffersQ {
	if request.ShouldFilter(requests.FilterTypeSaleParticipationParticipant) {
		q = q.FilterByOwnerID(request.Filters.Participant)
	}

	if request.ShouldFilter(requests.FilterTypeSaleParticipationQuoteAsset) {
		q = q.FilterByQuoteAssetCode(request.Filters.QuoteAsset)
	}

	if request.ShouldInclude(requests.IncludeTypeSaleParticipationQuoteAsset) {
		q = q.WithQuoteAsset()
	}
	if request.ShouldInclude(requests.IncludeTypeSaleParticipationBaseAsset) {
		q = q.WithBaseAsset()
	}

	return q
}

func populateSaleParticipationFilters(q history2.SaleParticipationQ, request *requests.GetSaleParticipation) history2.SaleParticipationQ {
	if request.ShouldFilter(requests.FilterTypeSaleParticipationParticipant) {
		q = q.FilterByParticipant(request.Filters.Participant)
	}

	if request.ShouldFilter(requests.FilterTypeSaleParticipationQuoteAsset) {
		q = q.FilterByQuoteAsset(request.Filters.QuoteAsset)
	}

	return q
}

func (h *getSaleParticipationHandler) fromHistory(
	participations []history2.SaleParticipation,
	request *requests.GetSaleParticipation,
	response *regources.SaleParticipationListResponse) error {

	for _, p := range participations {
		response.Data = append(response.Data,
			resources.NewSaleParticipation(p.ID,
				p.ParticipantID,
				p.BaseAsset,
				p.QuoteAsset,
				uint64(p.QuoteAmount)))
		err := h.populateIncludes(p.BaseAsset, p.QuoteAsset, request, response)
		if err != nil {
			return errors.Wrap(err, "failed to populate includes")
		}
	}

	return nil
}

func (h *getSaleParticipationHandler) populateIncludes(
	baseAsset, quoteAsset string,
	request *requests.GetSaleParticipation,
	response *regources.SaleParticipationListResponse) error {

	if request.ShouldInclude(requests.IncludeTypeSaleParticipationBaseAsset) {
		baseRaw, err := h.AssetsQ.GetByCode(baseAsset)
		if err != nil {
			return errors.Wrap(err, "failed to get asset by code", logan.F{
				"code": baseAsset,
			})
		}
		if baseRaw == nil {
			return errors.From(errors.New("asset not found"), logan.F{
				"code": baseAsset,
			})
		}

		base := resources.NewAsset(*baseRaw)

		response.Included.Add(&base)
	}

	if request.ShouldInclude(requests.IncludeTypeSaleParticipationQuoteAsset) {
		quoteRaw, err := h.AssetsQ.GetByCode(quoteAsset)
		if err != nil {
			return errors.Wrap(err, "failed to get asset by code", logan.F{
				"code": quoteAsset,
			})
		}
		if quoteRaw == nil {
			return errors.From(errors.New("asset not found"), logan.F{
				"code": quoteAsset,
			})
		}

		quote := resources.NewAsset(*quoteRaw)

		response.Included.Add(&quote)
	}
	return nil
}

func (h *getSaleParticipationHandler) fromCore(
	offers []core2.Offer,
	request *requests.GetSaleParticipation,
	response *regources.SaleParticipationListResponse) error {

	result := make([]regources.SaleParticipation, 0, len(offers))
	for _, offer := range offers {
		result = append(result,
			resources.NewSaleParticipation(offer.OfferID,
				offer.OwnerID,
				offer.BaseAssetCode,
				offer.QuoteAssetCode,
				offer.QuoteAmount))

		if request.ShouldInclude(requests.IncludeTypeSaleParticipationQuoteAsset) {
			if offer.QuoteAsset == nil {
				return errors.From(errors.New("asset not found"), logan.F{
					"code": offer.QuoteAssetCode,
				})
			}

			quote := resources.NewAsset(*offer.QuoteAsset)

			response.Included.Add(&quote)
		}

		if request.ShouldInclude(requests.IncludeTypeSaleParticipationBaseAsset) {
			if offer.BaseAsset == nil {
				return errors.From(errors.New("asset not found"), logan.F{
					"code": offer.BaseAssetCode,
				})
			}

			base := resources.NewAsset(*offer.BaseAsset)

			response.Included.Add(&base)
		}
	}

	return nil
}
