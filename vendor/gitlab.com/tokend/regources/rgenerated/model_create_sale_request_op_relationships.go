package rgenerated

type CreateSaleRequestOpRelationships struct {
	BaseAsset         *Relation           `json:"base_asset,omitempty"`
	DefaultQuoteAsset *Relation           `json:"default_quote_asset,omitempty"`
	QuoteAssets       *RelationCollection `json:"quote_assets,omitempty"`
	Request           *Relation           `json:"request,omitempty"`
}
