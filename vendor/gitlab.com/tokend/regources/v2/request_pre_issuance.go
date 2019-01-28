package regources

// PreIssuanceRequest - represents details of the `pre-issuance` reviewable request
type PreIssuanceRequest struct {
	Key
	Attributes    PreIssuanceRequestAttrs     `json:"attributes"`
	Relationships PreIssuanceRequestRelations `json:"relationships"`
}

// PreIssuanceRequestAttrs - attributes of the `pre_issuance` reviewable request
type PreIssuanceRequestAttrs struct {
	Asset     string `json:"asset"`
	Amount    string `json:"amount"`
	Signature string `json:"signature"`
	Reference string `json:"reference"`
}

// PreIssuanceRequestRelations - relationships of the `pre_issuance` reviewable request
type PreIssuanceRequestRelations struct {
	Asset *Relation `json:"asset"`
}
