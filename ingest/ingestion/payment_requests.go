package ingestion

import (
	"encoding/json"
	"time"

	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/tokend/go/xdr"
)

func (ingest *Ingestion) InsertPaymentRequest(
	ledger *core.LedgerHeader,
	paymentID uint64,
	details interface{},
	accepted *bool,
	requestType xdr.RequestType,
) error {
	ledgerCloseTime := time.Unix(ledger.CloseTime, 0)
	djson, err := json.Marshal(details)
	if err != nil {
		return err
	}

	sql := ingest.payment_requests.Values(
		paymentID,
		accepted,
		djson,
		ledgerCloseTime,
		ledgerCloseTime,
		int(requestType),
	)

	_, err = ingest.DB.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to execute sql query")
	}

	return nil
}

func (ingest *Ingestion) UpdatePaymentRequest(
	ledger *core.LedgerHeader,
	paymentID uint64,
	accept bool,
) error {
	sql := sq.Update("history_payment_requests").SetMap(sq.Eq{
		"accepted":   accept,
		"updated_at": time.Unix(ledger.CloseTime, 0),
	}).Where("payment_id = ?", paymentID)

	_, err := ingest.DB.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to update history_payment_request")
	}

	return nil
}

func (ingest *Ingestion) UpdatePayment(
	paymentID xdr.Uint64,
	accept bool,
	rejectReason *xdr.Longstring,
) error {
	state := history.OperationStateSuccess
	if !accept {
		state = history.OperationStateRejected
	}

	if rejectReason != nil {
		err := ingest.ingestRejectReason(string(*rejectReason), uint64(paymentID))
		if err != nil {
			return errors.Wrap(err, "failed to ingest reject reason")
		}
	}
	sql := sq.Update("history_operations").SetMap(sq.Eq{
		"state": state,
	}).Where("identifier = ?", paymentID)

	_, err := ingest.DB.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to update history_operations")
	}

	return nil
}
