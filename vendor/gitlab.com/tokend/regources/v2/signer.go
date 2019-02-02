package regources

import "encoding/json"

//SignersResponse - response on signer request
type SignersResponse struct {
	Data     []Signer `json:"data"`
	Included Included `json:"included"`
}

func (r SignersResponse) MarshalJSON() ([]byte, error) {
	if r.Data == nil {
		r.Data = []Signer{}
	}

	type temp SignersResponse
	return json.Marshal(temp(r))
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
