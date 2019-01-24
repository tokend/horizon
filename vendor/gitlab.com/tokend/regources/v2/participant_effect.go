package regources

//ParticipantEffectsResponse - response for history and movements handlers
type ParticipantEffectsResponse struct {
	Links    *Links              `json:"links"`
	Data     []ParticipantEffect `json:"data"`
	Included Included            `json:"included"`
}

//ParticipantEffect - represent account effected by operation
type ParticipantEffect struct {
	Key
	Relationships ParticipantEffectRelation `json:"relationships"`
}

//ParticipantEffectRelation - represents relations of resource
type ParticipantEffectRelation struct {
	Account   *Relation `json:"account"`
	Balance   *Relation `json:"balance,omitempty"`
	Asset     *Relation `json:"asset,omitempty"`
	Operation *Relation `json:"operation"`
	Effect    *Relation `json:"effect"`
}
