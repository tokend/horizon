package resource

import (
	"bullioncoin.githost.io/development/go/amount"
	"gitlab.com/distributed_lab/tokend/horizon/db2/core"
	"strconv"
)

type OrderBookEntry struct {
	PT string `json:"paging_token"`
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
