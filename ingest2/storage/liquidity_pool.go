package storage

import (
	"gitlab.com/tokend/horizon/db2/history2"

	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// LiquidityPool is helper struct to operate with `liquidity_pools`
type LiquidityPool struct {
	repo *pgdb.DB
}

// NewLiquidityPool - creates new instance of the `LiquidityPool`
func NewLiquidityPool(repo *pgdb.DB) *LiquidityPool {
	return &LiquidityPool{
		repo: repo,
	}
}

// Insert - inserts new liquidity pool
func (q *LiquidityPool) Insert(lp history2.LiquidityPool) error {
	pairAsset := func(balance string) (string, error) {
		var assetCode string
		query := sq.Select("b.asset_code").
			From("balances b").
			Where("b.address = ?", balance)

		err := q.repo.Get(&assetCode, query)

		return assetCode, err
	}

	firstAssetCode, err := pairAsset(lp.FirstBalanceID)
	secondAssetCode, err := pairAsset(lp.SecondBalanceID)

	stmt := sq.Insert("liquidity_pools").
		Columns("id", "account", "token_asset", "first_balance", "second_balance", "tokens_amount",
			"first_reserve", "second_reserve", "first_asset_code", "second_asset_code").
		Values(lp.ID, lp.Account, lp.TokenAsset, lp.FirstBalanceID, lp.SecondBalanceID, lp.TokensAmount, lp.FirstReserve,
			lp.SecondReserve, firstAssetCode, secondAssetCode)

	err = q.repo.Exec(stmt)
	if err != nil {
		return errors.Wrap(err, "failed to insert liquidity pool", logan.F{"liquidity_pool_id": lp.ID})
	}

	return nil
}

// Update - updates existing liquidity pool
func (q *LiquidityPool) Update(lp history2.LiquidityPool) error {
	stmt := sq.Update("liquidity_pools").SetMap(map[string]interface{}{
		"tokens_amount":  lp.TokensAmount,
		"first_reserve":  lp.FirstReserve,
		"second_reserve": lp.SecondReserve,
	}).Where("id = ?", lp.ID)

	err := q.repo.Exec(stmt)
	if err != nil {
		return errors.Wrap(err, "failed to update liquidity pool", logan.F{"liquidity_pool_id": lp.ID})
	}

	return nil
}
