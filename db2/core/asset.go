package core

import (
	sq "github.com/lann/squirrel"
)

type Asset struct {
	Code       string         `db:"code"`
	Policies   int32          `db:"policies"`
}

func (a *Asset) IsVisibleForUser(account *Account) bool {
	return true
}

func (q *Q) Assets() ([]Asset, error) {
	sql := selectAsset
	var assets []Asset
	err := q.Select(&assets, sql)
	return assets, err
}

func (q *Q) AssetByCode(code string) (*Asset, error) {
	sql := selectAsset.Where("code = ?", code)
	var result Asset
	err := q.Get(&result, sql)
	if q.NoRows(err) {
		return nil, nil
	}

	return &result, err
}

var selectAsset = sq.Select("a.code, a.policies").From("asset a")
