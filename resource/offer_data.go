package resource

import (
	"time"

	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/horizon/db2/core"
)

type OfferData struct {
	BaseAssetCode  string    `json:"base_asset_code"`
	QuoteAssetCode string    `json:"quote_asset_code"`
	IsBuy          bool      `json:"is_buy"`
	BaseAmount     string    `json:"base_amount"`
	QuoteAmount    string    `json:"quote_amount"`
	Price          string    `json:"price"`
	CreatedAt      time.Time `json:"created_at"`
}

func (o *OfferData) Populate(s *core.OrderBookEntry, baseAsset, quoteAsset string, isBuy bool) {
	o.BaseAssetCode = baseAsset
	o.QuoteAssetCode = quoteAsset
	o.IsBuy = isBuy
	o.BaseAmount = amount.String(s.BaseAmount)
	o.QuoteAmount = amount.String(s.QuoteAmount)
	o.Price = amount.String(s.Price)
	o.CreatedAt = time.Unix(s.CreatedAt, 0).UTC()
}
