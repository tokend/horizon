package ingestion

import (
	"encoding/json"

	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/tokend/go/xdr"
)

func (ingest *Ingestion) UpdateOfferDetails(newOfferDetails map[string]interface{}, stateToSet uint64) error {
	bytes, err := json.Marshal(newOfferDetails)
	if err != nil {
		return err
	}

	sql := sq.Update("history_operations").
		SetMap(map[string]interface{}{
			"details": bytes,
			"state":   stateToSet,
		}).Where("type = ? AND details->>'offer_id' = ?", xdr.OperationTypeManageOffer, newOfferDetails["offer_id"])

	_, err = ingest.DB.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to update history_operations")
	}

	return nil
}

func (ingest *Ingestion) UpdateOfferState(offerID, state uint64) error {
	sql := sq.Update("history_operations").
		Set("state", state).
		Where("type = ? AND details->>'offer_id' = ?", xdr.OperationTypeManageOffer, offerID)

	_, err := ingest.DB.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to update history_operations")
	}

	return nil
}

func (ingest *Ingestion) UpdateOrderBookState(orderBookID, state uint64, ignoreCanceled bool) error {
	sql := sq.Update("history_operations").
		Set("state", state).
		Where("type = ? AND details->>'order_book_id' = ?", xdr.OperationTypeManageOffer, orderBookID)
	if ignoreCanceled {
		sql = sql.Where("state <> ?", history.OperationStateCanceled)
	}

	_, err := ingest.DB.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to update history_operations")
	}

	return nil
}
