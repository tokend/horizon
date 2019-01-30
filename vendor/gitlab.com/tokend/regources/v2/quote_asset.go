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
	Price      Amount `json:"price"`
	CurrentCap Amount `json:"current_cap"`
	HardCap    Amount `json:"hard_cap"`
	SoftCap    Amount `json:"soft_cap"`
}

type SaleQuoteAssetRelations struct {
	Asset *Relation `json:"asset"`
}
