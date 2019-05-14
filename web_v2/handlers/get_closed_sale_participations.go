package handlers

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/amount"

	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/web_v2/requests"
	"gitlab.com/tokend/horizon/web_v2/resources"
	"gitlab.com/tokend/regources/generated"
)

// GetClosedSaleParticipations - returns closed sale participations by completed matches
func (h *getSaleParticipationHandler) GetClosedSaleParticipations(request *requests.GetSaleParticipation) (*regources.SaleParticipationsResponse, error) {
	matches, err := h.getMatches(request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get matches")
	}

	response := regources.SaleParticipationsResponse{
		Data: make([]regources.SaleParticipation, 0, len(matches)),
	}

	for _, m := range matches {
		response.Data = append(response.Data, resources.NewSaleParticipation(
			m.ID,
			m.ParticipantID,
			m.BaseAsset,
			m.QuoteAsset,
			amount.MustParseU(m.QuoteAmount),
		))

		if request.ShouldInclude(requests.IncludeTypeSaleParticipationBaseAsset) {
			coreBaseAsset, err := h.AssetsQ.GetByCode(m.BaseAsset)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get asset by code", logan.F{
					"code": m.BaseAsset,
				})
			}
			if coreBaseAsset == nil {
				return nil, errors.From(errors.New("asset not found"), logan.F{
					"code": m.BaseAsset,
				})
			}

			base := resources.NewAsset(*coreBaseAsset)
			response.Included.Add(&base)
		}

		if request.ShouldInclude(requests.IncludeTypeSaleParticipationQuoteAsset) {
			coreQuoteAsset, err := h.AssetsQ.GetByCode(m.QuoteAsset)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get asset by code", logan.F{
					"code": m.QuoteAsset,
				})
			}
			if coreQuoteAsset == nil {
				return nil, errors.From(errors.New("asset not found"), logan.F{
					"code": m.QuoteAsset,
				})
			}

			quote := resources.NewAsset(*coreQuoteAsset)
			response.Included.Add(&quote)
		}
	}

	if len(response.Data) > 0 {
		response.Links = request.GetCursorLinks(*request.PageParams, response.Data[len(response.Data)-1].ID)
	} else {
		response.Links = request.GetCursorLinks(*request.PageParams, "")
	}

	return &response, nil
}

func (h *getSaleParticipationHandler) getMatches(request *requests.GetSaleParticipation) ([]history2.SaleParticipation2, error) {
	q := h.ParticipationQ.FilterBySale(request.SaleID).Page(*request.PageParams)

	if request.ShouldFilter(requests.FilterTypeSaleParticipationParticipant) {
		q = q.FilterByParticipant(request.Filters.Participant)
	}

	if request.ShouldFilter(requests.FilterTypeSaleParticipationQuoteAsset) {
		q = q.FilterByQuoteAsset(request.Filters.QuoteAsset)
	}

	matches, err := q.Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load matches from db")
	}

	return matches, nil
}
