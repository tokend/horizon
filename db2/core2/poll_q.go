package core2

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

// PollsQ is a helper struct to aid in configuring queries that loads
// fee structs.
type PollsQ struct {
	repo     *db2.Repo
	selector sq.SelectBuilder
}

// NewPollsQ - creates new instance of Pollsq
func NewPollsQ(repo *db2.Repo) PollsQ {
	return PollsQ{
		repo: repo,
		selector: sq.Select(
			"p.permission_type", "p.number_of_choices",
			"p.type", "p.data", "p.start_time",
			"p.end_time", "p.owner_id", "p.result_provider_id",
			"p.vote_confirmation_required", "p.details", "p.lastmodified").
			From("polls p"),
	}
}

// Page - returns Q with specified limit and offset params
func (q PollsQ) Page(params db2.OffsetPageParams) PollsQ {
	q.selector = params.ApplyTo(q.selector, "polls.id")
	return q
}

//FilterByOwner - returns q with filter by owner
func (q PollsQ) FilterByOwner(owner string) PollsQ {
	q.selector = q.selector.Where("p.owner_id = ?", owner)
	return q
}

//FilterByResultProvider - returns q with filter by result provider
func (q PollsQ) FilterByResultProvider(resultProviderID string) PollsQ {
	q.selector = q.selector.Where("p.result_provider_id = ?", resultProviderID)
	return q
}

//FilterByPermissionType - returns q with filter by permission type
func (q PollsQ) FilterByPermissionType(permissionType uint64) PollsQ {
	q.selector = q.selector.Where("p.permission_type = ?", permissionType)
	return q
}

//FilterByType - returns q with filter by poll type
func (q PollsQ) FilterByType(pollType int32) PollsQ {
	q.selector = q.selector.Where("p.type = ?", pollType)
	return q
}

//FilterByVoteConfirmation- returns q with filter by vote confirmation
func (q PollsQ) FilterByVoteConfirmation(voteConfirmationRequired bool) PollsQ {
	q.selector = q.selector.Where("p.vote_confirmation_required = ?", voteConfirmationRequired)
	return q
}

func (q PollsQ) FilterByPollID(ID uint64) PollsQ {
	q.selector = q.selector.Where("p.id = ?", ID)
	return q
}

// GetByOfferID - loads a row from `offers` found with offer ID
// returns nil, nil - if such offer doesn't exist
func (q PollsQ) GetByPollID(id uint64) (*Poll, error) {
	return q.FilterByPollID(id).Get()
}

// Get - loads a row from `polls`
// returns nil, nil - if poll does not exists
// returns error if more than one poll found
func (q PollsQ) Get() (*Poll, error) {
	var result Poll
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if q.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load poll")
	}

	return &result, nil
}

//Select - selects slice from the db, if no polls found - returns nil, nil
func (q PollsQ) Select() ([]Poll, error) {
	var result []Poll
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		if q.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load polls")
	}

	return result, nil
}
