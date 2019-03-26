package storage

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
	"gitlab.com/tokend/horizon/db2/history2"
)

// Vote is helper struct to operate with `votes`
type Vote struct {
	repo *db2.Repo
}

// NewVote - creates new instance of the `Vote`
func NewVote(repo *db2.Repo) *Vote {
	return &Vote{
		repo: repo,
	}
}

// Insert - inserts new vote
func (q *Vote) Insert(vote history2.Vote) error {
	sql := sq.Insert("votes").
		Columns(
			"voter_id", "poll_id", "choices",
		).
		Values(
			vote.VoterID, vote.PollID, vote.Choices,
		)

	_, err := q.repo.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to insert vote", logan.F{
			"voter_id": vote.VoterID,
			"poll_id":  vote.PollID,
		})
	}

	return nil
}

// Update - updates existing vote
func (q *Vote) Update(vote history2.Vote) error {
	sql := sq.Update("votes").SetMap(map[string]interface{}{
		"choices": vote.Choices,
	}).Where("voter_id = ?", vote.VoterID).Where("poll_id = ?", vote.PollID)

	_, err := q.repo.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to update vote", logan.F{
			"voter_id": vote.VoterID,
			"poll_id":  vote.PollID,
		})
	}

	return nil
}

// Update - updates existing vote
func (q *Vote) Remove(voterID string, pollID uint64) error {
	sql := sq.Delete("votes").Where("voter_id = ?", voterID).Where("poll_id = ?", pollID)

	_, err := q.repo.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to remove vote", logan.F{
			"voter_id": voterID,
			"poll_id":  pollID,
		})
	}

	return nil
}
