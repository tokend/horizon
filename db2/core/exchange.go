package core

import (
	sq "github.com/lann/squirrel"
)

type Exchange struct {
	AccountID     string `db:"account_id"`
	Name          string `db:"name"`
	RequireReview bool   `db:"require_review"`
}

func (q *Q) ExchangeName(address string) (*string, error) {
	var result []byte
	sql := sq.Select("name").
		From("exchanges ex").
		Limit(1).
		Where("account_id = ?", address)
	err := q.Get(&result, sql)

	if q.NoRows(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	name := string(result[:])
	return &name, nil
}
