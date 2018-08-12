package history

import (
	sq "github.com/lann/squirrel"
		"time"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/swarmfund/horizon/db2"
)

type ContractsQI interface {
	ByStartTime(seconds int64) ContractsQI
	ByEndTime(seconds int64) ContractsQI
	ByDisputeState(isDisputing bool) ContractsQI
	Page(page db2.PageQuery) ContractsQI
	Select() ([]Contract, error)
}

type ContractsQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

func (q *Q) Contracts() ContractsQI {
	return &ContractsQ{
		parent: q,
		sql:    selectLedgerChanges,
	}
}

func (q *ContractsQ) ByStartTime(seconds int64) ContractsQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("start_time >= ?", time.Unix(seconds, 0).UTC())

	return q
}

func (q *ContractsQ) ByEndTime(seconds int64) ContractsQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("end_time >= ?", time.Unix(seconds, 0).UTC())

	return q
}

func (q *ContractsQ) ByDisputeState(isDisputing bool) ContractsQI {
	if q.Err != nil {
		return q
	}

	disputeState := int32(xdr.ContractStateDisputing)

	if isDisputing {
		q.sql = q.sql.Where("state & ? = ?", disputeState, disputeState)
	} else {
		q.sql = q.sql.Where("state & ? = 0", disputeState)
	}

	return q
}

func (q *ContractsQ) Page(page db2.PageQuery) ContractsQI {
	if q.Err != nil {
		return q
	}

	q.sql, q.Err = page.ApplyTo(q.sql, "contract_id")
	return q
}

func (q *ContractsQ) Select() ([]Contract, error) {
	if q.Err != nil {
		return nil, q.Err
	}

	var result []Contract
	q.Err = q.parent.Select(&result, q.sql)
	return result, q.Err
}