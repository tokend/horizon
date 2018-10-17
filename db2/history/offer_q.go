package history

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

var selectOffers = sq.Select(
	"id",
	"offer_id",
	"owner_id",
	"base_asset",
	"quote_asset",
	"is_buy",
	"initial_base_amount",
	"current_base_amount",
	"price",
	"is_canceled",
	"created_at").
	From("history_offer")

type OffersQI interface {
	ForBase(base string) OffersQI
	ForQuote(quote string) OffersQI
	ForOwnerID(ownerID string) OffersQI
	ForIsBuy(isBuy bool) OffersQI
	NoMatches() OffersQI
	PartiallyMatched() OffersQI
	FullyMatched() OffersQI
	Canceled() OffersQI
	Page(page db2.PageQuery) OffersQI
	Select() ([]Offer, error)
	Insert(offer Offer) error
	UpdateBaseAmount(baseAmount, offerID int64) error
	Cancel(offerID int64) error
}

type OffersQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

func (q *Q) Offers() OffersQI {
	return &OffersQ{
		parent: q,
		sql:    selectOffers,
	}
}

func (q *OffersQ) ForBase(base string) OffersQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where(sq.Eq{"base_asset": base})
	return q
}

func (q *OffersQ) ForQuote(quote string) OffersQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where(sq.Eq{"quote_asset": quote})
	return q
}

func (q *OffersQ) ForOwnerID(ownerID string) OffersQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where(sq.Eq{"owner_id": ownerID})
	return q
}

func (q *OffersQ) ForIsBuy(isBuy bool) OffersQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where(sq.Eq{"is_buy": isBuy})
	return q
}

func (q *OffersQ) NoMatches() OffersQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("initial_base_amount = current_base_amount")
	return q
}

func (q *OffersQ) PartiallyMatched() OffersQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("initial_base_amount != current_base_amount").
		Where("current_base_amount > 0")
	return q
}

func (q *OffersQ) FullyMatched() OffersQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("current_base_amount = 0")
	return q
}

func (q *OffersQ) Canceled() OffersQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("is_canceled")
	return q
}

func (q *OffersQ) Page(page db2.PageQuery) OffersQI {
	if q.Err != nil {
		return q
	}

	q.sql, q.Err = page.ApplyTo(q.sql, "id")
	return q
}

func (q *OffersQ) Select() ([]Offer, error) {
	if q.Err != nil {
		return nil, errors.Wrap(q.Err, "failed to build query before select")
	}

	var offers []Offer
	err := q.parent.Select(&offers, q.sql)
	if q.parent.NoRows(err) {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to execute select query")
	}

	return offers, nil
}

func (q *OffersQ) Insert(offer Offer) error {
	query := sq.Insert("history_offer").SetMap(sq.Eq{
		"offer_id":            offer.OfferID,
		"owner_id":            offer.OwnerID,
		"base_asset":          offer.BaseAsset,
		"quote_asset":         offer.QuoteAsset,
		"is_buy":              offer.IsBuy,
		"initial_base_amount": offer.InitialBaseAmount,
		"current_base_amount": offer.CurrentBaseAmount,
		"price":               offer.Price,
		"is_canceled":         offer.IsCanceled,
		"created_at":          offer.CreatedAt,
	})

	_, err := q.parent.Exec(query)
	if err != nil {
		return errors.Wrap(err, "failed to execute insert query")
	}

	return nil
}

func (q *OffersQ) UpdateBaseAmount(baseAmount, offerID int64) error {
	query := sq.Update("history_offer").SetMap(sq.Eq{
		"current_base_amount": baseAmount,
	}).Where(sq.Eq{"offer_id": offerID})

	_, err := q.parent.Exec(query)
	if err != nil {
		return errors.Wrap(err, "failed to execute update query")
	}

	return nil
}

func (q *OffersQ) Cancel(offerID int64) error {
	query := sq.Update("history_offer").
		Set("is_canceled", true).
		Where(sq.Eq{"offer_id": offerID})

	_, err := q.parent.Exec(query)
	if err != nil {
		return errors.Wrap(err, "failed to execute update query")
	}

	return nil
}
