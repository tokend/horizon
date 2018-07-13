package ingestion

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/swarmfund/horizon/db2/history"
)

func (ingest *Ingestion) UpdateOfferState(offerID, state uint64) error {
	sql := sq.Update("history_operations").
		Set("state", state).
		Where("type = ? AND details->>'offer_id' = ?", xdr.OperationTypeManageOffer, offerID)

	_, err := ingest.DB.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

func (ingest *Ingestion) UpdateOrderBookState(orderBookID, state uint64, ignoreCanceled bool) error {
	var sql sq.UpdateBuilder
	if ignoreCanceled {
		sql = sq.Update("history_operations").
			Set("state", state).
			Where("type = ? AND state <> ? AND details->>'order_book_id' = ?", xdr.OperationTypeManageOffer,
			history.OperationStateCanceled, orderBookID)
	} else {
		sql = sq.Update("history_operations").
			Set("state", state).
			Where("type = ? AND details->>'order_book_id' = ?", xdr.OperationTypeManageOffer, orderBookID)
	}

	_, err := ingest.DB.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}
