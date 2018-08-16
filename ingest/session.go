package ingest

import (
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/ingest/participants"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
)

// Run starts an attempt to ingest the range of ledgers specified in this
// session.
func (is *Session) Run() {
	err := is.Ingestion.Start()
	if err != nil {
		is.log.WithError(err).Error("failed to start ingestion")
		return
	}

	defer is.Ingestion.Rollback()

	for {
		isNextLegerLoaded, err := is.Cursor.NextLedger()
		if err != nil {
			is.log.WithError(err).Error("failed to load next ledger")
			return
		}

		if !isNextLegerLoaded {
			break
		}

		err = is.ingestLedger()
		if err != nil {
			is.log.WithError(err).Error("failed to ingest ledger")
			return
		}

		err = is.flush()
		if err != nil {
			is.log.WithError(err).Error("failed to flush")
			return
		}
	}

	err = is.Ingestion.Close()
	if err != nil {
		is.log.WithError(err).Error("failed to close ingestion")
		return
	}

	err = is.CoreConnector.SetCursor("HORIZON", is.Cursor.LastLedger)
	if err != nil {
		is.log.WithError(err).Error("failed to set cursor")
		return
	}
}

func (is *Session) flush() (err error) {
	err = is.Ingestion.Flush()
	if err != nil {
		return errors.Wrap(err, "failed to flush")
	}
	return nil
}

// ingestLedger ingests the current ledger
func (is *Session) ingestLedger() (err error) {

	start := time.Now()
	err = is.Ingestion.Ledger(
		is.Cursor.LedgerID(),
		is.Cursor.Ledger(),
		is.Cursor.SuccessfulTransactionCount(),
		is.Cursor.SuccessfulLedgerOperationCount(),
	)

	if err != nil {
		return errors.Wrap(err, "failed to add ledger to current ingestion")
	}

	// ingest accounts from genesis block
	if is.Cursor.LedgerSequence() == 1 {
		systemAccounts := []string{
			is.CoreInfo.MasterAccountID,
			is.CoreInfo.CommissionAccountID,
			is.CoreInfo.OperationalAccountID,
		}
		for _, address := range systemAccounts {
			_, err = is.Ingestion.TryIngestAccount(address)
			if err != nil {
				return errors.Wrap(err, "failed to ingest account")
			}
		}
	}

	for is.Cursor.NextTx() {
		err = is.ingestTransaction()
		if err != nil {
			return errors.Wrap(err, "failed to ingest transaction")
		}
	}

	is.Ingested++
	if is.Metrics != nil {
		is.Metrics.IngestLedgerTimer.Update(time.Since(start))
	}

	return nil
}

func (is *Session) storeTrades(orderBookID uint64, result xdr.ManageOfferSuccessResult) error {
	return is.Ingestion.StoreTrades(orderBookID, result, is.Cursor.Ledger().CloseTime)
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

func (is *Session) processManageOfferLedgerChanges(offerID uint64) error {
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
				return errors.Wrap(err, "failed to update offer details", logan.F{
					"offer_id": uint64(change.Updated.Data.Offer.OfferId),
				})
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
				return errors.Wrap(err, "failed to update offer state", logan.F{
					"offer_id": uint64(change.Removed.Offer.OfferId),
				})
			}
		}
	}
	return nil
}

func (is *Session) permanentReject(op xdr.ReviewRequestOp) error {
	err := is.Ingestion.HistoryQ().ReviewableRequests().PermanentReject(uint64(op.RequestId), string(op.Reason))
	if err != nil {
		return errors.Wrap(err, "failed to permanently reject request")
	}

	return nil
}

func (is *Session) handleCheckSaleState(result xdr.CheckSaleStateSuccess) error {
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
		return errors.Wrap(err, "failed to set state", logan.F{"sale_id": uint64(result.SaleId)})
	}

	err = is.Ingestion.UpdateOrderBookState(uint64(result.SaleId), offerState, true)
	if err != nil {
		return errors.Wrap(err, "failed to set offers states", logan.F{"sale_id": uint64(result.SaleId)})
	}
	return nil
}

func (is *Session) handleManageSale(op *xdr.ManageSaleOp) error {
	if op.Data.Action != xdr.ManageSaleActionCancel {
		return nil
	}

	err := is.Ingestion.HistoryQ().Sales().SetState(uint64(op.SaleId), history.SaleStateCanceled)
	if err != nil {
		return errors.Wrap(err, "failed to set state", logan.F{"sale_id": uint64(op.SaleId)})
	}

	err = is.Ingestion.UpdateOrderBookState(uint64(op.SaleId), uint64(history.OperationStateCanceled), false)
	if err != nil {
		return errors.Wrap(err, "failed to set offers states", logan.F{"sale_id": uint64(op.SaleId)})
	}
	return nil
}

