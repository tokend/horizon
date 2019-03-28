package resources

import "gitlab.com/tokend/regources/rgenerated"

// newQuoteAssetKey - creates new key for resource type `quote-assets`
func newQuoteAssetKey(quoteAssetCode string) rgenerated.Key {
	return rgenerated.Key{
		ID:   quoteAssetCode,
		Type: rgenerated.QUOTE_ASSETS,
	}
}

// newQuoteAssetKey - creates new key for resource type `sale-quote-assets`
func newSaleQuoteAssetKey(quoteAssetCode string) rgenerated.Key {
	return rgenerated.Key{
		ID:   quoteAssetCode,
		Type: rgenerated.SALE_QUOTE_ASSETS,
	}
}
