package history2

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

type SaleConvertedBalancesQ struct {
	repo     *db2.Repo
	selector sq.SelectBuilder
}

func NewSaleConvertedBalancesQ(repo *db2.Repo) SaleConvertedBalancesQ {
	return SaleConvertedBalancesQ{
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
			From("participant_effects pe").
			Join("accounts a ON pe.account_id = a.id").
			Join("sale s ON (pe.effect#>>'{matched,order_book_id}')::int = s.id").
			Where("(pe.effect#>>'{type}')::int = ?", EffectTypeMatched).
			Where("pe.effect#>>'{matched,charged,asset_code}' = s.base_asset"),
	}
}

func (q SaleConvertedBalancesQ) FilterByParticipant(id string) SaleConvertedBalancesQ {
	q.selector = q.selector.Where("a.address = ?", id)
	return q
}

func (q SaleConvertedBalancesQ) Select() ([]SaleParticipation, error) {
	var result []SaleParticipation
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load sale participations")
	}

	return result, nil
}
