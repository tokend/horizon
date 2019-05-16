package handlers

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/horizon/web_v2/requests"
	"gitlab.com/tokend/horizon/web_v2/resources"
	"gitlab.com/tokend/regources/generated"
)

type pendingParticipationsQ struct {
	offersQ core2.OffersQ
}

func newPendingParticipationQ(request *requests.GetSaleParticipations, q core2.OffersQ) pendingParticipationsQ {
	q = q.
		FilterByIsBuy(true).
		FilterByOrderBookID(int64(request.SaleID)).
		CursorPage(*request.PageParams)

	return pendingParticipationsQ{
		offersQ: q,
	}
}

// FilterByParticipant - filters out participations by participant address
func (p pendingParticipationsQ) FilterByParticipant(address string) participationsQ {
	p.offersQ = p.offersQ.FilterByOwnerID(address)
	return p
}

// FilterByQuoteAsset - filters out participations by quote asset
func (p pendingParticipationsQ) FilterByQuoteAsset(code string) participationsQ {
	p.offersQ = p.offersQ.FilterByQuoteAssetCode(code)
	return p
}

// Select - select records from db and wraps them to participations
func (p pendingParticipationsQ) Select() ([]regources.SaleParticipation, error) {
	offers, err := p.offersQ.Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load offers from db")
	}

	result := make([]regources.SaleParticipation, 0, len(offers))

	for _, offer := range offers {
		result = append(result, resources.NewSaleParticipation(
			offer.OfferID,
			offer.OwnerID,
			offer.BaseAssetCode,
			offer.QuoteAssetCode,
			offer.QuoteAmount),
		)
	}

	return result, nil
}
