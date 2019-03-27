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
	return regources.Vote{
		Key: NewVoteKey(record.VoterID),
		Attributes: regources.VoteAttrs{
			VoteData: record.VoteData,
		},
		Relationships: regources.VoteRelations{
			Voter: NewAccountKey(record.VoterID).AsRelation(),
			Poll:  NewPollKey(record.PollID).AsRelation(),
		},
	}
}
