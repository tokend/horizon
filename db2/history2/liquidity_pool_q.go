package history2

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// LiquidityPoolQ is a helper struct to aid in configuring queries that load liquidity pools
type LiquidityPoolQ struct {
	repo     *pgdb.DB
	selector sq.SelectBuilder
}

// NewLiquidityPoolQ - creates new instance of LiquidityPoolQ
func NewLiquidityPoolQ(repo *pgdb.DB) LiquidityPoolQ {
	return LiquidityPoolQ{
		repo: repo,
		selector: sq.Select(
			"lp.id",
			"lp.account",
			"lp.token_asset",
			"lp.first_balance",
			"lp.second_balance",
			"lp.tokens_amount",
			"lp.first_reserve",
			"lp.second_reserve",
			"lp.first_asset",
			"lp.second_asset",
		).From("liquidity_pools lp"),
	}
}

// FilterByID - returns q with filter by start_time
func (q LiquidityPoolQ) FilterByID(id uint64) LiquidityPoolQ {
	q.selector = q.selector.Where("lp.id = ?", id)
	return q
}

// FilterByLPAsset - returns q with filter by LP token asset
func (q LiquidityPoolQ) FilterByLPAsset(lpAsset string) LiquidityPoolQ {
	q.selector = q.selector.Where("lp.token_asset = ?", lpAsset)
	return q
}

// FilterByPairAssets - returns q with filter by LP pair asset codes
func (q LiquidityPoolQ) FilterByPairAssets(pairAssets []string) LiquidityPoolQ {
	q.selector = q.selector.Where(sq.Or{
		sq.Eq{
			"lp.first_asset": pairAssets,
		},
		sq.Eq{
			"lp.second_asset": pairAssets,
		},
	})
	return q
}

// GetByID - loads a row from `liquidity_pools` by ID
// returns nil, nil if liquidity pool doesn't exist
func (q LiquidityPoolQ) GetByID(id uint64) (*LiquidityPool, error) {
	return q.FilterByID(id).Get()
}

// GetByLPAsset - loads a row from `liquidity_pools` by LP token asset
// returns nil, nil if liquidity pool doesn't exist
func (q LiquidityPoolQ) GetByLPAsset(lpAsset string) (*LiquidityPool, error) {
	return q.FilterByLPAsset(lpAsset).Get()
}

// SelectByPairAssets - loads a slice of liquidity pools
// returns nil, nil if no liquidity pools found
func (q LiquidityPoolQ) SelectByPairAssets(pairAssets []string) ([]LiquidityPool, error) {
	return q.FilterByPairAssets(pairAssets).Select()
}

// Page - returns q with specified limit and offset params
func (q LiquidityPoolQ) Page(params pgdb.OffsetPageParams) LiquidityPoolQ {
	q.selector = params.ApplyTo(q.selector, "lp.id")
	return q
}

// CursorPage - returns q with specified limit and offset params
func (q LiquidityPoolQ) CursorPage(params pgdb.CursorPageParams) LiquidityPoolQ {
	q.selector = params.ApplyTo(q.selector, "lp.id")
	return q
}

// Get - loads a row from `liquidity_pools`
// returns nil, nil if liquidity pool doesn't exists
func (q LiquidityPoolQ) Get() (*LiquidityPool, error) {
	var result LiquidityPool

	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load liquidity pool")
	}

	return &result, nil
}

// Select - loads a slice of liquidity pools
// returns nil, nil if no liquidity pools found
func (q LiquidityPoolQ) Select() ([]LiquidityPool, error) {
	var result []LiquidityPool

	err := q.repo.Select(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load liquidity pools")
	}

	return result, err
}
