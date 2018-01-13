package resource

import (
	"strconv"

	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/horizon/db2/core"
)

type OrderBookEntry struct {
	PT      string `json:"paging_token"`
	OwnerID string `json:"owner_id,omitempty"`
	OfferData
}

func (o *OrderBookEntry) Populate(s *core.OrderBookEntry, baseAsset, quoteAsset string, isBuy bool) {
	o.OfferData.Populate(s, baseAsset, quoteAsset, isBuy)
	o.PT = o.PagingToken()

}

// PagingToken implementation for hal.Pageable
func (o *OrderBookEntry) PagingToken() string {
	return strconv.FormatInt(int64(amount.MustParse(o.Price)), 10)
}
