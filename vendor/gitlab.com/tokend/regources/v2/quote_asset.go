package regources

type QuoteAsset struct {
	Key
	Attributes QuoteAssetAttrs
}

type QuoteAssetAttrs struct {
	Price Amount `json:"price"`
}

type SaleQuoteAsset struct {
	Key
	Attributes    SaleQuoteAssetAttrs
	Relationships SaleQuoteAssetRelations
}

type SaleQuoteAssetAttrs struct {
	Price      string `json:"price"`
	CurrentCap string `json:"current_cap"`
	HardCap    string `json:"hard_cap"`
	SoftCap    string `json:"soft_cap,omitempty"`
}

type SaleQuoteAssetRelations struct {
	Asset *Relation `json:"asset"`
}
