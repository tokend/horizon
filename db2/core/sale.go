package core

type Sale struct {
	ID                uint64 `db:"id"`
	OwnerID           string `db:"owner_id"`
	BaseAsset         string `db:"base_asset"`
	DefaultQuoteAsset string `db:"default_quote_asset"`
	StartTime         int64  `db:"start_time"`
	EndTime           int64  `db:"end_time"`
	SoftCap           uint64 `db:"soft_cap"`
	HardCap           uint64 `db:"hard_cap"`
	HardCapInBase     uint64 `db:"hard_cap_in_base"`
	CurrentCapInBase  uint64 `db:"current_cap_in_base"`
	Details           string `db:"details"`
	BaseBalance       string `db:"base_balance"`
}
