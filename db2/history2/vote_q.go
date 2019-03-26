package history2

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

// VotesQ is a helper struct to aid in configuring queries that loads
// poll structures.
type VotesQ struct {
	repo     *db2.Repo
	selector sq.SelectBuilder
}

// NewVotesQ - creates new instance of VotesQ
func NewVotesQ(repo *db2.Repo) VotesQ {
	return VotesQ{
		repo: repo,
		selector: sq.Select(
			"v.poll_id",
			"v.voter_id",
			"v.choices",
		).From("votes p"),
	}
}

// FilterByID - returns q with filter by sale ID
func (q VotesQ) FilterByVoter(voterID string) VotesQ {
	q.selector = q.selector.Where("votes.voter_id = ?", voterID)
	return q
}

// FilterByOwner - returns q with filter by Owner
func (q VotesQ) FilterByPollID(pollID int64) VotesQ {
	q.selector = q.selector.Where("votes.poll_id = ?", pollID)
	return q
}

// Page - returns Q with specified limit and offset params
func (q VotesQ) Page(params db2.OffsetPageParams) VotesQ {
	q.selector = params.ApplyTo(q.selector, "votes.voter_id")
	return q
}

// Get - loads a row from `sales`
// returns nil, nil - if sale does not exists
// returns error if more than one Vote found
func (q VotesQ) Get() (*Vote, error) {
	var result Vote
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if q.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load sale")
	}

	return &result, nil
}

//Select - selects slice from the db, if no sales found - returns nil, nil
func (q VotesQ) Select() ([]Vote, error) {
	var result []Vote
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load sales")
	}

	return result, nil
}
