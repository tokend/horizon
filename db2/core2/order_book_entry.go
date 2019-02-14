package core2

import "time"

// OrderBookEntry - db representation of offer
type OrderBookEntry struct {
	ID             string    `db:"id"`
	OrderBookID    uint64    `db:"order_book_id"`
	BaseAssetCode  string    `db:"base_asset_code"`
	QuoteAssetCode string    `db:"quote_asset_code"`
	BaseAmount     uint64    `db:"base_amount"`
	QuoteAmount    uint64    `db:"quote_amount"`
	IsBuy          bool      `db:"is_buy"`
	Price          int64     `db:"price"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`

	BaseAsset  *Asset `db:"base_assets"`
	QuoteAsset *Asset `db:"quote_assets"`
}
