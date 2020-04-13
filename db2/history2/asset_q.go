package history2

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

var assetColumns = []string{
	"asset.code",
	"asset.owner",
	"asset.preissued_asset_signer",
	"asset.details",
	"asset.max_issuance_amount",
	"asset.available_for_issuance",
	"asset.issued",
	"asset.pending_issuance",
	"asset.policies",
	"asset.trailing_digits",
	"asset.type",
	"asset.state",
}

// AssetQ is a helper struct to aid in configuring queries that loads assets
type AssetQ struct {
	repo     *db2.Repo
	selector sq.SelectBuilder
}

// NewAssetQ- creates new instance of AssetQ
func NewAssetQ(repo *db2.Repo) AssetQ {
	return AssetQ{
		repo:     repo,
		selector: sq.Select(assetColumns...).From("asset asset"),
	}
}

// GetByCode - get asset by code
func (q AssetQ) GetByCode(code string) (*Asset, error) {
	q.selector = q.selector.Where("asset.code = ?", code)
	return q.Get()
}

//GetByOwner - gets account by owner address, returns nil, nil if one does not exist
func (q AssetQ) ByOwner(address string) (*Asset, error) {
	q.selector = q.selector.Where("accounts.owner = ?", address)
	return q.Get()
}

//Get - selects account from db, returns nil, nil if one does not exists
func (q AssetQ) Get() (*Asset, error) {
	var result Asset
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if q.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load account")
	}

	return &result, nil
}
