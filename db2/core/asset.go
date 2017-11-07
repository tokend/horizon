package core

import (
	"bullioncoin.githost.io/development/go/xdr"
	sq "github.com/lann/squirrel"
)

type Asset struct {
	Code       string         `db:"code"`
	Policies   int32          `db:"policies"`
	Token      string         `db:"token"`
	AssetForms xdr.AssetForms `db:"forms"`
}

func (a *Asset) IsVisibleForUser(account *Account) bool {
	isAccountAllowedToSeeTokens := account.Policies&int32(xdr.AccountPoliciesAllowToTransferTokens) == int32(xdr.AccountPoliciesAllowToTransferTokens)
	return isAccountAllowedToSeeTokens || a.Code == "XAAU" || a.Code == "USD" || a.Code == "XAAG"
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

var selectAsset = sq.Select("a.code, a.policies, a.token, a.forms").From("asset a")
