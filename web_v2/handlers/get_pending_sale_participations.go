package handlers

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/horizon/web_v2/requests"
	"gitlab.com/tokend/horizon/web_v2/resources"
	"gitlab.com/tokend/regources/generated"
)

// GetPendingSaleParticipations - returns pending sale participations by open offers
func (h *getSaleParticipationsHandler) GetPendingSaleParticipations(request *requests.GetSaleParticipations) (*regources.SaleParticipationListResponse, error) {
	offers, err := h.getPendingOffers(request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get offers")
	}

	response := regources.SaleParticipationListResponse{
		Data: make([]regources.SaleParticipation, 0, len(offers)),
	}

	for _, offer := range offers {
		response.Data = append(response.Data, resources.NewSaleParticipation(
			offer.OfferID,
			offer.OwnerID,
			offer.BaseAssetCode,
			offer.QuoteAssetCode,
			offer.QuoteAmount),
		)

		if request.ShouldInclude(requests.IncludeTypeSaleParticipationsQuoteAsset) {
			if offer.QuoteAsset == nil {
				return nil, errors.From(errors.New("quote asset not found"), logan.F{
					"code": offer.QuoteAssetCode,
				})
			}

			quote := resources.NewAsset(*offer.QuoteAsset)
			response.Included.Add(&quote)
		}

		if request.ShouldInclude(requests.IncludeTypeSaleParticipationsBaseAsset) {
			if offer.BaseAsset == nil {
				return nil, errors.From(errors.New("base asset not found"), logan.F{
					"code": offer.BaseAssetCode,
				})
			}

			base := resources.NewAsset(*offer.BaseAsset)
			response.Included.Add(&base)
		}
	}

	if len(response.Data) > 0 {
		response.Links = request.GetCursorLinks(*request.PageParams, response.Data[len(response.Data)-1].ID)
	} else {
		response.Links = request.GetCursorLinks(*request.PageParams, "")
	}

	return &response, nil
}

func (h *getSaleParticipationsHandler) getPendingOffers(request *requests.GetSaleParticipations) ([]core2.Offer, error) {
	q := h.OffersQ.
		FilterByIsBuy(true).
		FilterByOrderBookID(int64(request.SaleID)).
		CursorPage(*request.PageParams)

	if request.ShouldFilter(requests.FilterTypeSaleParticipationsParticipant) {
		q = q.FilterByOwnerID(request.Filters.Participant)
	}

	if request.ShouldFilter(requests.FilterTypeSaleParticipationsQuoteAsset) {
		q = q.FilterByQuoteAssetCode(request.Filters.QuoteAsset)
	}

	if request.ShouldInclude(requests.IncludeTypeSaleParticipationsQuoteAsset) {
		q = q.WithQuoteAsset()
	}

	if request.ShouldInclude(requests.IncludeTypeSaleParticipationsBaseAsset) {
		q = q.WithBaseAsset()
	}

	offers, err := q.Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load offers from db")
	}

	return offers, nil
}
