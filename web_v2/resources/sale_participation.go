package resources

import (
	"fmt"

	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/horizon/db2/history2"
	regources "gitlab.com/tokend/regources/generated"
)

func NewSaleParticipationKey(id uint64) regources.Key {
	return regources.Key{
		ID:   fmt.Sprintf("%d", id),
		Type: regources.SALE_PARTICIPATION,
	}
}

func NewSaleParticipation(id uint64, address, base, quote string, amount uint64) regources.SaleParticipation {
	return regources.SaleParticipation{
		Key: NewSaleParticipationKey(id),
		Relationships: regources.SaleParticipationRelationships{
			Participant: NewAccountKey(address).AsRelation(),
			QuoteAsset:  NewAssetKey(quote).AsRelation(),
			BaseAsset:   NewAssetKey(base).AsRelation(),
		},
		Attributes: regources.SaleParticipationAttributes{
			Amount: regources.Amount(amount),
		},
	}
}

func SaleParticipationsFromCore(offers []core2.Offer) []regources.SaleParticipation {
	result := make([]regources.SaleParticipation, 0, len(offers))
	for _, offer := range offers {
		result = append(result,
			NewSaleParticipation(offer.OfferID,
				offer.OwnerID,
				offer.BaseAsset.Code,
				offer.QuoteAsset.Code,
				offer.QuoteAmount))
	}
	return result
}

func SaleParticipationsFromHistory(participations []history2.SaleParticipation) []regources.SaleParticipation {
	result := make([]regources.SaleParticipation, 0, len(participations))
	for _, p := range participations {
		result = append(result,
			NewSaleParticipation(p.ID,
				p.ParticipantID,
				p.BaseAsset,
				p.QuoteAsset,
				uint64(p.QuoteAmount)))
	}
	return result
}
