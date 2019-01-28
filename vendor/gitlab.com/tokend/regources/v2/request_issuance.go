package regources

// IssuanceRequest - represents details of the `issuance` reviewable request
type IssuanceRequest struct {
	Key
	Attributes    IssuanceRequestAttrs     `json:"attributes"`
	Relationships IssuanceRequestRelations `json:"relationships"`
}

// IssuanceRequestAttrs - attributes of the `issuance` reviewable request
type IssuanceRequestAttrs struct {
	Amount  string  `json:"amount"`
	Details Details `json:"external_details"`
}

// IssuanceRequestRelations - relationships of the `issuance` reviewable request
type IssuanceRequestRelations struct {
	Asset    *Relation `json:"asset"`
	Receiver *Relation `json:"receiver"`
}
