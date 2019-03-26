package regources

//VoteResponse - response for vote handler
type VoteResponse struct {
	Data     Vote     `json:"data"`
	Included Included `json:"included"`
}

type VotesResponse struct {
	Links    *Links   `json:"links"`
	Data     []Vote   `json:"data"`
	Included Included `json:"included"`
}

// Vote - Resource object representing VoteEntry
type Vote struct {
	Key
	Attributes    VoteAttrs     `json:"attributes"`
	Relationships VoteRelations `json:"relationships"`
}

type VoteAttrs struct {
	Choices []uint64 `json:"choices"`
}

type VoteRelations struct {
	Voter Relation `json:"voter"`
	Poll  Relation `json:"poll"`
}
