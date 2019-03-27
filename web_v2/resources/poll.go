package resources

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/v2"
)

func NewPollKey(ID int64) regources.Key {
	return regources.NewKeyInt64(ID, regources.TypePolls)
}

func NewPoll(record history2.Poll, votes []history2.Vote) regources.Poll {
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
			VotesCount:               VoteCount(votes),
		},
		Relationships: regources.PollRelations{
			ResultProvider: NewAccountKey(record.ResultProviderID).AsRelation(),
			Owner:          NewAccountKey(record.OwnerID).AsRelation(),
			Votes:          VotesAsRelations(votes),
		},
	}
}

func VotesAsRelations(votes []history2.Vote) *regources.RelationCollection {
	keys := make([]regources.Key, 0, len(votes))
	for _, vote := range votes {
		keys = append(keys, regources.Key{
			ID:   vote.VoterID,
			Type: regources.TypeVotes,
		})
	}

	return &regources.RelationCollection{
		Data: keys,
	}
}

func VoteCount(votes []history2.Vote) []regources.VoteCount {
	voteCount := make([]regources.VoteCount, 0)

	helperMap := make(map[uint64]uint32)
	for _, vote := range votes {
		switch vote.VoteData.Type {
		case xdr.PollTypeSingleChoice:
			helperMap[vote.VoteData.SingleChoice]++
		default:
			panic("Unexpected vote data type (poll type)")
		}
	}

	for ch, count := range helperMap {
		choice := regources.VoteCount{
			Choice: uint64(ch),
			Count:  count,
		}
		voteCount = append(voteCount, choice)
	}
	return voteCount
}
