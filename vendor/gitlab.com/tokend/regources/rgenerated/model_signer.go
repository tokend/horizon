package rgenerated

type Signer struct {
	Key
	Attributes    *SignerAttributes    `json:"attributes,omitempty"`
	Relationships *SignerRelationships `json:"relationships,omitempty"`
}
type SignerResponse struct {
	Data     Signer   `json:"data"`
	Included Included `json:"included"`
}

type SignersResponse struct {
	Data     []Signer `json:"data"`
	Included Included `json:"included"`
	Links    *Links   `json:"links"`
}

// MustSigner - returns Signer from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustSigner(key Key) *Signer {
	var signer Signer
	if c.tryFindEntry(key, &signer) {
		return &signer
	}
	return nil
}
