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
			"m.id",
			"m.operation_id",
			"m.offer_id",
			"m.base_amount",
			"m.quote_amount",
			"m.base_asset",
			"m.quote_asset",
			"m.price",
		).From("matches m"),
	}
}

// WithCreatedAt - returns Q with `created_at` column
func (q MatchQ) WithCreatedAt() MatchQ {
	q.selector = q.selector.Join("operations op ON op.id = m.operation_id").Columns("ledger_close_time created_at")
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

// Page - apply paging params to the query
func (q MatchQ) Page(pageParams db2.CursorPageParams) MatchQ {
	q.selector = pageParams.ApplyTo(q.selector, "m.id")
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
