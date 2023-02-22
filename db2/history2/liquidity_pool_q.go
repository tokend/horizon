package history2

import (
	"database/sql"

	"gitlab.com/tokend/horizon/db2"

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
			"lp.first_asset_code",
			"lp.second_asset_code",
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

// FilterByPairAsset - returns q with filter by LP pair asset code
func (q LiquidityPoolQ) FilterByPairAsset(pairAsset string) LiquidityPoolQ {
	q.selector = q.selector.Where(sq.Or{
		sq.Eq{
			"lp.first_asset_code": pairAsset,
		},
		sq.Eq{
			"lp.second_asset_code": pairAsset,
		},
	})
	return q
}

// FilterByExistentBalances - returns q, filtered by assets, that the user with `accountID` has
func (q LiquidityPoolQ) FilterByExistentBalances(accountID string, excludedAssets ...string) LiquidityPoolQ {
	q.selector = q.selector.
		Join("balances ba ON ba.asset_code IN (lp.first_asset_code, lp.second_asset_code)").
		Join("accounts a ON a.id = ba.account_id").
		Where(sq.And{
			sq.Eq{
				"a.address": accountID,
			},
			sq.NotEq{
				"ba.asset_code": excludedAssets,
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

// WithAssets - returns q with joined assets
func (q LiquidityPoolQ) WithAssets() LiquidityPoolQ {
	q.selector = q.selector.
		Columns(db2.GetColumnsForJoin(assetColumns, "first_asset")...).
		Columns(db2.GetColumnsForJoin(assetColumns, "second_asset")...).
		Columns(db2.GetColumnsForJoin(assetColumns, "lp_tokens_asset")...).
		LeftJoin("asset first_asset ON first_asset.code = lp.first_asset_code").
		LeftJoin("asset second_asset ON second_asset.code = lp.second_asset_code").
		LeftJoin("asset lp_tokens_asset ON lp_tokens_asset.code = lp.token_asset")

	return q
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
