package regources

//SaleQuoteAsset - represents quote asset of the sale
type SaleQuoteAsset struct {
	QuoteAsset string `json:"quote_asset"`
	Price      Amount `json:"price"`
}
