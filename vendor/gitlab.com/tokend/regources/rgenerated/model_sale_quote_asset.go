package rgenerated

type SaleQuoteAsset struct {
	Key
	Attributes    *SaleQuoteAssetAttributes    `json:"attributes,omitempty"`
	Relationships *SaleQuoteAssetRelationships `json:"relationships,omitempty"`
}
type SaleQuoteAssetResponse struct {
	Data     SaleQuoteAsset `json:"data"`
	Included Included       `json:"included"`
}

type SaleQuoteAssetsResponse struct {
	Data     []SaleQuoteAsset `json:"data"`
	Included Included         `json:"included"`
	Links    *Links           `json:"links"`
}

// MustSaleQuoteAsset - returns SaleQuoteAsset from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustSaleQuoteAsset(key Key) *SaleQuoteAsset {
	var saleQuoteAsset SaleQuoteAsset
	if c.tryFindEntry(key, &saleQuoteAsset) {
		return &saleQuoteAsset
	}
	return nil
}
