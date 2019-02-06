package regources

//SignersResponse - response on signer request
type SignersResponse struct {
	Data     []Signer `json:"data"`
	Included Included `json:"included"`
}

type Signer struct {
	Key
	Attributes    SignerAttrs    `json:"attributes"`
	Relationships SignerRelation `json:"relationships"`
}

type SignerAttrs struct {
	Weight   int                    `json:"weight"`
	Identity int                    `json:"identity"`
	Details  map[string]interface{} `json:"details"`
}

type SignerRelation struct {
	Role *Relation `json:"role"`
}
