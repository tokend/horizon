package ingest

import (
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/ingest/participants"
	"gitlab.com/tokend/go/amount"
)

// Run starts an attempt to ingest the range of ledgers specified in this
// session.
func (is *Session) Run() {
	is.Err = is.Ingestion.Start()
	if is.Err != nil {
		return
	}

	defer is.Ingestion.Rollback()

	for {
		isNextLegerLoaded, err := is.Cursor.NextLedger()
		if err != nil {
			is.Err = errors.Wrap(err, "failed to load next ledger")
			return
		}

		if !isNextLegerLoaded {
			break
		}

		if is.Err != nil {
			return
		}

		is.ingestLedger()
		is.flush()
	}

	if is.Err != nil {
		is.Ingestion.Rollback()
		return
	}

	is.Err = is.Ingestion.Close()
	if is.Err != nil {
		return
	}

	is.Err = is.CoreConnector.SetCursor("HORIZON", is.Cursor.LastLedger)
}

func (is *Session) flush() {
	if is.Err != nil {
		return
	}
	is.Err = is.Ingestion.Flush()
}

// ingestLedger ingests the current ledger
func (is *Session) ingestLedger() {
	if is.Err != nil {
		return
	}

	start := time.Now()
	is.Err = is.Ingestion.Ledger(
		is.Cursor.LedgerID(),
		is.Cursor.Ledger(),
		is.Cursor.SuccessfulTransactionCount(),
		is.Cursor.SuccessfulLedgerOperationCount(),
	)

	if is.Err != nil {
		return
	}

	// ingest accounts from genesis block
	if is.Cursor.LedgerSequence() == 1 {
		systemAccounts := []string{
			is.CoreInfo.MasterAccountID,
			is.CoreInfo.CommissionAccountID,
			is.CoreInfo.OperationalAccountID,
		}
		for _, address := range systemAccounts {
			_, is.Err = is.Ingestion.TryIngestAccount(address)
			if is.Err != nil {
				return
			}
		}
	}

	for is.Cursor.NextTx() {
		is.ingestTransaction()
	}

	is.Ingested++
	if is.Metrics != nil {
		is.Metrics.IngestLedgerTimer.Update(time.Since(start))
	}

	return
}

func (is *Session) storeTrades(orderBookID uint64, result xdr.ManageOfferSuccessResult) {
	if is.Err != nil {
		return
	}

	is.Err = is.Ingestion.StoreTrades(orderBookID, result, is.Cursor.Ledger().CloseTime)
}

func parseOfferEntryToDetails(offerEntry xdr.OfferEntry) map[string]interface{} {
	result := map[string]interface{}{}
	result["fee"] = amount.String(int64(offerEntry.Fee))
	result["price"] = amount.String(int64(offerEntry.Price))
	result["amount"] = amount.String(int64(offerEntry.BaseAmount))
	result["is_buy"] = offerEntry.IsBuy
	result["offer_id"] = uint64(offerEntry.OfferId)
	result["is_deleted"] = false
	result["order_book_id"] = uint64(offerEntry.OrderBookId)
	return result
}

func (is *Session) processManageOfferLedgerChanges(offerID uint64) {
	ledgerChanges := is.Cursor.OperationChanges()
	for _, change := range ledgerChanges {
		switch change.Type {
		case xdr.LedgerEntryChangeTypeUpdated:
			if change.Updated.Data.Type != xdr.LedgerEntryTypeOfferEntry {
				continue
			}
			if uint64(change.Updated.Data.Offer.OfferId) == offerID {
				continue
			}
			offerDetails := parseOfferEntryToDetails(*change.Updated.Data.Offer)
			err := is.Ingestion.UpdateOfferDetails(offerDetails, uint64(history.OperationStatePartiallyMatched))
			if err != nil {
				is.Err = errors.Wrap(err, "failed to update offer details", logan.F{
					"offer_id": uint64(change.Updated.Data.Offer.OfferId),
				})
				return
			}
		case xdr.LedgerEntryChangeTypeRemoved:
			if change.Removed.Type != xdr.LedgerEntryTypeOfferEntry {
				continue
			}
			if uint64(change.Removed.Offer.OfferId) == offerID {
				continue
			}
			err := is.Ingestion.UpdateOfferState(uint64(change.Removed.Offer.OfferId),
				uint64(history.OperationStateExternallyFullyMatched))
			if err != nil {
				is.Err = errors.Wrap(err, "failed to update offer state", logan.F{
					"offer_id": uint64(change.Removed.Offer.OfferId),
				})
				return
			}
		}
	}
}

func (is *Session)  permanentReject(op xdr.ReviewRequestOp) error {
	err := is.Ingestion.HistoryQ().ReviewableRequests().PermanentReject(uint64(op.RequestId), string(op.Reason))
	if err != nil {
		return errors.Wrap(err, "failed to permanently reject request")
	}

	err = is.Ingestion.UpdatePayment(op.RequestId, false, nil)
	if err != nil {
		return errors.Wrap(err, "failed to permanently reject operation")
	}

	return nil
}

