package reviewablerequest2

import (
	"time"

	"gitlab.com/tokend/regources/valueflag"
)

type SaleCreationRequest struct {
	BaseAsset           string                 `json:"base_asset"`
	DefaultQuoteAsset   string                 `json:"default_quote_asset"`
	StartTime           time.Time              `json:"start_time"`
	EndTime             time.Time              `json:"end_time"`
	SoftCap             string                 `json:"soft_cap"`
	HardCap             string                 `json:"hard_cap"`
	SaleType            valueflag.Flag         `json:"sale_type"`
	BaseAssetForHardCap string                 `json:"base_asset_for_hard_cap"`
	Details             map[string]interface{} `json:"details"`
	QuoteAssets         []SaleQuoteAsset       `json:"quote_assets"`
	State               valueflag.Flag         `json:"state"`
}

type SaleQuoteAsset struct {
	QuoteAsset string `json:"quote_asset"`
	Price      string `json:"price"`
}
