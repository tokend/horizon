package core2

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// AtomicSwapBidQI is a helper interface to aid in configuring queries
// that loads slices or entry of AtomicSwapBid structs.
type AtomicSwapQuoteAssetQI interface {
	// ByCode - returns atomic swap bid by code, if not found returns nil, nil
	ByID(bidIDs []int64) AtomicSwapQuoteAssetQI
	// Select - selects atomic swap bid for specified filter
	Select() ([]AtomicSwapQuoteAsset, error)
}

type AtomicSwapQuoteAssetQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

func (q *Q) AtomicSwapQuoteAsset() AtomicSwapQuoteAssetQI {
	return &AtomicSwapQuoteAssetQ{
		parent: q,
		sql:    selectAtomicSwapQuoteAsset,
	}
}

func (q *AtomicSwapQuoteAssetQ) ByID(bidIDs []int64) AtomicSwapQuoteAssetQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where(sq.Eq{"asqa.bid_id": bidIDs})
	return q
}

func (q *AtomicSwapQuoteAssetQ) Select() ([]AtomicSwapQuoteAsset, error) {
	if q.Err != nil {
		return nil, errors.Wrap(q.Err, "failed before select")
	}

	var result []AtomicSwapQuoteAsset
	err := q.parent.Select(&result, q.sql)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select atomic swap quote assets")
	}

	return result, nil
}

var selectAtomicSwapQuoteAsset = sq.Select(
	"asqa.bid_id",
	"asqa.quote_asset",
	"asqa.price",
).From("atomic_swap_quote_asset asqa")