func (is *Session) processManageAsset(op *xdr.ManageAssetOp) error {
	if op.Request.Action != xdr.ManageAssetActionCancelAssetRequest {
		return nil
	}

	err := is.Ingestion.HistoryQ().ReviewableRequests().Cancel(uint64(op.RequestId))
	if err != nil {
		return errors.Wrap(err, "failed to cancel reviewable request", map[string]interface{}{
			"request_id": uint64(op.RequestId),
		})
	}
	return nil
}

func (is *Session) ingestOperationParticipants() error {
	// Find the participants
	var opParticipants []participants.Participant
	opParticipants, err := participants.ForOperation(
		is.Ingestion.DB,
		&is.Cursor.Transaction().Envelope.Tx,
		is.Cursor.Operation(),
		*is.Cursor.OperationResult(),
		is.Cursor.OperationChanges(),
		is.Cursor.Ledger(),
	)

	if err != nil {
		return errors.Wrap(err, "failed to load operation participants")
	}

	if len(opParticipants) == 0 {
		return nil
	}

	err = is.Ingestion.OperationParticipants(is.Cursor.OperationID(), opParticipants)
	if err != nil {
		return errors.Wrap(err, "failed to insert participants info into database")
	}
	return nil
}

func (is *Session) ingestTransaction() error {
	// skip ingesting failed transactions
	if !is.Cursor.Transaction().IsSuccessful() {
		return nil
	}

	err := is.Ingestion.Transaction(
		is.Cursor.Ledger(),
		is.Cursor.TransactionID(),
		is.Cursor.Transaction(),
		is.Cursor.TransactionFee(),
	)
	if err != nil {
		return errors.Wrap(err, "failed to ingest transaction", map[string]interface{}{
			"tx_id": is.Cursor.TransactionID(),
		})
	}

	for is.Cursor.NextOp() {
		err = is.operation()
		if err != nil {
			return errors.Wrap(err, "failed to ingest operation")
		}
	}

	err = is.ingestTransactionParticipants()
	if err != nil {
		return errors.Wrap(err, "failed to ingest transactions participants")
	}
	return nil
}

func (is *Session) ingestTransactionParticipants() (err error) {
	// Find the participants

	var p []xdr.AccountId
	p, err = participants.ForTransaction(
		is.Ingestion.DB,
		&is.Cursor.Transaction().Envelope.Tx,
		*is.Cursor.Transaction().Result.Result.Result.Results,
		&is.Cursor.Transaction().ResultMeta,
		&is.Cursor.TransactionFee().Changes,
		is.Cursor.Ledger(),
	)
	if err != nil {
		return errors.Wrap(err, "failed to get participants ids for transaction")
	}

	err = is.Ingestion.TransactionParticipants(is.Cursor.TransactionID(), p)
	if err != nil {
		return errors.Wrap(err, "failed to load transaction participants")
	}
	return nil
}

func (is *Session) processManageInvoiceRequest(op xdr.ManageInvoiceRequestOp,
	result xdr.ManageInvoiceRequestResult,
) error {
	if result.Code != xdr.ManageInvoiceRequestResultCodeSuccess {
		return nil
	}
	if op.Details.Action == xdr.ManageInvoiceRequestActionCreate {
		return nil
	}

	err := is.Ingestion.HistoryQ().ReviewableRequests().Cancel(uint64(*op.Details.RequestId))
	if err != nil {
		return errors.Wrap(err, "failed to update invoice request state to cancel", logan.F{
			"request_id": uint64(*op.Details.RequestId),
		})
	}

	return nil
}

func (is *Session) processManageContractRequest(
	op xdr.ManageContractRequestOp,
	result xdr.ManageContractRequestResult,
) error {
	if result.Code != xdr.ManageContractRequestResultCodeSuccess {
		return nil
	}
	if op.Details.Action == xdr.ManageContractRequestActionCreate {
		return nil
	}

	err := is.Ingestion.HistoryQ().ReviewableRequests().Cancel(uint64(*op.Details.RequestId))
	if err != nil {
		return errors.Wrap(err, "failed to update contract request state to cancel", logan.F{
			"request_id": uint64(*op.Details.RequestId),
		})
	}

	return nil
}
