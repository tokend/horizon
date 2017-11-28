package core

import (
	"database/sql"
	sq "github.com/lann/squirrel"
	"math"
)

type Limits struct {
	DailyOut   int64 `db:"daily_out"`
	WeeklyOut  int64 `db:"weekly_out"`
	MonthlyOut int64 `db:"monthly_out"`
	AnnualOut  int64 `db:"annual_out"`
}

func DefaultLimits() Limits {
	return Limits{
		DailyOut:   math.MaxInt64,
		WeeklyOut:  math.MaxInt64,
		MonthlyOut: math.MaxInt64,
		AnnualOut:  math.MaxInt64,
	}
}

type AccountLimits struct {
	Accountid string `db:"accountid"`
	Limits
}

// SignersByAddress loads all signer rows for `addy`
func (q *Q) LimitsByAddress(addy string) (*AccountLimits, error) {
	query := selectLimit.Where("accountid = ?", addy)
	var result AccountLimits
	err := q.Get(&result, query)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (q *Q) LimitsForAccount(accountID string, accountType int32) (Limits, error) {
	accountLimits, err := q.LimitsByAddress(accountID)
	if err != nil {
		return Limits{}, err
	}

	if accountLimits != nil {
		return accountLimits.Limits, nil
	}

	accountTypeLimits, err := q.LimitsByAccountType(accountType)
	if err != nil {
		return Limits{}, err
	}

	if accountTypeLimits != nil {
		return accountTypeLimits.Limits, nil
	}

	return DefaultLimits(), nil
}

var selectLimit = sq.Select(
	"li.daily_out",
	"li.weekly_out",
	"li.monthly_out",
	"li.annual_out",
).From("account_limits li").OrderBy("accountid DESC")
