package history2

import (
	"database/sql"
	"github.com/lann/builder"

	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

// ParticipantEffectsQ - helper struct to get participants from db
type ParticipantEffectsQ struct {
	repo     *pgdb.DB
	selector sq.SelectBuilder
}

// NewParticipantEffectsQ - creates new ParticipantEffectsQ
func NewParticipantEffectsQ(repo *pgdb.DB) ParticipantEffectsQ {
	return ParticipantEffectsQ{
		repo: repo,
		selector: sq.Select("effects.id", "effects.account_id", "effects.balance_id", "effects.asset_code",
			"effects.effect", "effects.operation_id").From("participant_effects effects"),
	}
}

// WithOperation - left joins operations
func (q ParticipantEffectsQ) WithOperation() ParticipantEffectsQ {
	q.selector = q.selector.Columns(db2.GetColumnsForJoin(operationColumns, "operations")...).
		LeftJoin("operations operations ON effects.operation_id = operations.id")
	return q
}

// WithBalance - left joins balances
func (q ParticipantEffectsQ) WithBalance() ParticipantEffectsQ {
	q.selector = q.selector.Columns("balances.address balance_address").
		LeftJoin("balances balances ON balances.id = effects.balance_id")
	return q
}

// WithAccount - left joins accounts
func (q ParticipantEffectsQ) WithAccount() ParticipantEffectsQ {
	q.selector = q.selector.Columns("accounts.address account_address").
		LeftJoin("accounts accounts ON accounts.id = effects.account_id")
	return q
}

// ForBalance - adds filter by balance ID
func (q ParticipantEffectsQ) ForBalance(id ...uint64) ParticipantEffectsQ {
	q.selector = q.selector.Where(sq.Eq{
		"effects.balance_id": id,
	})
	return q
}

// ForEffect - adds filter by effectType
func (q ParticipantEffectsQ) ForEffect(types ...EffectType) ParticipantEffectsQ {
	q.selector = q.selector.Where(sq.Eq{"(effect->>'type')::integer": types})
	return q
}

// ForAsset - adds filter by asset
func (q ParticipantEffectsQ) ForAsset(asset string) ParticipantEffectsQ {
	q.selector = q.selector.Where("effects.asset_code = ?", asset)
	return q
}

// Movements - filters out non movement effects
func (q ParticipantEffectsQ) Movements() ParticipantEffectsQ {
	q.selector = q.selector.Where("effects.balance_id is not null")
	return q
}

// ForAccount - adds filter by accounts ID
func (q ParticipantEffectsQ) ForAccount(id uint64) ParticipantEffectsQ {
	q.selector = q.selector.Where("effects.account_id = ?", id)
	return q
}

func (q ParticipantEffectsQ) FilterByID(ids ...uint64) ParticipantEffectsQ {
	q.selector = q.selector.Where(sq.Eq{"effects.id": ids})
	return q
}

// Page - apply cursor paging params to the query
func (q ParticipantEffectsQ) Page(pageParams pgdb.CursorPageParams) ParticipantEffectsQ {
	q.selector = pageParams.ApplyTo(q.selector, "effects.id")
	return q
}

// PageOffset - apply offset paging params to the query
func (q ParticipantEffectsQ) PageOffset(pageParams pgdb.OffsetPageParams) ParticipantEffectsQ {
	q.selector = pageParams.ApplyTo(q.selector, "effects.id")
	return q
}

// Count - return total number of records with applied filters
func (q ParticipantEffectsQ) Count() (uint64, error) {
	var result uint64

	// replace default select columns
	selector := builder.Delete(q.selector, "Columns").(sq.SelectBuilder)
	selector = selector.Columns("COUNT(*)")

	err := q.repo.Get(&result, selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}

		return 0, errors.Wrap(err, "failed to get count")
	}

	return result, nil
}

// Select - selects ParticipantEffect from db using specified filters. Returns nil, nil - if one does not exists
func (q ParticipantEffectsQ) Select() ([]ParticipantEffect, error) {
	var result []ParticipantEffect
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to select participant effects")
	}

	return result, nil
}

func (q ParticipantEffectsQ) Get() (*ParticipantEffect, error) {
	var result ParticipantEffect

	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load poll")
	}

	return &result, nil
}
