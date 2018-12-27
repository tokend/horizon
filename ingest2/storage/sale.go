package storage

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
	"gitlab.com/tokend/horizon/db2/history2"
)

// Sale is helper struct to operate with `sales`
type Sale struct {
	repo *db2.Repo
}

// NewSale - creates new instance of the `Sale`
func NewSale(repo *db2.Repo) *Sale {
	return &Sale{
		repo: repo,
	}
}

// Insert - inserts new sale
func (q *Sale) Insert(sale history2.Sale) error {
	sql := sq.Insert("sale").
		Columns(
			"id", "owner_id", "base_asset", "default_quote_asset", "start_time", "end_time",
			"quote_assets", "soft_cap", "hard_cap", "details", "state", "base_current_cap",
			"base_hard_cap, sale_type",
		).
		Values(
			sale.ID, sale.OwnerAddress, sale.BaseAsset, sale.DefaultQuoteAsset, sale.StartTime, sale.EndTime,
			sale.QuoteAssets, sale.SoftCap, sale.HardCap, sale.Details, sale.State,
			sale.BaseCurrentCap, sale.BaseHardCap, sale.SaleType,
		)

	_, err := q.repo.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to insert sale", logan.F{"sale_id": sale.ID})
	}

	return nil
}

// Update - updates existing sale
func (q *Sale) Update(sale history2.Sale) error {
	sql := sq.Update("sale").SetMap(map[string]interface{}{
		"owner_id":            sale.OwnerAddress,
		"base_asset":          sale.BaseAsset,
		"default_quote_asset": sale.DefaultQuoteAsset,
		"start_time":          sale.StartTime,
		"end_time":            sale.EndTime,
		"quote_assets":        sale.QuoteAssets,
		"soft_cap":            sale.SoftCap,
		"hard_cap":            sale.HardCap,
		"details":             sale.Details,
		"state":               sale.State,
		"base_hard_cap":       sale.BaseHardCap,
		"base_current_cap":    sale.BaseCurrentCap,
		"sale_type":           sale.SaleType,
	}).Where("id = ?", sale.ID)

	_, err := q.repo.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to update sale", logan.F{"sale_id": sale.ID})
	}

	return nil
}

// SetState - sets state
func (q *Sale) SetState(saleID uint64, state history2.SaleState) error {
	sql := sq.Update("sale").Set("state", state).Where("id = ?", saleID)
	_, err := q.repo.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to set state", logan.F{"sale_id": saleID})
	}

	return nil
}
