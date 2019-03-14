package core2

type AtomicSwapBid struct {
	BidID           int64  `db:"bid_id"`
	OwnerID         string `db:"owner_id"`
	BaseAsset       string `db:"base_asset_code"`
	BaseBalanceID   string `db:"base_balance_id"`
	AvailableAmount uint64 `db:"base_amount"`
	LockedAmount    uint64 `db:"locked_amount"`
	IsCanceled      bool   `db:"is_cancelled"`
	Details         string `db:"details"`
	CreatedAt       int64  `db:"created_at"`
}

type AtomicSwapQuoteAsset struct {
	BidID      int64  `db:"bid_id"`
	QuoteAsset string `db:"quote_asset"`
	Price      uint64 `db:"price"`
}
