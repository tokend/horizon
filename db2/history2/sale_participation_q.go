package history2

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

// SaleParticipationQ is a helper struct to aid in configuring queries that loads
// sale participation structures.
type SaleParticipationQ struct {
	repo     *db2.Repo
	selector sq.SelectBuilder
}

// NewSaleParticipationQ - creates new instance of SaleParticipationQ
func NewSaleParticipationQ(repo *db2.Repo) SaleParticipationQ {
	return SaleParticipationQ{
		repo: repo,
		selector: sq.Select(
			"sp.id",
			"sp.participant_id",
			"sp.sale_id",
			"sp.base_amount",
			"sp.quote_amount",
			"sp.base_asset",
			"sp.quote_asset",
			"sp.price",
		).From("sale_participation sp"),
	}
}

// FilterByID - returns q with filter by sale ID
func (q SaleParticipationQ) FilterByID(id uint64) SaleParticipationQ {
	q.selector = q.selector.Where("sp.id = ?", id)
	return q
}

// GetByID loads a row from `sales`, by ID
// returns nil, nil - if sale does not exists
func (q SaleParticipationQ) GetByID(id uint64) (*SaleParticipation, error) {
	return q.FilterByID(id).Get()
}

// FilterByOwner - returns q with filter by participant
func (q SaleParticipationQ) FilterByParticipant(address string) SaleParticipationQ {
	q.selector = q.selector.Where("sp.participant_id = ?", address)
	return q
}

// FilterByBaseAsset - returns q with filter by base asset
func (q SaleParticipationQ) FilterByBaseAsset(asset string) SaleParticipationQ {
	q.selector = q.selector.Where("sp.base_asset = ?", asset)
	return q
}

// FilterByBaseAsset - returns q with filter by base asset
func (q SaleParticipationQ) FilterByQuoteAsset(asset string) SaleParticipationQ {
	q.selector = q.selector.Where("sp.quote_asset = ?", asset)
	return q
}

// FilterBySaleID- returns q with filter by sale id
func (q SaleParticipationQ) FilterBySale(saleID uint64) SaleParticipationQ {
	q.selector = q.selector.Where("sp.sale_id = ?", saleID)
	return q
}

// Page - returns Q with specified limit and offset params
func (q SaleParticipationQ) Page(params db2.CursorPageParams) SaleParticipationQ {
	q.selector = params.ApplyTo(q.selector, "sp.id")
	return q
}

// Get - loads a row from `sales`
// returns nil, nil - if sale does not exists
// returns error if more than one Sale found
func (q SaleParticipationQ) Get() (*SaleParticipation, error) {
	var result SaleParticipation
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if q.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load sale participation")
	}

	return &result, nil
}

//Select - selects slice from the db, if no sales found - returns nil, nil
func (q SaleParticipationQ) Select() ([]SaleParticipation, error) {
	var result []SaleParticipation
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load sale participation")
	}

	return result, nil
}
