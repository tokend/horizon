package core

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

// AtomicSwapBidQI is a helper interface to aid in configuring queries
// that loads slices or entry of AtomicSwapBid structs.
type AtomicSwapBidQI interface {
	// ByCode - returns atomic swap bid by code, if not found returns nil, nil
	ByID(bidID int64) (*AtomicSwapBidEntry, error)
	// ForOwner - filters atomic swap bid by owner
	ForOwner(owner string) AtomicSwapBidQI
	// ForCodes - filters atomic swap bid by base asset code
	ForBaseAsset(codes []string) AtomicSwapBidQI
	// Page - applies page params
	Page(page db2.PageQuery) AtomicSwapBidQI
	// Select - selects atomic swap bid for specified filter
	Select() ([]AtomicSwapBidEntry, error)
}

type AtomicSwapBidQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

func (q *Q) AtomicSwapBid() AtomicSwapBidQI {
	return &AtomicSwapBidQ{
		parent: q,
		sql:    selectAtomicSwapBid,
	}
}

func (q *AtomicSwapBidQ) ByID(bidID int64) (*AtomicSwapBidEntry, error) {
	if q.Err != nil {
		return nil, q.Err
	}

	q.sql = q.sql.Where(sq.Eq{"asb.bid_id": bidID})

	var result AtomicSwapBidEntry
	err := q.parent.Get(&result, q.sql)
	if q.parent.NoRows(err) {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to get atomic swap bid")
	}

	return &result, err
}

func (q *AtomicSwapBidQ) ForOwner(ownerID string) AtomicSwapBidQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where(sq.Eq{"asb.owner_id": ownerID})
	return q
}

func (q *AtomicSwapBidQ) ForBaseAsset(codes []string) AtomicSwapBidQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where(sq.Eq{"asb.base_asset_code": codes})
	return q
}

func (q *AtomicSwapBidQ) Page(page db2.PageQuery) AtomicSwapBidQI {
	if q.Err != nil {
		return q
	}

	q.sql, q.Err = page.ApplyTo(q.sql, "bid_id")
	return q
}

func (q *AtomicSwapBidQ) Select() ([]AtomicSwapBidEntry, error) {
	if q.Err != nil {
		return nil, errors.Wrap(q.Err, "failed before select")
	}

	var result []AtomicSwapBidEntry
	err := q.parent.Select(&result, q.sql)
	if q.parent.NoRows(err) {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to select atomic swap bids")
	}

	return result, nil
}

var selectAtomicSwapBid = sq.Select(
	"asb.bid_id",
	"asb.owner_id",
	"asb.base_asset_code",
	"asb.base_balance_id",
	"asb.base_amount",
	"asb.locked_amount",
	"asb.is_cancelled",
	"asb.details",
	"asb.created_at",
).From("atomic_swap_bid asb")
