package core2

import (
	"fmt"
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

var assetColumns = []string{"code", "owner", "preissued_asset_signer", "details",
	"max_issuance_amount", "available_for_issueance", "issued",
	"pending_issuance", "policies", "trailing_digits"}

func getAssetColumns(tableName string) []string {
	result := make([]string, 0, len(assetColumns))
	for _, column := range assetColumns {
		result = append(result, fmt.Sprintf(`%s.%s "%s.%s"`, tableName, column, tableName, column))
	}

	return result
}

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

func (q AssetsQ) FilterByOwner(ownerID string) AssetsQ {
	q.selector = q.selector.Where("assets.owner = ?", ownerID)
	return q
}

func (q AssetsQ) FilterByPolicy(mask uint64) AssetsQ {
	q.selector = q.selector.Where("assets.policies & ? = ?", mask, mask)
	return q
}

func (q AssetsQ) Page(limit, offset uint64) AssetsQ {
	q.selector = q.selector.Limit(limit).Offset(offset)
	return q
}

func (q AssetsQ) Select() ([]Asset, error) {
	var result []Asset
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		if q.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load assets")
	}

	return result, nil
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
