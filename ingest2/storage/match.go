package storage

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
	"gitlab.com/tokend/horizon/db2/history2"
	"math"
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

func (q *Match) latestMatchID() (int32, error) {
	var result int32
	err := q.repo.GetRaw(&result, "SELECT coalesce(max(m.id), 0) from matches m")
	if err != nil {
		return 0, errors.Wrap(err, "failed to load max ID")
	}

	return result, nil
}

// Insert - inserts new match
func (q *Match) Insert(match history2.Match) error {
	latestID, err := q.latestMatchID()
	if err != nil {
		return errors.Wrap(err, "failed to get latest match ID")
	}
	if latestID == math.MaxInt32 {
		return errors.New("failed to generate new match ID - latest match ID is already maxInt32")
	}

	query := sq.Insert("matches").SetMap(map[string]interface{}{
		"id":             latestID + 1,
		"order_book_id":  match.OrderBookID,
		"operation_id":   match.OperationID,
		"offer_id":       match.OfferID,
		"participant_id": match.ParticipantID,
		"base_amount":    match.BaseAmount,
		"quote_amount":   match.QuoteAmount,
		"base_asset":     match.BaseAsset,
		"quote_asset":    match.QuoteAsset,
		"price":          match.Price,
	})

	_, err = q.repo.Exec(query)
	if err != nil {
		return errors.Wrap(err, "failed to insert match")
	}

	return nil
}
