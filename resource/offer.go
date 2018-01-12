package resource

import (
	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/horizon/db2/core"
	"strconv"
)

type Offer struct {
	PT             string `json:"paging_token"`
	OwnerID        string `json:"owner_id"`
	OfferID        uint64 `json:"offer_id"`
	BaseBalanceID  string `json:"base_balance_id"`
	QuoteBalanceID string `json:"quote_balance_id"`
	Fee            string `json:"fee"`
	OfferData
}

func (o *Offer) Populate(s *core.Offer) {
	o.OwnerID = s.OwnerID
	o.OfferID = s.OfferID
	o.OfferData.Populate(&s.OrderBookEntry, s.BaseAssetCode, s.QuoteAssetCode, s.IsBuy)
	o.BaseBalanceID = s.BaseBalanceID
	o.QuoteBalanceID = s.QuoteBalanceID
	o.Fee = amount.StringU(s.Fee)
	o.PT = o.PagingToken()
}

// PagingToken implementation for hal.Pageable
func (o *Offer) PagingToken() string {
	return strconv.FormatUint(o.OfferID, 10)
}
