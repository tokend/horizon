package db2

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const (
	// Default limit is the default amount of records will be returned
	DefaultLimit  = 15
	DefaultOffset = 0
)

// NewPageQueryV2 accepts the page derived for request
// URL params and calculates the proper offset for the
// sql-select
func NewPageQueryV2(page uint64) (result PageQueryV2, err error) {
	if err != nil {
		err = errors.Wrap(err, "Got invalid page")
		return
	}

	result = PageQueryV2{
		Limit:  DefaultLimit,
		Offset: page * DefaultLimit,
		Page: page,
	}

	return
}

// ApplyTo returns a new SelectBuilder after applying the paging effects of
// `p` to `sql`.  This method provides the default case for paging: int64
// cursor-based paging by an id column.
func (p PageQueryV2) ApplyTo(sql sq.SelectBuilder) (sq.SelectBuilder, error) {
	return sql.Limit(p.Limit).Offset(p.Offset), nil
}
