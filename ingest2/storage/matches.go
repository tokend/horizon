package storage

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
	"gitlab.com/tokend/horizon/db2/history2"
)

// Matches - helper struct to store `matches`
type Matches struct {
	repo *db2.Repo
}

// NewMatches - creates new instance of `Matches`
func NewMatches(repo *db2.Repo) *Matches {
	return &Matches{
		repo: repo,
	}
}

func convertMatch(match history2.Match) []interface{} {
	return []interface{}{
		match.ID,
		match.OrderBookID,
		match.OperationID,
		match.OfferID,
		match.BaseAmount,
		match.QuoteAmount,
		match.BaseAsset,
		match.QuoteAsset,
		match.Price,
	}
}

// Insert - inserts new match
func (s *Matches) Insert(matches []history2.Match) error {
	columns := []string{
		"id",
		"order_book_id",
		"operation_id",
		"offer_id",
		"base_amount",
		"quote_amount",
		"base_asset",
		"quote_asset",
		"price",
	}

	err := matchesBatchInsert(s.repo, matches, "matches", columns, convertMatch)
	if err != nil {
		return errors.Wrap(err, "failed to insert matches")
	}

	return nil
}
