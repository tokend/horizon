package history2

import (
	"time"

	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

// PollsQ is a helper struct to aid in configuring queries that loads
// poll structures.
type PollsQ struct {
	repo     *db2.Repo
	selector sq.SelectBuilder
}

// NewPollsQ - creates new instance of PollsQ
func NewPollsQ(repo *db2.Repo) PollsQ {
	return PollsQ{
		repo: repo,
		selector: sq.Select(
			"p.id",
			"p.permission_type",
			"p.number_of_choices",
			"p.type",
			"p.start_time",
			"p.end_time",
			"p.owner_id",
			"p.result_provider_id",
			"p.vote_confirmation_required",
			"p.details",
			"p.state",
		).From("polls p"),
	}
}

// FilterByID - returns q with filter by poll ID
func (q PollsQ) FilterByID(id int64) PollsQ {
	q.selector = q.selector.Where("polls.id = ?", id)
	return q
}

// GetByID loads a row from `polls`, by ID
// returns nil, nil - if poll does not exists
func (q PollsQ) GetByID(id int64) (*Poll, error) {
	return q.FilterByID(id).Get()
}

// FilterByOwner - returns q with filter by Owner
func (q PollsQ) FilterByOwner(ownerID string) PollsQ {
	q.selector = q.selector.Where("polls.owner_id = ?", ownerID)
	return q
}

// FilterByResultProvider - returns q with filter by Owner
func (q PollsQ) FilterByResultProvider(resultProviderID string) PollsQ {
	q.selector = q.selector.Where("polls.result_provider_id = ?", resultProviderID)
	return q
}

// FilterByPermissionType - returns q with filter by BaseAsset
func (q PollsQ) FilterByPermissionType(permissionType uint64) PollsQ {
	q.selector = q.selector.Where("polls.permission_type = ?", permissionType)
	return q
}

// FilterByMaxEndTime - returns q with filter by max end time
func (q PollsQ) FilterByMaxEndTime(time time.Time) PollsQ {
	q.selector = q.selector.Where("polls.end_time <= ?", time)
	return q
}

// FilterByMaxStartTime - returns q with filter by start_time
func (q PollsQ) FilterByMaxStartTime(time time.Time) PollsQ {
	q.selector = q.selector.Where("polls.start_time <= ?", time)
	return q
}

// FilterByMinStartTime - returns q with filter by start_time
func (q PollsQ) FilterByMinStartTime(time time.Time) PollsQ {
	q.selector = q.selector.Where("polls.start_time >= ?", time)
	return q
}

// FilterByMinEndTime - returns q with filter by end_time
func (q PollsQ) FilterByMinEndTime(time time.Time) PollsQ {
	q.selector = q.selector.Where("polls.end_time >= ?", time)
	return q
}

// FilterByPollType - returns q with filter by type
func (q PollsQ) FilterByPollType(pollType uint64) PollsQ {
	q.selector = q.selector.Where("polls.type = ?", pollType)
	return q
}

// FilterByPollState - returns q with filter by type
func (q PollsQ) FilterByState(state int32) PollsQ {
	q.selector = q.selector.Where("polls.state = ?", state)
	return q
}

// FilterByVoteConfirmation- returns q with filter by type
func (q PollsQ) FilterByVoteConfirmation(voteConfirmation bool) PollsQ {
	q.selector = q.selector.Where("polls.vote_confirmation_required = ?", voteConfirmation)
	return q
}

// Page - returns Q with specified limit and offset params
func (q PollsQ) Page(params db2.OffsetPageParams) PollsQ {
	q.selector = params.ApplyTo(q.selector, "polls.id")
	return q
}

// Get - loads a row from `polls`
// returns nil, nil - if poll does not exists
// returns error if more than one Poll found
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
		return nil, errors.Wrap(err, "failed to load polls")
	}

	return result, nil
}
