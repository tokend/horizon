package history

import (
	"gitlab.com/tokend/horizon/db2"
	sq "github.com/lann/squirrel"
)

var selectCoinsEmissionRequest = sq.Select("her.*").
	From("history_emission_requests her")

type CoinsEmissionRequestsQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

type CoinsEmissionRequestsQI interface {
	Page(page db2.PageQuery) CoinsEmissionRequestsQI
	ForAsset(asset string) CoinsEmissionRequestsQI
	ForExchange(aid string) CoinsEmissionRequestsQI
	ForReference(reference string) CoinsEmissionRequestsQI
	ForState(state *bool) CoinsEmissionRequestsQI
	Select(dest interface{}) error
}

// CoinsEmissionRequests provides a helper to filter rows from the `history_emission_requests` table
// with pre-defined filters.  See `CoinsEmissionRequestsQ` methods for the available filters.
func (q *Q) CoinsEmissionRequests() CoinsEmissionRequestsQI {
	return &CoinsEmissionRequestsQ{
		parent: q,
		sql:    selectCoinsEmissionRequest,
	}
}

func (q *Q) CoinsEmissionRequestByRequestID(dest interface{}, requestID string) error {
	sql := selectCoinsEmissionRequest.Where("her.request_id = ?", requestID)
	return q.Get(dest, sql)
}

func (q *Q) CoinsEmissionRequestByReference(dest interface{}, reference string) error {
	sql := selectCoinsEmissionRequest.Where("her.reference = ?", reference)
	return q.Get(dest, sql)
}

func (q *CoinsEmissionRequestsQ) ForAsset(asset string) CoinsEmissionRequestsQI {
	if q.Err != nil {
		return q
	}
	q.sql = q.sql.Where("her.asset = ?", asset)
	return q
}

func (q *CoinsEmissionRequestsQ) ForExchange(accountId string) CoinsEmissionRequestsQI {
	if q.Err != nil {
		return q
	}
	q.sql = q.sql.Where("her.issuer = ?", accountId)
	return q
}

func (q *CoinsEmissionRequestsQ) ForReference(reference string) CoinsEmissionRequestsQI {
	if q.Err != nil {
		return q
	}
	q.sql = q.sql.Where("her.reference = ?", reference)
	return q
}

func (q *CoinsEmissionRequestsQ) ForState(state *bool) CoinsEmissionRequestsQI {
	if q.Err != nil {
		return q
	}
	q.sql = q.sql.Where("her.approved = ?", state)
	return q
}

// Page specifies the paging constraints for the query being built by `q`.
func (q *CoinsEmissionRequestsQ) Page(page db2.PageQuery) CoinsEmissionRequestsQI {
	if q.Err != nil {
		return q
	}

	q.sql, q.Err = page.ApplyTo(q.sql, "her.request_id")
	return q
}

// Select loads the results of the query specified by `q` into `dest`.
func (q *CoinsEmissionRequestsQ) Select(dest interface{}) error {
	if q.Err != nil {
		return q.Err
	}

	q.Err = q.parent.Select(dest, q.sql)
	return q.Err
}
