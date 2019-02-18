package resources

import (
	"fmt"
	"time"

	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/regources/v2"
)

// NewOffer creates new instance of Offer from provided one
func NewOffer(record core2.Offer) regources.Offer {
	return regources.Offer{
		Key: regources.Key{
			ID:   fmt.Sprint(record.OfferID),
			Type: regources.TypeOffers,
		},
		Attributes: regources.OfferAttrs{
			IsBuy:       record.IsBuy,
			OrderBookID: fmt.Sprint(record.OrderBookID),
			CreatedAt:   time.Unix(record.CreatedAt, 0).UTC().Format(time.RFC3339),
			BaseAmount:  regources.Amount(record.BaseAmount),
			QuoteAmount: regources.Amount(record.QuoteAmount),
			Price:       regources.Amount(record.Price),
			Fee: regources.Fee{
				CalculatedPercent: regources.Amount(record.Fee),
			},
		},
	}
}
