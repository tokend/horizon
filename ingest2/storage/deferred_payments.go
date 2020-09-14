package storage

import (
	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2/history2"
)

// CreateDeferredPayment is helper struct to operate with `deferredPayments`
type DeferredPayment struct {
	repo    *pgdb.DB
	q       history2.DeferredPaymentQ
	updater sq.UpdateBuilder
}

// NewDeferredPayment - creates new instance of the `CreateDeferredPayment`
func NewDeferredPayment(repo *pgdb.DB) *DeferredPayment {
	return &DeferredPayment{
		repo:    repo,
		updater: sq.Update("deferred_payments"),
	}
}

// Insert - inserts new deferredPayment
func (q *DeferredPayment) Insert(deferredPayment history2.DeferredPayment) error {

	sql := sq.Insert("deferredPayment").
		Columns(
			"id",
			"amount",
			"source_pays_for_dest",
			"source_fixed_fee",
			"source_percent_fee",
			"destination_fixed_fee",
			"destination_percent_fee",
			"source_account",
			"source_balance",
			"destination_account",
		).
		Values(
			deferredPayment.ID, deferredPayment.Amount, deferredPayment.SourcePaysForDest, deferredPayment.SourceFixedFee,
			deferredPayment.SourcePercentFee, deferredPayment.DestinationFixedFee, deferredPayment.DestinationPercentFee,
			deferredPayment.SourceAccount, deferredPayment.SourceBalance, deferredPayment.DestinationAccount,
		)

	err := q.repo.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to insert deferredPayment", logan.F{"deferredPayment_id": deferredPayment.ID})
	}

	return nil
}

// Update - updates existing deferredPayment
func (q *DeferredPayment) Update(deferredPayment history2.DeferredPayment) error {
	sql := sq.Update("deferred_payments").SetMap(map[string]interface{}{
		"amount":                  deferredPayment.Amount,
		"source_pays_for_dest":    deferredPayment.SourcePaysForDest,
		"source_fixed_fee":        deferredPayment.SourceFixedFee,
		"source_percent_fee":      deferredPayment.SourcePercentFee,
		"destination_fixed_fee":   deferredPayment.DestinationFixedFee,
		"destination_percent_fee": deferredPayment.DestinationPercentFee,
		"source_account":          deferredPayment.SourceAccount,
		"source_balance":          deferredPayment.SourceBalance,
		"destination_account":     deferredPayment.DestinationAccount,
	}).Where(sq.Eq{"id": deferredPayment.ID})

	err := q.repo.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to update deferredPayment", logan.F{"deferredPayment_id": deferredPayment.ID})
	}

	return nil
}

// Remove - removes existing deferredPayment
func (q *DeferredPayment) Remove(deferredPaymentID int64) error {
	sql := sq.Delete("deferred_payments").Where(sq.Eq{"id": deferredPaymentID})

	err := q.repo.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to remove deferredPayment", logan.F{"deferredPayment_id": deferredPaymentID})
	}

	return nil
}

func (q *DeferredPayment) MustDeferredPayment(id int64) history2.DeferredPayment {
	deferredPayment, err := q.getDeferredPayment(id)
	if err != nil {
		panic(err)
	}

	return deferredPayment
}

func (q *DeferredPayment) getDeferredPayment(id int64) (history2.DeferredPayment, error) {
	deferredPayment, err := q.q.GetByID(id)
	if err != nil {
		return history2.DeferredPayment{}, errors.Wrap(err, "failed to get deferredPayment by id", logan.F{
			"deferredPayment_id": id,
		})
	}

	if deferredPayment == nil {
		return history2.DeferredPayment{}, errors.From(errors.New("deferredPayment missing"), logan.F{
			"deferredPayment_id": id,
		})
	}

	return *deferredPayment, nil
}
