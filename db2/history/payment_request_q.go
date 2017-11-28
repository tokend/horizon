package history

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2"
	sq "github.com/lann/squirrel"
)

var selectPaymentRequest = sq.Select("hpr.*, ho.state").
	From("history_payment_requests hpr").
	LeftJoin("history_operations ho ON hpr.payment_id = ho.identifier")

type PaymentRequestsQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

type PaymentRequestsQI interface {
	Page(page db2.PageQuery) PaymentRequestsQI
	ForAccount(accountId string) PaymentRequestsQI
	ForBalance(balanceId string) PaymentRequestsQI
	ForState(state *bool) PaymentRequestsQI
	ForfeitsOnly() PaymentRequestsQI
	PaymentsOnly() PaymentRequestsQI
	Select(dest interface{}) error
}

func (q *Q) PaymentRequests() PaymentRequestsQI {
	return &PaymentRequestsQ{
		parent: q,
		sql:    selectPaymentRequest,
	}
}

func (q *Q) PaymentRequestByID(dest interface{}, id uint64) error {
	sql := selectPaymentRequest.Limit(1).Where("hpr.id = ?", id)
	return q.Get(dest, sql)
}
func (q *Q) PaymentRequestByPaymentID(dest interface{}, id uint64) error {
	sql := selectPaymentRequest.Limit(1).Where("hpr.payment_id = ?", id)
	return q.Get(dest, sql)
}

func (q *PaymentRequestsQ) ForBalance(balanceId string) PaymentRequestsQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("(hpr.details->>'from_balance' = ? OR hpr.details->>'to_balance' = ?)",
		balanceId, balanceId)

	return q
}

func (q *PaymentRequestsQ) ForAccount(accountId string) PaymentRequestsQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("(hpr.details->>'from' = ? OR hpr.details->>'to' = ?)",
		accountId, accountId)

	return q
}

func (q *PaymentRequestsQ) ForState(state *bool) PaymentRequestsQI {
	if state == nil {
		q.sql = q.sql.Where("hpr.accepted IS NULL")
	} else {
		q.sql = q.sql.Where("hpr.accepted = ?", state)
	}
	return q
}

func (q *PaymentRequestsQ) ForfeitsOnly() PaymentRequestsQI {
	q.sql = q.sql.Where("hpr.request_type != ?", int(xdr.RequestTypeRequestTypePayment))
	return q
}

func (q *PaymentRequestsQ) PaymentsOnly() PaymentRequestsQI {
	q.sql = q.sql.Where("hpr.request_type == ?", int(xdr.RequestTypeRequestTypePayment))
	return q
}

// Page specifies the paging constraints for the query being built by `q`.
func (q *PaymentRequestsQ) Page(page db2.PageQuery) PaymentRequestsQI {
	if q.Err != nil {
		return q
	}

	q.sql, q.Err = page.ApplyTo(q.sql, "hpr.id")
	return q
}

// Select loads the results of the query specified by `q` into `dest`.
func (q *PaymentRequestsQ) Select(dest interface{}) error {
	if q.Err != nil {
		return q.Err
	}

	q.Err = q.parent.Select(dest, q.sql)
	return q.Err
}
