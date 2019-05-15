package resources

import (
	"fmt"

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
