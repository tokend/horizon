package history

import (
	"time"

	sq "github.com/lann/squirrel"
	"gitlab.com/swarmfund/horizon/db2"
	"gitlab.com/tokend/go/xdr"
)

var selectContracts = sq.Select(
	"hc.id",
	"hc.contractor",
	"hc.customer",
	"hc.escrow",
	"hc.disputer",
	"hc.start_time",
	"hc.end_time",
	"hc.details",
	"hc.invoices",
	"hc.dispute_reason",
	"hc.state",
).From("history_contracts hc")

type ContractsQI interface {
	ByStartTime(seconds int64) ContractsQI
	ByEndTime(seconds int64) ContractsQI
	ByDisputeState(isDisputing bool) ContractsQI
	ByContractorID(contractorID string) ContractsQI
	ByCustomerID(customerID string) ContractsQI
	Page(page db2.PageQuery) ContractsQI
	Select() ([]Contract, error)
	ByID(id int64) (Contract, error)
	Update(contract Contract) error
}

type ContractsQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

func (q *Q) Contracts() ContractsQI {
	return &ContractsQ{
		parent: q,
		sql:    selectContracts,
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

	q.sql, q.Err = page.ApplyTo(q.sql, "id")
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

func (q *ContractsQ) ByID(id int64) (Contract, error) {
	if q.Err != nil {
		return Contract{}, q.Err
	}

	q.sql = q.sql.Where(sq.Eq{"id": id})

	var result Contract
	q.Err = q.parent.Get(&result, q.sql)
	return result, q.Err
}

func (q *ContractsQ) ByContractorID(contractorID string) ContractsQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where(sq.Eq{"contractor": contractorID})

	return q
}

func (q *ContractsQ) ByCustomerID(customerID string) ContractsQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where(sq.Eq{"customer": customerID})

	return q
}

// Update - update contract using it's ID
func (q *ContractsQ) Update(contract Contract) error {
	if q.Err != nil {
		return q.Err
	}

	query := sq.Update("history_contracts").SetMap(map[string]interface{}{
		"contractor":     contract.Contractor,
		"customer":       contract.Customer,
		"escrow":         contract.Escrow,
		"disputer":       contract.Disputer,
		"start_time":     contract.StartTime,
		"end_time":       contract.EndTime,
		"details":        contract.Details,
		"invoices":       contract.Invoices,
		"dispute_reason": contract.DisputeReason,
		"state":          contract.State,
	}).Where("id = ?", contract.ID)

	_, err := q.parent.Exec(query)
	return err
}
