package resources

import (
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/v2"
)

func NewPollKey(ID int64) regources.Key {
	return regources.NewKeyInt64(ID, regources.TypePolls)
}

func NewParticipationKey(ID int64) regources.Key {
	return regources.NewKeyInt64(ID, regources.TypePollsParticipation)
}

func NewPoll(record history2.Poll) regources.Poll {
	return regources.Poll{
		Key: NewPollKey(record.ID),
		Attributes: regources.PollAttributes{
			VoteConfirmationRequired: record.VoteConfirmationRequired,
			PermissionType:           record.PermissionType,
			NumberOfChoices:          record.NumberOfChoices,
			StartTime:                record.StartTime,
			EndTime:                  record.EndTime,
			Details:                  record.Details,
			PollData:                 record.Data,
			PollState:                record.State,
		},
		Relationships: regources.PollRelations{
			ResultProvider: NewAccountKey(record.ResultProviderID).AsRelation(),
			Owner:          NewAccountKey(record.OwnerID).AsRelation(),
			Participation:  NewParticipationKey(record.ID).AsRelation(),
		},
	}
}
func NewParticipation(id int64, historyVotes []history2.Vote) regources.PollParticipation {
	outcome := regources.PollParticipation{
		Key: NewParticipationKey(id),
		Relationships: regources.PollParticipationRelations{
			Votes: &regources.RelationCollection{},
		},
	}
	for _, v := range historyVotes {
		vote := NewVoteKey(v.VoterID)
		outcome.Relationships.Votes.Data = append(outcome.Relationships.Votes.Data,
			vote)
	}

	return outcome
}
