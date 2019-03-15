package core2

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

// AtomicSwapBidQ is a helper struct to aid in configuring queries that loads atomic swap bids
type AtomicSwapBidQ struct {
	repo     *db2.Repo
	selector sq.SelectBuilder
}

// NewAtomicSwapBidQ - creates new instance of AtomicSwapBidQ with no filters
func NewAtomicSwapBidQ(repo *db2.Repo) AtomicSwapBidQ {
	return AtomicSwapBidQ{
		repo: repo,
		selector: sq.Select(
			"b.bid_id",
			"b.owner_id",
			"b.base_asset_code",
			"b.base_balance_id",
			"b.base_amount",
			"b.locked_amount",
			"b.is_cancelled",
			"b.details",
			"b.created_at",
		).From("atomic_swap_bid b"),
	}
}

// GetByID - loads a row from `atomic_swap_bids` found with matching id
// returns nil, nil - if such atomic swap bid doesn't exists
func (q AtomicSwapBidQ) GetByID(id uint64) (*AtomicSwapBid, error) {
	return q.FilterByID(id).Get()
}

// FilterByCode - returns q with filter by code
func (q AtomicSwapBidQ) FilterByID(id uint64) AtomicSwapBidQ {
	q.selector = q.selector.Where("b.bid_id = ?", id)
	return q
}

// FilterByIDs - returns q with filter by ids
func (q AtomicSwapBidQ) FilterByIDs(ids []uint64) AtomicSwapBidQ {
	q.selector = q.selector.Where(sq.Eq{"b.bid_id": ids})
	return q
}

// FilterByOwner - returns q with filter by owner ID
func (q AtomicSwapBidQ) FilterByOwner(ownerID string) AtomicSwapBidQ {
	q.selector = q.selector.Where("b.owner_id = ?", ownerID)
	return q
}

func (q AtomicSwapBidQ) FilterByBaseAssets(codes []string) AtomicSwapBidQ {
	q.selector = q.selector.Where(sq.Eq{"b.base_asset_code": codes})
	return q
}

func (q AtomicSwapBidQ) FilterByQuoteAssets(codes []string) AtomicSwapBidQ {
	q.selector = q.selector.LeftJoin("atomic_swap_quote_asset qa ON qa.bid_id = b.bid_id").
		Where(sq.Eq{"qa.quote_asset": codes})
	return q
}

// Page - returns Q with specified limit and offset params
func (q AtomicSwapBidQ) Page(params db2.OffsetPageParams) AtomicSwapBidQ {
	q.selector = params.ApplyTo(q.selector, "b.bid_id")
	return q
}

// Get - loads a row from `atomic_swap_bids`
// returns nil, nil - if atomic swap bid does not exists
// returns error if more than one asset found
func (q AtomicSwapBidQ) Get() (*AtomicSwapBid, error) {
	var result AtomicSwapBid
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if q.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load atomic swap bid")
	}

	return &result, nil
}

func (q *AtomicSwapBidQ) Select() ([]AtomicSwapBid, error) {
	var result []AtomicSwapBid
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		if q.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load atomic swap bids")
	}

	return result, nil
}
