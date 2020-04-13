package storage

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
	"gitlab.com/tokend/horizon/db2/history2"
	regources "gitlab.com/tokend/regources/generated"
)

type Asset struct {
	repo *db2.Repo
}

func NewAsset(repo *db2.Repo) *Asset {
	return &Asset{
		repo: repo,
	}
}

func (q *Asset) Insert(asset history2.Asset) error {
	sql := sq.Insert("asset").
		Columns(
			"code",
			"owner",
			"preissued_asset_signer",
			"details.max_issuance_amount",
			"available_for_issueance",
			"issued",
			"pending_issuance",
			"policies",
			"trailing_digits",
			"type",
			"state",
		).
		Values(
			asset.Code,
			asset.Owner,
			asset.PreIssuanceAssetSigner,
			asset.MaxIssuanceAmount,
			asset.PendingIssuance,
			asset.Details,
			asset.Issued,
			asset.PendingIssuance,
			asset.Policies,
			asset.TrailingDigits,
			asset.Type,
			asset.State,
		)

	_, err := q.repo.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to insert asset", logan.F{
			"asset_code":  asset.Code,
		})
	}

	return nil
}

func (q *Asset) SetState(code string, state regources.AssetState) error {
	sql := sq.Update("asset").Set("state", state).Where("code = ?", code)
	_, err := q.repo.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to set state", logan.F{"asset_code": code})
	}

	return nil
}

func (q *Asset) Update(asset history2.Asset) error{
	sql := sq.Update("asset").SetMap(
		map[string]interface{}{
			"owner":                       asset.Owner,
			"preissued_asset_signer":      asset.PreIssuanceAssetSigner,
			"details.max_issuance_amount": asset.Details,
			"available_for_issueance":     asset.AvailableForIssuance,
			"issued":                      asset.Issued,
			"pending_issuance":            asset.PendingIssuance,
			"policies":                    asset.Policies,
			"trailing_digits":             asset.TrailingDigits,
			"type":                        asset.Type,
			"state":                       asset.State,
		}).Where("code = ?", asset.Code)

	_, err := q.repo.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to update poll", logan.F{
			"asset_code": asset.Code,
		})
	}

	return nil
}
