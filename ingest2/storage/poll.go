package storage

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
	"gitlab.com/tokend/horizon/db2/history2"
	regources "gitlab.com/tokend/regources/v2"
)

// CreatePoll is helper struct to operate with `polls`
type Poll struct {
	repo *db2.Repo
}

// NewPoll - creates new instance of the `CreatePoll`
func NewPoll(repo *db2.Repo) *Poll {
	return &Poll{
		repo: repo,
	}
}

// Insert - inserts new poll
func (q *Poll) Insert(poll history2.Poll) error {

	sql := sq.Insert("polls").
		Columns(
			"id", "permission_type", "number_of_choices", "data", "start_time",
			"end_time", "owner_id", "result_provider_id",
			"vote_confirmation_required", "details", "state",
		).
		Values(
			poll.ID, poll.PermissionType, poll.NumberOfChoices, poll.Data, poll.StartTime,
			poll.EndTime, poll.OwnerID, poll.ResultProviderID, poll.VoteConfirmationRequired,
			poll.Details, poll.State,
		)

	_, err := q.repo.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to insert poll", logan.F{"poll_id": poll.ID})
	}

	return nil
}

// Update - updates existing poll
func (q *Poll) Update(poll history2.Poll) error {
	sql := sq.Update("polls").SetMap(map[string]interface{}{
		"permission_type":            poll.PermissionType,
		"number_of_choices":          poll.NumberOfChoices,
		"data":                       poll.Data,
		"start_time":                 poll.StartTime,
		"end_time":                   poll.EndTime,
		"owner_id":                   poll.OwnerID,
		"result_provider_id":         poll.ResultProviderID,
		"vote_confirmation_required": poll.VoteConfirmationRequired,
		"details":                    poll.Details,
		"state":                      poll.State,
	}).Where("id = ?", poll.ID)

	_, err := q.repo.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to update poll", logan.F{"poll_id": poll.ID})
	}

	return nil
}

// SetState - sets state
func (q *Poll) SetState(pollID uint64, state regources.PollState) error {
	sql := sq.Update("polls").Set("state", state).Where("id = ?", pollID)
	_, err := q.repo.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to set state", logan.F{"poll_id": pollID})
	}

	return nil
}
