package rgenerated

type SaleRelationships struct {
	BaseAsset         *Relation           `json:"base_asset,omitempty"`
	DefaultQuoteAsset *Relation           `json:"default_quote_asset,omitempty"`
	Owner             *Relation           `json:"owner,omitempty"`
	QuoteAssets       *RelationCollection `json:"quote_assets,omitempty"`
}
