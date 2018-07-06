package history

import (
	"time"

	"fmt"

	sq "github.com/lann/squirrel"
	"gitlab.com/swarmfund/horizon/db2"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/logan/v3"
)

var selectOperation = sq.Select("distinct on (ho.id) ho.*").
	From("history_operations ho").
	LeftJoin("history_operation_participants hop on hop.history_operation_id = ho.id")

// OperationsQ is a helper struct to aid in configuring queries that loads
// slices of Operation structs.
type OperationsQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

type OperationsQI interface {
	ForAccount(aid string) OperationsQI
	ForAccountType(accountType int32) OperationsQI
	ForBalance(bid string) OperationsQI
	ForAsset(asset string) OperationsQI
	ForTypes(opTypes []xdr.OperationType) OperationsQI
	ForTx(txhash string) OperationsQI
	ForReference(reference string) OperationsQI
	Since(ts time.Time) OperationsQI
	To(ts time.Time) OperationsQI
	// required prior to exchange and asset filtering
	JoinOnBalance() OperationsQI
	// JoinOnAccount required to filter on account ID and account type
	JoinOnAccount() OperationsQI
	Page(page db2.PageQuery) OperationsQI
	Update(op Operation) error
	Select(dest interface{}) error

	Participants(dest map[int64]*OperationParticipants) error
}

// Operations provides a helper to filter the operations table with pre-defined
// filters.  See `OperationsQ` for the available filters.
func (q *Q) Operations() OperationsQI {
	return &OperationsQ{
		parent: q,
		sql:    selectOperation,
	}
}

// OperationByID loads a single operation with `id` into `dest`
func (q *Q) OperationByID(dest interface{}, id int64) error {
	sql := selectOperation.
		Limit(1).
		Where("ho.id = ?", id)

	return q.Get(dest, sql)
}

func (q *OperationsQ) Participants(dest map[int64]*OperationParticipants) error {
	ids := []int64{}
	for opid, _ := range dest {
		ids = append(ids, opid)
	}

	if len(ids) == 0 {
		return nil
	}

	stmt := sq.Select("hop.history_operation_id, ha.address as account_id, hop.balance_id, hop.effects").
		From("history_operation_participants hop").
		LeftJoin("history_accounts ha on ha.id = hop.history_account_id"). // join to get account addresses
		Where(sq.Eq{"hop.history_operation_id": ids})

	var participants []*Participant
	err := q.parent.Select(&participants, stmt)
	if err != nil {
		return err
	}

	for _, participant := range participants {
		opid := participant.OperationID
		dest[opid].Participants = append(dest[opid].Participants, participant)
	}

	return nil
}

func (q *OperationsQ) JoinOnAccount() OperationsQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Join("history_accounts ha on ha.id = hop.history_account_id")

	return q
}

// ForAccount filters the operations collection to a specific account
func (q *OperationsQ) ForAccount(address string) OperationsQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("ha.address = ?", address)

	return q
}

func (q *OperationsQ) ForAccountType(accountType int32) OperationsQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("ha.account_type = ?", accountType)

	return q
}

func (q *OperationsQ) JoinOnBalance() OperationsQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.LeftJoin("history_balances hb on hb.balance_id=hop.balance_id")

	return q
}

func (q *OperationsQ) ForAsset(asset string) OperationsQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("hb.asset = ?", asset)

	return q
}

func (q *OperationsQ) ForBalance(balanceID string) OperationsQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("hop.balance_id = ?", balanceID)

	return q
}

func (q *OperationsQ) ForTypes(opTypes []xdr.OperationType) OperationsQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where(sq.Eq{"ho.type": opTypes})

	return q
}

func (q *OperationsQ) ForTx(txhash string) OperationsQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Join("history_transactions ht on ht.id=ho.transaction_id").
		Where("ht.transaction_hash = ?", txhash)

	return q
}

func (q *OperationsQ) ForReference(reference string) OperationsQI {
	if q.Err != nil {
		return q
	}

	// FIXME might(will) not work for all operation types, works at least for payments and issuances
	q.sql = q.sql.Where(fmt.Sprintf("ho.details->>'reference' ilike '%%%s%%'", reference))

	return q
}

func (q *OperationsQ) Since(ts time.Time) OperationsQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("ho.ledger_close_time >= ?", ts)

	return q
}

func (q *OperationsQ) To(ts time.Time) OperationsQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("ho.ledger_close_time <= ?", ts)

	return q
}

// Page specifies the paging constraints for the query being built by `q`.
func (q *OperationsQ) Page(page db2.PageQuery) OperationsQI {
	if q.Err != nil {
		return q
	}

	q.sql, q.Err = page.ApplyTo(q.sql, "ho.id")
	return q
}

// Update - updates existing operation
func (q *OperationsQ) Update(op Operation) error {
	sql := sq.Update("history_operations").SetMap(map[string]interface{}{
		"id": op.ID,
		"transaction_id": op.TransactionID,
		"application_order": op.ApplicationOrder,
		"type": op.Type,
		"details": op.DetailsString,
		"source_account": op.SourceAccount,
		"ledger_close_time": op.LedgerCloseTime,
		"identifier": op.Identifier,
		"state": op.State,
	}).Where("id = ?", op.ID)
	
	_, err := q.parent.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to update operation", logan.F{"operation_id": op.ID})
	}

	return nil
}

// Select loads the results of the query specified by `q` into `dest`.
func (q *OperationsQ) Select(dest interface{}) error {
	if q.Err != nil {
		return q.Err
	}

	q.Err = q.parent.Select(dest, q.sql)
	return q.Err
}
