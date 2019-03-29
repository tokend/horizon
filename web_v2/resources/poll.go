package resources

import (
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/v2"
)

func NewPollKey(ID int64) regources.Key {
	return regources.NewKeyInt64(ID, regources.TypePolls)
}

func NewOutcomeKey(ID int64) regources.Key {
	return regources.NewKeyInt64(ID, regources.TypePollOutcome)
}

func NewPoll(record history2.Poll) regources.Poll {
	return regources.Poll{
		Key: NewPollKey(record.ID),
		Attributes: regources.PollAttrs{
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
			Outcome:        NewOutcomeKey(record.ID).AsRelation(),
		},
	}
}
