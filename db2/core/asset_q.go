package core

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// AssetQI is a helper interface to aid in configuring queries that loads
// slices or entry of Asset structs.
type AssetQI interface {
	// ByCode - returns asset by code, if not found returns nil, nil
	ByCode(assetCode string) (*Asset, error)
	// ForOwner - filters assets by owner
	ForOwner(owner string) AssetQI
	// Select - selects assets for specified filter
	Select() ([]Asset, error)
}

// assetQ is a helper struct to aid in configuring queries that loads
// slices or entry of Asset structs.
type assetQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

// ByCode - returns asset by code, if not found returns nil, nil
func (q *assetQ) ByCode(code string) (*Asset, error) {
	sql := selectAsset.Where("code = ?", code)
	var result Asset
	err := q.parent.Get(&result, sql)
	if q.parent.NoRows(err) {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to load asset", map[string]interface{}{
			"code": code,
		})
	}
	return &result, err
}

// ForOwner - filters assets by owner
func (q *assetQ) ForOwner(ownerID string) (AssetQI) {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("a.owner = ?", ownerID)
	return q
}
// Select - selects assets for specified filter
func (q *assetQ) Select() ([]Asset, error) {
	if q.Err != nil {
		return nil, errors.Wrap(q.Err, "failed before select. builder had error")
	}

	var result []Asset
	err := q.parent.Select(&result, q.sql)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select assets")
	}

	return result, nil
}

var selectAsset = sq.Select("a.code, a.owner, a.name, a.preissued_asset_signer, a.description, a.external_resource_link," +
	"a.max_issuance_amount, a.available_for_issueance, a.issued, a.policies").From("asset a")
