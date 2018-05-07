package resource

import (
	"strconv"

	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/tokend/go/amount"
)

type OrderBookEntry struct {
	PT      string `json:"paging_token"`
	OfferID uint64 `json:"offer_id,omitempty"`
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
