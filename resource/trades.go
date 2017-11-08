package resource

import (
	"strconv"
	"time"

	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/horizon/db2/history"
)

type Trades struct {
	PT          string    `json:"paging_token"`
	ID          int64     `json:"id"`
	BaseAsset   string    `json:"base_asset"`
	QuoteAsset  string    `json:"quote_asset"`
	BaseAmount  string    `json:"base_amount"`
	QuoteAmount string    `json:"quote_asset"`
	Price       string    `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
}

func (t *Trades) Populate(h *history.Trades) {
	t.ID = h.ID
	t.PT = strconv.FormatInt(h.ID, 10)
	t.BaseAsset = h.BaseAsset
	t.QuoteAsset = h.QuoteAsset
	t.BaseAmount = amount.String(h.BaseAmount)
	t.QuoteAmount = amount.String(h.QuoteAmount)
	t.Price = amount.String(h.Price)
	t.CreatedAt = h.CreatedAt
}

// PagingToken implementation for hal.Pageable
func (t *Trades) PagingToken() string {
	return t.PT
}
