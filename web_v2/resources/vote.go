package resources

import (
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/v2"
)

func NewVoteKey(voterID string) regources.Key {
	return regources.Key{
		ID:   voterID,
		Type: regources.TypeVotes,
	}
}

func NewVote(record history2.Vote) regources.Vote {
	choices := make([]uint64, len(record.Choices))
	for _, c := range record.Choices {
		choices = append(choices, uint64(c))
	}
	return regources.Vote{
		Key: NewVoteKey(record.VoterID),
		Attributes: regources.VoteAttrs{
			Choices: choices,
		},
		Relationships: regources.VoteRelations{
			Voter: *NewAccountKey(record.VoterID).AsRelation(),
			Poll:  *NewPollKey(record.PollID).AsRelation(),
		},
	}
}
