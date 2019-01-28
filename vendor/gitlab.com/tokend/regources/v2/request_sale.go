package regources

import (
	"gitlab.com/tokend/go/xdr"
	"time"
)

// SaleRequest - represents details of the `sale` reviewable request
type SaleRequest struct {
	Key
	Attributes SaleRequestAttrs `json:"attributes"`
}

// SaleRequestAttrs - attributes of the `sale` reviewable request
type SaleRequestAttrs struct {
	BaseAssetForHardCap string       `json:"base_asset_for_hard_cap"`
	StartTime           time.Time    `json:"start_time"`
	EndTime             time.Time    `json:"end_time"`
	SaleType            xdr.SaleType `json:"sale_type"`
	Details             Details      `json:"details"`
}

// SaleRequestRelations - attributes of the `sale` reviewable request
type SaleRequestRelations struct {
	BaseAsset         *Relation           `json:"base_asset"`
	QuoteAssets       *RelationCollection `json:"quote_assets"`
	DefaultQuoteAsset *Relation           `json:"default_quote_asset"`
}
