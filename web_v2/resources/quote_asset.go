package resources

import "gitlab.com/tokend/regources/v2"

// newQuoteAssetKey - creates new key for resource type `quote-assets`
func newQuoteAssetKey(quoteAssetCode string) regources.Key {
	return regources.Key{
		ID:   quoteAssetCode,
		Type: regources.TypeQuoteAssets,
	}
}

// newQuoteAssetKey - creates new key for resource type `sale-quote-assets`
func newSaleQuoteAssetKey(quoteAssetCode string) regources.Key {
	return regources.Key{
		ID:   quoteAssetCode,
		Type: regources.TypeSaleQuoteAssets,
	}
}
