package resource

import (
	"strconv"

	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/tokend/regources"
	"time"
)

func PopulateOffer(o core.Offer) regources.Offer {
	return regources.Offer{
		OwnerID:        o.OwnerID,
		OfferID:        o.OfferID,
		OrderBookID:    o.OrderBookID,
		OfferData:      PopulateOfferData(o),
		BaseBalanceID:  o.BaseBalanceID,
		QuoteBalanceID: o.QuoteBalanceID,
		Fee:            regources.Amount(o.Fee),
		PT:             strconv.FormatUint(o.OfferID, 10),
	}
}

func PopulateOfferData(o core.Offer) regources.OfferData {
	return regources.OfferData{
		BaseAssetCode:  o.BaseAssetCode,
		QuoteAssetCode: o.QuoteAssetCode,
		IsBuy:          o.IsBuy,
		BaseAmount:     regources.Amount(o.BaseAmount),
		QuoteAmount:    regources.Amount(o.QuoteAmount),
		Price:          regources.Amount(o.Price),
		CreatedAt:      time.Unix(o.CreatedAt, 0).UTC(),
	}
}
