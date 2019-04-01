package rgenerated

type Asset struct {
	Key
	Attributes    *AssetAttributes    `json:"attributes,omitempty"`
	Relationships *AssetRelationships `json:"relationships,omitempty"`
}
type AssetResponse struct {
	Data     Asset    `json:"data"`
	Included Included `json:"included"`
}

type AssetsResponse struct {
	Data     []Asset  `json:"data"`
	Included Included `json:"included"`
	Links    *Links   `json:"links"`
}

// MustAsset - returns Asset from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustAsset(key Key) *Asset {
	var asset Asset
	if c.tryFindEntry(key, &asset) {
		return &asset
	}
	return nil
}
