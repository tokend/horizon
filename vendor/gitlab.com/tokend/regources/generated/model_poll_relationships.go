package regources

type PollRelationships struct {
	Owner          *Relation `json:"owner,omitempty"`
	Participation  *Relation `json:"participation,omitempty"`
	ResultProvider *Relation `json:"result_provider,omitempty"`
}
