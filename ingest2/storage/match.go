package storage

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
	"gitlab.com/tokend/horizon/db2/history2"
)

// Match - helper struct to store `matches`
type Match struct {
	repo *db2.Repo
}

// NewMatch - creates new instance of `Match`
func NewMatch(repo *db2.Repo) *Match {
	return &Match{
		repo: repo,
	}
}

// Insert - inserts new match
func (q *Match) Insert(match history2.Match) error {
	query := sq.Insert("matches").SetMap(map[string]interface{}{
		"order_book_id":  match.OrderBookID,
		"participant_id": match.ParticipantID,
		"base_amount":    match.BaseAmount,
		"quote_amount":   match.QuoteAmount,
		"base_asset":     match.BaseAsset,
		"quote_asset":    match.QuoteAsset,
		"price":          match.Price,
	})

	_, err := q.repo.Exec(query)
	if err != nil {
		return errors.Wrap(err, "failed to insert match")
	}

	return nil
}
