package regources

// PublicKeyResponse - response on /public_key request
type PublicKeyResponse struct {
	Data     PublicKey `json:"data"`
	Included Included  `json:"included"`
}

// PublicKey - Resource object representing "public key" resource
type PublicKey struct {
	Key
	Relationships PublicKeyRelationships `json:"relationships"`
}

// PublicKeyRelationships - the relationships of the PublicKey resource
type PublicKeyRelationships struct {
	Accounts *RelationCollection `json:"accounts"`
}
