package history2

import (
	"database/sql"
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

var effectColumns = []string{
	"type",
	"funded",
	"issued",
	"charged",
	"withdrawn",
	"locked",
	"unlocked",
	"chargedFromLocked",
	"matched",
}

type EffectQ struct {
	repo     *pgdb.DB
	selector sq.SelectBuilder
}

// NewEffectQ - creates new instance of EffectQ
func NewEffectQ(repo *pgdb.DB) EffectQ {
	return EffectQ{
		repo: repo,
		selector: sq.Select(effectColumns...).From("effects effects"),
	}
}

// FilterByTypes - returns data by given EffectTypes
func (q EffectQ) FilterByTypes(types []int64) (*Effect, error) {
	q.selector = q.selector.Where("effects.address = ?", types)
	return q.Get()
}

//Get - selects effect from db, returns nil, nil if one does not exists
func (q EffectQ) Get() (*Effect, error) {
	var result Effect
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load account")
	}

	return &result, nil
}