func (is *Session) handleCheckSaleState(result xdr.CheckSaleStateSuccess) {
	if is.Err != nil {
		return
	}

	state := history.SaleStateClosed
	if result.Effect.Effect == xdr.CheckSaleStateEffectCanceled {
		state = history.SaleStateCanceled
	}

	var offerState uint64
	switch state {
	case history.SaleStateCanceled:
		offerState = uint64(history.OperationStateCanceled)
	case history.SaleStateClosed:
		offerState = uint64(history.OperationStateFullyMatched)
	}

	err := is.Ingestion.HistoryQ().Sales().SetState(uint64(result.SaleId), state)
	if err != nil {
		is.Err = errors.Wrap(err, "failed to set state", logan.F{"sale_id": uint64(result.SaleId)})
		return
	}

	err = is.Ingestion.UpdateOrderBookState(uint64(result.SaleId), offerState, true)
	if err != nil {
		is.Err = errors.Wrap(err, "failed to set offers states", logan.F{"sale_id": uint64(result.SaleId)})
		return
	}
}

func (is *Session) handleManageSale(op *xdr.ManageSaleOp) {
	if is.Err != nil {
		return
	}

	if op.Data.Action != xdr.ManageSaleActionCancel {
		return
	}

	err := is.Ingestion.HistoryQ().Sales().SetState(uint64(op.SaleId), history.SaleStateCanceled)
	if err != nil {
		is.Err = errors.Wrap(err, "failed to set state", logan.F{"sale_id": uint64(op.SaleId)})
		return
	}

	err = is.Ingestion.UpdateOrderBookState(uint64(op.SaleId), uint64(history.OperationStateCanceled), false)
	if err != nil {
		is.Err = errors.Wrap(err, "failed to set offers states", logan.F{"sale_id": uint64(op.SaleId)})
		return
	}
}

func (is *Session) processManageAsset(op *xdr.ManageAssetOp) {
	if is.Err != nil {
		return
	}

	if op.Request.Action != xdr.ManageAssetActionCancelAssetRequest {
		return
	}

	err := is.Ingestion.HistoryQ().ReviewableRequests().Cancel(uint64(op.RequestId))
	if err != nil {
		is.Err = errors.Wrap(err, "failed to cancel reviewable request", map[string]interface{}{
			"request_id": uint64(op.RequestId),
		})
		return
	}
}

func (is *Session) ingestOperationParticipants() {
	if is.Err != nil {
		return
	}

	// Find the participants
	var opParticipants []participants.Participant
	opParticipants, is.Err = participants.ForOperation(
		is.Ingestion.DB,
		&is.Cursor.Transaction().Envelope.Tx,
		is.Cursor.Operation(),
		*is.Cursor.OperationResult(),
		is.Cursor.OperationChanges(),
		is.Cursor.Ledger(),
	)

	if is.Err != nil {
		return
	}

	if len(opParticipants) == 0 {
		return
	}

	is.Err = is.Ingestion.OperationParticipants(is.Cursor.OperationID(), opParticipants)
	if is.Err != nil {
		return
	}
}

func (is *Session) ingestTransaction() {
	if is.Err != nil {
		return
	}

	// skip ingesting failed transactions
	if !is.Cursor.Transaction().IsSuccessful() {
		return
	}

	is.Ingestion.Transaction(
		is.Cursor.Ledger(),
		is.Cursor.TransactionID(),
		is.Cursor.Transaction(),
		is.Cursor.TransactionFee(),
	)

	for is.Cursor.NextOp() {
		is.operation()
	}

	is.ingestTransactionParticipants()
}

func (is *Session) ingestTransactionParticipants() {
	if is.Err != nil {
		return
	}

	// Find the participants

	var p []xdr.AccountId
	p, is.Err = participants.ForTransaction(
		is.Ingestion.DB,
		&is.Cursor.Transaction().Envelope.Tx,
		*is.Cursor.Transaction().Result.Result.Result.Results,
		&is.Cursor.Transaction().ResultMeta,
		&is.Cursor.TransactionFee().Changes,
		is.Cursor.Ledger(),
	)
	if is.Err != nil {
		return
	}

	is.Err = is.Ingestion.TransactionParticipants(is.Cursor.TransactionID(), p)
	if is.Err != nil {
		return
	}

}

func (is *Session) processPayment(paymentOp xdr.PaymentOp, source xdr.AccountId, result xdr.PaymentResponse) {
	if is.Err != nil {
		return
	}
}

func (is *Session) processManageInvoiceRequest(op xdr.ManageInvoiceRequestOp, result xdr.ManageInvoiceRequestResult) {
	if is.Err != nil {
		return
	}
	if result.Code != xdr.ManageInvoiceRequestResultCodeSuccess {
		return
	}
	if op.Details.Action == xdr.ManageInvoiceRequestActionCreate {
		return
	}

	err := is.Ingestion.HistoryQ().ReviewableRequests().Cancel(uint64(*op.Details.RequestId))
	if err != nil {
		is.Err = errors.Wrap(err, "failed to update invoice request state to cancel", logan.F{
			"request_id": uint64(*op.Details.RequestId),
		})
	}
}