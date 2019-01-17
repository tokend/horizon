package core2

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

var assetColumns = []string{"assets.code", "assets.owner", "assets.preissued_asset_signer", "assets.details",
	"assets.max_issuance_amount", "assets.available_for_issueance", "assets.issued",
	"assets.pending_issuance", "assets.policies", "assets.trailing_digits"}

//AssetsQ - helper struct to load assets from db
type AssetsQ struct {
	repo     *db2.Repo
	selector sq.SelectBuilder
}

func NewAssetsQ(repo *db2.Repo) AssetsQ {
	return AssetsQ{
		repo:     repo,
		selector: sq.Select(assetColumns...).From("asset AS assets"),
	}
}

func (q AssetsQ) GetByCode(code string) (*Asset, error) {
	return q.FilterByCode(code).Get()
}

func (q AssetsQ) FilterByCode(code string) AssetsQ {
	q.selector = q.selector.Where("assets.code = ?", code)
	return q
}

func (q AssetsQ) Get() (*Asset, error) {
	var result Asset
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if q.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load asset")
	}

	return &result, nil
}
