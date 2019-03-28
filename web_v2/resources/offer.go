package resources

import (
	"fmt"
	"time"

	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/regources/rgenerated"
)

// NewOffer creates new instance of Offer from provided one
func NewOffer(record core2.Offer) rgenerated.Offer {
	return rgenerated.Offer{
		Key: rgenerated.Key{
			ID:   fmt.Sprint(record.OfferID),
			Type: rgenerated.OFFERS,
		},
		Attributes: rgenerated.OfferAttributes{
			IsBuy:       record.IsBuy,
			OrderBookId: fmt.Sprint(record.OrderBookID),
			CreatedAt:   time.Unix(record.CreatedAt, 0).UTC().Format(time.RFC3339),
			BaseAmount:  rgenerated.Amount(record.BaseAmount),
			QuoteAmount: rgenerated.Amount(record.QuoteAmount),
			Price:       rgenerated.Amount(record.Price),
			Fee: rgenerated.Fee{
				CalculatedPercent: rgenerated.Amount(record.Fee),
			},
		},
	}
}
