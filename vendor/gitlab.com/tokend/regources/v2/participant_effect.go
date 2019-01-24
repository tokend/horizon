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
	Attributes    ParticipantEffectAttrs    `json:"attributes"`
}

//ParticipantEffectRelation - represents relations of resource
type ParticipantEffectRelation struct {
	Account   *Relation `json:"account"`
	Balance   *Relation `json:"balance"`
	Asset     *Relation `json:"asset"`
	Operation *Relation `json:"operation"`
}

//ParticipantEffectAttrs - attributes of participant effect
type ParticipantEffectAttrs struct {
	Effect *Effect `json:"effect"`
}
