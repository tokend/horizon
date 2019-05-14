package history2

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

// SaleParticipationQ is a helper struct to aid in configuring queries that load
// sale participation structures from `participant_effects` table`.
type SaleParticipationQ struct {
	repo     *db2.Repo
	selector sq.SelectBuilder
}

// NewSaleParticipationQ - creates new instance of SaleParticipationQ
func NewSaleParticipationQ(repo *db2.Repo) SaleParticipationQ {
	return SaleParticipationQ{
		repo: repo,
		selector: sq.Select(
			"(pe.effect#>>'{matched,offer_id}')::int id",
			"a.address participant_id",
			"(pe.effect#>>'{matched,order_book_id}')::int sale_id",
			"pe.effect#>>'{matched,charged,asset_code}' quote_asset",
			"pe.effect#>>'{matched,charged,amount}' quote_amount",
			"pe.effect#>>'{matched,funded,asset_code}' base_asset",
			"pe.effect#>>'{matched,funded,amount}' base_amount",
		).
			Distinct().
			From("participant_effects pe").
			Join("accounts a ON pe.account_id = a.id").
			Where("(pe.effect#>>'{type}')::int = ?", EffectTypeMatched).
			Where("(pe.effect#>>'{matched,offer_id}')::int != ?", 0),
	}
}

// FilterByQuoteAsset - returns q with filter by quote asset
func (q SaleParticipationQ) FilterByQuoteAsset(asset string) SaleParticipationQ {
	q.selector = q.selector.Where("pe.effect#>>'{matched,charged,asset_code}' = ?", asset)

	return q
}

// FilterByParticipant - returns q with filter by participant
func (q SaleParticipationQ) FilterByParticipant(id string) SaleParticipationQ {
	q.selector = q.selector.Where("a.address = ?", id)
	return q
}

// FilterBySaleID - returns q with filter by sale id
func (q SaleParticipationQ) FilterBySale(id uint64) SaleParticipationQ {
	q.selector = q.selector.Where("(pe.effect#>>'{matched,order_book_id}')::int = ?", id)
	return q
}

// Page - returns Q with specified cursor params
func (q SaleParticipationQ) Page(params db2.CursorPageParams) SaleParticipationQ {
	q.selector = params.ApplyTo(q.selector, "(pe.effect#>>'{matched,offer_id}')::int")
	return q
}

// Select - selects slice from the db, if no sales found - returns nil, nil
func (q SaleParticipationQ) Select() ([]SaleParticipation2, error) {
	var result []SaleParticipation2
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load sale participations")
	}

	return result, nil
}
