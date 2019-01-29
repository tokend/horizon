package regources

import "time"

//CreateSaleRequestAttrs -details of corresponding op
type CreateSaleRequest struct {
	Key
	Attributes CreateSaleRequestAttrs `json:"attributes"`
}

//CreateSaleRequestAttrs -details of corresponding op
type CreateSaleRequestAttrs struct {
	RequestID         int64            `json:"request_id"`
	BaseAsset         string           `json:"base_asset"`
	DefaultQuoteAsset string           `json:"default_quote_asset"`
	StartTime         time.Time        `json:"start_time"`
	EndTime           time.Time        `json:"end_time"`
	SoftCap           Amount           `json:"soft_cap"`
	HardCap           Amount           `json:"hard_cap"`
	QuoteAssets       []SaleQuoteAsset `json:"quote_assets"`
	Details           Details          `json:"details"`
}
