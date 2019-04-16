package resources

import (
	"gitlab.com/tokend/horizon/db2/history2"
	regources "gitlab.com/tokend/regources/generated"
)

func NewVoteKey(voterID string) regources.Key {
	return regources.Key{
		ID:   voterID,
		Type: regources.VOTES,
	}
}

func NewVote(record history2.Vote) regources.Vote {
	return regources.Vote{
		Key: NewVoteKey(record.VoterID),
		Attributes: regources.VoteAttributes{
			VoteData: regources.VoteData(record.VoteData),
		},
		Relationships: regources.VoteRelationships{
			Voter: NewAccountKey(record.VoterID).AsRelation(),
			Poll:  NewPollKey(record.PollID).AsRelation(),
		},
	}
}
