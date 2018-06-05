package ingestion

import (
	"fmt"

	"gitlab.com/tokend/go/xdr"
	sq "github.com/lann/squirrel"
	"gitlab.com/swarmfund/horizon/db2/history"
)

func (ingest *Ingestion) UpdateInvoice(
	invoiceID xdr.Uint64,
	state history.OperationState,
	rejectReason *xdr.Longstring,
) error {
	if rejectReason != nil {
		err := ingest.ingestRejectReason(string(*rejectReason), uint64(invoiceID))
		if err != nil {
			return err
		}
	}
	sql := sq.Update("history_operations").SetMap(sq.Eq{
		"state": state,
	}).Where("identifier = ?", invoiceID)

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
