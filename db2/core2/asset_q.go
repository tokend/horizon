package core2

import sq "github.com/lann/squirrel"

/*
import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

// AssetsQ is a helper struct to aid in configuring queries that loads
// asset structs.
type AssetsQ struct {
	repo *db2.Repo
}

// NewAssetsQ - creates new instance of AssetsQ
func NewAssetsQ(repo *db2.Repo) *AssetsQ {
	return &AssetsQ{
		repo: repo,
	}
}

// ByCode loads a row from `assets`, by code
// returns nil, nil - if asset does not exists
func (q *AssetsQ) ByCode(code string) (*Account, error) {
	var result Account
	err := q.repo.Get(&result, assetsFilter{Query: assetSelector}.ByCode(code).Query)
	if err != nil {
		if q.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load asset by address", logan.F{"address": address})
	}

	return &result, nil
}*/

var assetColumns = []string{"assets.code", "assets.owner", "assets.preissued_asset_signer", "assets.details",
	"assets.max_issuance_amount", "assets.available_for_issueance", "assets.issued",
	"assets.pending_issuance", "assets.policies", "assets.trailing_digits"}
var assetSelector = sq.Select(assetColumns...).From("asset assets")
