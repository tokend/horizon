package ingestion

import (
			"gitlab.com/tokend/go/xdr"
		"gitlab.com/swarmfund/horizon/db2/history"
	sq "github.com/lann/squirrel"
	"fmt"
)

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
			return err
		}
	}
	sql := sq.Update("history_operations").SetMap(sq.Eq{
		"state": state,
	}).Where("identifier = ?", paymentID)

	_, err := ingest.DB.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

func (ingest *Ingestion) ingestRejectReason(rejectReason string, operationID uint64) error {
	addRejectReasonQuery := fmt.Sprintf("UPDATE history_operations SET details = jsonb_set(details, '{reject_reason}', '\"%s\"') "+
		"WHERE identifier = %v", rejectReason, operationID)
	_, err := ingest.DB.ExecRaw(addRejectReasonQuery)
	return err
}
