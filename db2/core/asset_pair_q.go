package core

import (
	sq "github.com/lann/squirrel"
)

// AssetPairsQ is a helper interface to aid in configuring queries that loads
// slices or entry of Asset Pair structs.
type AssetPairsQ interface {
	// ByCode returns nil, if asset pair not found
	ByCode(base, quote string) (*AssetPair, error)
	// Select - selects all assets pairs with specified filters
	Select() ([]AssetPair, error)
}

// assetQ is a helper struct to aid in configuring queries that loads
// slices or entry of Asset structs.
type assetPairQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

func (q *assetPairQ) Select() ([]AssetPair, error) {
	if q.Err != nil {
		return nil, q.Err
	}

	var assetPairs []AssetPair
	err := q.parent.Select(&assetPairs, q.sql)
	if q.parent.NoRows(err) {
		return nil, nil
	}

	return assetPairs, err
}

// returns nil, if not found
func (q *assetPairQ) ByCode(base, quote string) (*AssetPair, error) {
	sql := selectAssetPair.Where("base = ? AND quote = ?", base, quote)
	var result AssetPair
	err := q.parent.Get(&result, sql)
	if q.parent.Repo.NoRows(err) {
		return nil, nil
	}

	return &result, err
}

var selectAssetPair = sq.Select("a.base, a.quote, a.current_price, a.physical_price, a.physical_price_correction, a.max_price_step, a.policies").From("asset_pair a")
