package core

import(
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"database/sql"
	"fmt"
)

var selectLimitsV2 = sq.Select("lim.id",
									   "lim.account_type",
									   "lim.account_id",
									   "lim.stats_op_type",
									   "lim.asset_code",
									   "lim.is_convert_needed",
									   "lim.daily_out",
									   "lim.weekly_out",
									   "lim.monthly_out",
									   "lim.annual_out").From("limits_v2 lim")
type LimitsV2QI interface {
	ForAccountByAccountType(accountID string, accountType int32) ([]LimitsV2Entry, error)
	ForAccountByStatsOpType(statsOpType int32, accountID string) ([]LimitsV2Entry, error)
	ForAccountTypeByStatsOpType(statsOpType, accountType int32) ([]LimitsV2Entry, error)
	ForAccount(accountID string) ([]LimitsV2Entry, error)
	ForAccountType(accountType int32) ([]LimitsV2Entry, error)
	Select(statsOpType int32) ([]LimitsV2Entry, error)
}

type LimitsV2Q struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

func (q *LimitsV2Q) ForAccountByAccountType(accountID string, accountType int32) ([]LimitsV2Entry, error) {
	if q.Err != nil {
		return nil, q.Err
	}

	accountIDStr := fmt.Sprintf("'%s'", accountID)

	query := fmt.Sprintf("select distinct on (account_id, account_type)  id, " +
		"account_type, account_id, stats_op_type, asset_code, is_convert_needed, daily_out, " +
		"weekly_out, monthly_out, annual_out " +
		"from limits_v2 " +
		"where (account_type = %d or account_type is null) and (account_id = %s or account_id is null) " +
		"order by account_id, account_type, stats_op_type desc", accountType, accountIDStr)

	var result []LimitsV2Entry
	err := q.parent.SelectRaw(&result, query)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "Failed to load limits_v2 for account")
	}

	return result, nil
}

func (q *LimitsV2Q) ForAccount(accountID string) ([]LimitsV2Entry, error) {
	if q.Err != nil {
		return nil, q.Err
	}

	query := selectLimitsV2.Where("lim.account_id = ?", accountID)
	var result []LimitsV2Entry
	err := q.parent.Select(&result, query)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "Failed to load limits_v2 for account")
	}

	return result, nil
}

func (q *LimitsV2Q) ForAccountByStatsOpType(statsOpType int32, accountID string) ([]LimitsV2Entry, error) {
	if q.Err != nil {
		return nil, q.Err
	}

	query := selectLimitsV2.Where("lim.account_id = ? AND lim.stats_op_type = ?", accountID, statsOpType)
	var result []LimitsV2Entry
	err := q.parent.Select(&result, query)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "Failed to load limits_v2 for account")
	}
	if result != nil {
		return result, nil
	}

	defaultLimits := new(LimitsV2Entry)
	defaultLimits.SetDefaultLimits()
	defaultLimits.AccountId = &accountID
	defaultLimits.StatsOpType = statsOpType
	result = append(result, *defaultLimits)

	return result, nil
}

func (q *LimitsV2Q) ForAccountType(accountType int32) ([]LimitsV2Entry, error) {
	if q.Err != nil {
		return nil, q.Err
	}

	var result []LimitsV2Entry
	query := selectLimitsV2.Where("lim.account_type = ?", accountType)
	err := q.parent.Select(&result, query)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "Failed to load limits_v2 for account type")
	}

	return result, nil
}

func (q *LimitsV2Q) ForAccountTypeByStatsOpType(statsOpType, accountType int32) ([]LimitsV2Entry, error){
	if q.Err != nil {
		return nil, q.Err
	}

	var result []LimitsV2Entry
	query := selectLimitsV2.Where("lim.account_type = ? AND lim.stats_op_type = ?", accountType, statsOpType)
	err := q.parent.Select(&result, query)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "Failed to load limits_v2 for account type by stats op type")
	}

	if result != nil {
		return result, nil
	}

	defaultLimits := new(LimitsV2Entry)
	defaultLimits.SetDefaultLimits()
	defaultLimits.AccountType = &accountType
	defaultLimits.StatsOpType = statsOpType
	result = append(result, *defaultLimits)

	return result, nil
}

func (q *LimitsV2Q) Select(statsOpType int32) ([]LimitsV2Entry, error){
	if q.Err != nil {
		return nil, q.Err
	}

	var result []LimitsV2Entry
	query := selectLimitsV2.Where("lim.stats_op_type = ?", statsOpType)
	err := q.parent.Select(&result, query)
	if q.parent.NoRows(err) {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "Failed to load limits_v2 by stats op type")
	}
	if result != nil {
		return result, nil
	}

	defaultLimits := new(LimitsV2Entry)
	defaultLimits.SetDefaultLimits()
	defaultLimits.StatsOpType = statsOpType
	result = append(result, *defaultLimits)

	return result, nil
}