package core2

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

// VotesQ is a helper struct to aid in configuring queries that loads
// fee structs.
type VotesQ struct {
	repo     *db2.Repo
	selector sq.SelectBuilder
}

// NewVotesQ - creates new instance of Votesq
func NewVotesQ(repo *db2.Repo) VotesQ {
	return VotesQ{
		repo: repo,
		selector: sq.Select("v.poll_id", "v.voter_id", "v.choice").
			From("votes v"),
	}
}

// Page - returns Q with specified limit and offset params
func (q VotesQ) Page(params db2.OffsetPageParams) VotesQ {
	q.selector = params.ApplyTo(q.selector, "votes.poll_id", "votes.voter_id")
	return q
}

//FilterByVoter - returns q with filter by voter
func (q VotesQ) FilterByVoter(voter string) VotesQ {
	q.selector = q.selector.Where("v.voter_id = ?", voter)
	return q
}

//FilterByPoll - returns q with filter by poll
func (q VotesQ) FilterByPoll(pollID uint64) VotesQ {
	q.selector = q.selector.Where("v.poll_id = ?", pollID)
	return q
}

// Get - loads a row from `polls`
// returns nil, nil - if poll does not exists
// returns error if more than one poll found
func (q VotesQ) Get() (*Vote, error) {
	var result Vote
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
func (q VotesQ) Select() ([]Vote, error) {
	var result []Vote
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		if q.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load polls")
	}

	return result, nil
}
