package core2

import (
	"github.com/lann/squirrel"
	"gitlab.com/tokend/horizon/db2"
)

type KeyValueQ struct {
	repo     *db2.Repo
	selector squirrel.SelectBuilder
}

func NewKeyValueQ(repo *db2.Repo) *KeyValueQ {
	return &KeyValueQ{
		repo,
		squirrel.Select("key", "value").From("key_value_entry"),
	}
}

func (q *KeyValueQ) ByKey(key string) (*KeyValue, error) {
	var result KeyValue
	stmt := q.selector.Where("key = ?", key)
	err := q.repo.Get(&result, stmt)
	if err != nil {
		if q.repo.NoRows(err) {
			return nil, nil
		}

		return nil, err
	}

	return &result, nil
}

func (q *KeyValueQ) Page(params *db2.OffsetPageParams) *KeyValueQ {
	q.selector = params.ApplyTo(q.selector, "key")
	return q
}

func (q *KeyValueQ) Select() ([]KeyValue, error) {
	var result []KeyValue
	err := q.repo.Select(&result, q.selector)
	return result, err
}
