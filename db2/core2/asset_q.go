package core2

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

var assetColumns = []string{"assets.code", "assets.owner", "assets.preissued_asset_signer", "assets.details",
	"assets.max_issuance_amount", "assets.available_for_issueance", "assets.issued",
	"assets.pending_issuance", "assets.policies", "assets.trailing_digits"}

//AssetQ - helper struct to load assets from db
type AssetQ struct {
	repo     *db2.Repo
	selector sq.SelectBuilder
}

func NewAssetQ(repo *db2.Repo) AssetQ {
	return AssetQ{
		repo:     repo,
		selector: sq.Select(assetColumns...).From("asset AS assets"),
	}
}

func (q AssetQ) GetByCode(code string) (*Asset, error) {
	return q.FilterByCode(code).Get()
}

func (q AssetQ) FilterByCode(code string) AssetQ {
	q.selector = q.selector.Where("assets.code = ?")
	return q
}

func (q AssetQ) Get() (*Asset, error) {
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
