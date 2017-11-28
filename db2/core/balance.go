package core

import (
	"database/sql"

	sq "github.com/lann/squirrel"
)

type Balance struct {
	BalanceID                string `db:"balance_id"`
	AccountID                string `db:"account_id"`
	Asset                    string `db:"asset"`
	Amount                   int64  `db:"amount"`
	Locked                   int64  `db:"locked"`

	IncentivePerCoin int64
}

func (q *Q) BalancesByAddress(dest interface{}, addy string) error {
	sql := selectBalance.Where("ba.account_id = ?", addy)
	return q.Select(dest, sql)
}

func (q *Q) BalanceByID(dest interface{}, bid string) error {
	sql := selectBalance.Where("ba.balance_id = ?", bid)
	return q.Get(dest, sql)
}

var selectBalance = sq.Select(
	"ba.balance_id",
	"ba.account_id",
	"ba.asset",
	"ba.amount",
	"ba.locked",
).From("balance ba")

var selectCoinsInCirculationAmounts = sq.Select(
	"b.asset as asset, sum(b.amount + b.locked) as amount").
	From("balance b").
	GroupBy("b.asset")

func (q *Q) AssetStats(masterAccountID string) ([]AssetStat, error) {
	var result []AssetStat
	stmt := sq.Select(
		"asset",
		`sum((amount+locked)/1000000) as hundreds`,
		`sum((amount+locked)::decimal/10000 % 100)::int as ones`,
		`((sum((amount+locked)::decimal/10000 % 100) % 1) * 10000)::bigint as remainder`,
	).Where("NOT account_id = ?", masterAccountID).From("balance").GroupBy("asset")

	err := q.Select(&result, stmt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return result, err
}

func (q *Q) CoinsInCirculation(masterAccountID string) ([]AssetAmount, error) {
	var result []AssetAmount
	stmt := selectCoinsInCirculationAmounts.Where("account_id != ?", masterAccountID)
	err := q.Select(&result, stmt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return result, err
}

func (q *Q) MustCoinsInCirculationForAsset(masterAccountID, asset string) (AssetAmount, error) {
	var result AssetAmount
	query := selectCoinsInCirculationAmounts.Where("account_id != ?", masterAccountID).
		Where("asset = ?", asset)

	err := q.Get(&result, query)
	return result, err
}

var selectBalanceAmounts = sq.Select(
	"b.asset as asset, b.amount as amount").
	From("balance b")

func (q *Q) AvailableEmissions(masterAccountID string) ([]AssetAmount, error) {
	var result []AssetAmount
	stmt := selectBalanceAmounts.Where("account_id = ?", masterAccountID)
	err := q.Select(&result, stmt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return result, err
}
