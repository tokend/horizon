package history2

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

// MatchQ is a helper struct to aid in configuring queries that loads matches
type MatchQ struct {
	repo     *db2.Repo
	selector sq.SelectBuilder
}

// NewMatchQ returns new instance of MatchQ
func NewMatchQ(repo *db2.Repo) MatchQ {
	return MatchQ{
		repo: repo,
		selector: sq.Select(
			"m.participant_id",
			"m.order_book_id",
			"m.base_amount",
			"m.quote_amount",
			"m.base_asset",
			"m.quote_asset",
			"m.price",
		).From("matches m"),
	}
}

// NewSquashedMatchesQ returns new instance of MatchQ with matches grouped by operation id and price
func NewSquashedMatchesQ(repo *db2.Repo) MatchQ {
	return MatchQ{
		repo: repo,
		selector: sq.Select(
			"encode(format('%s:%s:%s:%s:%s', op.id, m.price, m.base_asset, m.quote_asset, m.order_book_id)::bytea, 'base64') id",
			"m.order_book_id",
			"sum(m.base_amount) base_amount",
			"sum(m.quote_amount) quote_amount",
			"m.base_asset",
			"m.quote_asset",
			"m.price",
			"op.ledger_close_time created_at",
		).
			From("matches m").
			LeftJoin("operations op ON m.operation_id = op.id").
			GroupBy(
				"op.id",
				"m.price",
				"m.base_asset",
				"m.quote_asset",
				"m.order_book_id",
			),
	}
}

// FilterByOrderBookID - returns Q with filter by order book ID
func (q MatchQ) FilterByOrderBookID(id uint64) MatchQ {
	q.selector = q.selector.Where("m.order_book_id = ?", id)
	return q
}

// FilterByBaseAsset - returns Q with filter by base asset
func (q MatchQ) FilterByBaseAsset(asset string) MatchQ {
	q.selector = q.selector.Where("m.base_asset = ?", asset)
	return q
}

// FilterByQuoteAsset - returns Q with filter by quote asset
func (q MatchQ) FilterByQuoteAsset(asset string) MatchQ {
	q.selector = q.selector.Where("m.quote_asset = ?", asset)
	return q
}

// Select - selects slice from the db, if no matches found - returns nil, nil
func (q MatchQ) Select() ([]Match, error) {
	var result []Match
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		if q.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load matches")
	}

	return result, nil
}
