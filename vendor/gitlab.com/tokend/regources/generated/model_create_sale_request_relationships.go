package regources

type CreateSaleRequestRelationships struct {
	BaseAsset         *Relation           `json:"base_asset,omitempty"`
	DefaultQuoteAsset *Relation           `json:"default_quote_asset,omitempty"`
	QuoteAssets       *RelationCollection `json:"quote_assets,omitempty"`
}
