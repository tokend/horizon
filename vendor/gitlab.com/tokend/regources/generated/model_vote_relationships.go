package regources

type VoteRelationships struct {
	Poll  *Relation `json:"poll,omitempty"`
	Voter *Relation `json:"voter,omitempty"`
}
