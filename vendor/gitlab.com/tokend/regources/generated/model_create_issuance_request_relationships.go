package regources

type CreateIssuanceRequestRelationships struct {
	Asset    *Relation `json:"asset,omitempty"`
	Receiver *Relation `json:"receiver,omitempty"`
}
