package ingest

import (
	"time"

	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/core"
	"gitlab.com/tokend/horizon/db2/history"
	"gitlab.com/tokend/horizon/ingest/participants"
	"gitlab.com/tokend/horizon/resource/operations"
)

// Run starts an attempt to ingest the range of ledgers specified in this
// session.
func (is *Session) Run() {
	is.Err = is.Ingestion.Start()
	if is.Err != nil {
		return
	}

	defer is.Ingestion.Rollback()

	for is.Cursor.NextLedger() {
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

	// if this is ledger 1 or we are in paranoid mode,
	// create default system accounts and ingest their balances:
	// master, commission, storageFeeManage and operational
	if is.Paranoid || is.Cursor.LedgerSequence() == 1 {
		systemAccounts := []string{
			is.CoreInfo.MasterAccountID,
			is.CoreInfo.CommissionAccountID,
			is.CoreInfo.StorageFeeManageAccountID,
			is.CoreInfo.OperationalAccountID,
		}
		for _, address := range systemAccounts {
			_, is.Err = is.Ingestion.tryIngestAccount(address)
			if is.Err != nil {
				return
			}
			balances := []core.Balance{}
			is.Err = is.Ingestion.CoreQ.BalancesByAddress(&balances, address)
			if is.Err != nil {
				return
			}
			var created bool
			for _, balance := range balances {

				for _, asset := range is.CoreInfo.BaseAssets {
					if asset != balance.Asset {
						continue
					}
					created, is.Err = is.Ingestion.tryIngestBalance(
						balance.BalanceID, balance.Asset, balance.AccountID,
						balance.ExchangeID, balance.ExchangeName)
					if is.Err != nil {
						return
					}
					if !created {
						// balance already existed, expecting updates to exist too
						continue
					}

					// we can't always reliably determine balance amount,
					// zero is just safe default (hopefully)
					is.Err = is.Ingestion.tryIngestBalanceUpdate(
						balance.BalanceID, 0, is.Cursor.Ledger().CloseTime,
					)
					if is.Err != nil {
						return
					}
					break
				}
			}
		}
	}

	for is.Cursor.NextTx() {
		is.ingestTransaction()
	}

	is.priceHistory()

	is.Ingested++
	if is.Metrics != nil {
		is.Metrics.IngestLedgerTimer.Update(time.Since(start))
	}

	return
}

func (is *Session) priceHistory() {
	if is.Err != nil {
		return
	}

	priceHistory, err := is.Cursor.PriceHistoryProvider().ToPricePoints()
	if err != nil {
		is.Err = err
		return
	}

	if len(priceHistory) == 0 {
		return
	}

	q := is.Ingestion.priceHistory
	for _, price := range priceHistory {
		q = q.Values(price.BaseAsset, price.QuoteAsset, price.Timestamp, price.Price)
	}

	_, is.Err = is.Ingestion.DB.Exec(q)
}

func (is *Session) manageOffer(source xdr.AccountId, result xdr.ManageOfferResult) {
	if is.Err != nil {
		return
	}

	if result.Success == nil || len(result.Success.OffersClaimed) == 0 {
		return
	}

	q := is.Ingestion.trades
	for i := range result.Success.OffersClaimed {
		claimed := result.Success.OffersClaimed[i]
		q = q.Values(string(result.Success.BaseAsset),
			string(result.Success.QuoteAsset), int64(claimed.BaseAmount),
			int64(claimed.QuoteAmount), int64(claimed.CurrentPrice), time.Unix(is.Cursor.Ledger().CloseTime, 0).UTC())
	}

	_, is.Err = is.Ingestion.DB.Exec(q)
}

func (is *Session) processManageInvoice(op xdr.ManageInvoiceOp, result xdr.ManageInvoiceResult) {
	if is.Err != nil {
		return
	}
	if op.InvoiceId == 0 || op.Amount != 0 {
		return
	}
	is.Ingestion.UpdateInvoice(op.InvoiceId, history.CANCELED, nil)

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

func (is *Session) processForfeit(operation xdr.Operation, source xdr.AccountId, result xdr.OperationResultTr) {
	if is.Err != nil {
		return
	}
	forfeitOp := operation.Body.ForfeitOp
	state := true
	paymentId := result.ForfeitResult.MustSuccess().PaymentId

	details := operations.BasePayment{
		FromBalance:           forfeitOp.Balance.AsString(),
		ToBalance:             "",
		From:                  source.Address(),
		To:                    "",
		Amount:                amount.String(int64(forfeitOp.Amount)),
		SourcePaymentFee:      amount.String(0),
		DestinationPaymentFee: amount.String(0),
		SourceFixedFee:        amount.String(0),
		DestinationFixedFee:   amount.String(0),
		SourcePaysForDest:     false,
	}

	is.Err = is.Ingestion.InsertPaymentRequest(
		is.Cursor.Ledger(),
		uint64(paymentId),
		source.Address(),
		details,
		&state,
		forfeitOp.Type,
	)
	if is.Err != nil {
		return
	}
}

func (is *Session) processManageForfeitRequest(operation xdr.Operation, source xdr.AccountId, result xdr.OperationResultTr) {
	if is.Err != nil {
		return
	}
	manageRequestOp := operation.Body.ManageForfeitRequestOp
	manageRequestResult := result.MustManageForfeitRequestResult()

	details := operations.BasePayment{
		FromBalance:           manageRequestOp.Balance.AsString(),
		ToBalance:             "",
		From:                  source.Address(),
		To:                    "",
		Amount:                amount.String(int64(manageRequestOp.Amount)),
		SourcePaymentFee:      amount.String(0),
		DestinationPaymentFee: amount.String(0),
		SourceFixedFee:        amount.String(0),
		DestinationFixedFee:   amount.String(0),
		SourcePaysForDest:     false,
		UserDetails:           manageRequestOp.Details,
		Items:                 manageForfeitRequestToForfeitTimes(manageRequestResult),
	}
	is.Err = is.Ingestion.InsertPaymentRequest(
		is.Cursor.Ledger(),
		uint64(manageRequestResult.ForfeitRequestDetails.PaymentId),
		manageRequestResult.ForfeitRequestDetails.Exchange.Address(),
		details,
		nil,
		xdr.RequestTypeRequestTypeRedeem,
	)

	if is.Err != nil {
		return
	}
}

func (is *Session) processPayment(paymentOp xdr.PaymentOp, source xdr.AccountId, result xdr.PaymentResponse) {
	if is.Err != nil {
		return
	}

	invoiceReference := paymentOp.InvoiceReference
	if invoiceReference != nil {
		if invoiceReference.Accept && len(result.Exchanges) == 0 {
			is.Ingestion.UpdateInvoice(invoiceReference.InvoiceId, history.SUCCESS, nil)
		} else if !invoiceReference.Accept {
			is.Ingestion.UpdateInvoice(invoiceReference.InvoiceId, history.REJECTED, nil)
		}
	}

	if len(result.Exchanges) == 0 {
		return
	}
	for _, exchange := range result.Exchanges {
		details := operations.BasePayment{
			FromBalance:           paymentOp.SourceBalanceId.AsString(),
			ToBalance:             paymentOp.DestinationBalanceId.AsString(),
			From:                  source.Address(),
			To:                    result.Destination.Address(),
			Amount:                amount.String(int64(paymentOp.Amount)),
			SourcePaymentFee:      amount.String(int64(paymentOp.FeeData.SourceFee.PaymentFee)),
			DestinationPaymentFee: amount.String(int64(paymentOp.FeeData.DestinationFee.PaymentFee)),
			SourceFixedFee:        amount.String(int64(paymentOp.FeeData.SourceFee.PaymentFee)),
			DestinationFixedFee:   amount.String(int64(paymentOp.FeeData.DestinationFee.PaymentFee)),
			SourcePaysForDest:     paymentOp.FeeData.SourcePaysForDest,
		}
		is.Err = is.Ingestion.InsertPaymentRequest(
			is.Cursor.Ledger(),
			uint64(result.PaymentId),
			exchange.Address(),
			details,
			nil,
			xdr.RequestTypeRequestTypePayment,
		)

		if is.Err != nil {
			return
		}

	}
}

func (is *Session) processDemurrage(result xdr.OperationResultTr) {
	if is.Err != nil {
		return
	}
	opResult := result.MustDemurrageResult()
	if len(opResult.DemurrageInfo.PaymentRequests) == 0 {
		return
	}
	is.Err = is.Ingestion.InsertPaymentRequests(opResult.DemurrageInfo.PaymentRequests)
}

func (is *Session) processReviewEmissionRequest(operation xdr.Operation, result xdr.OperationResultTr) {
	if is.Err != nil {
		return
	}

	reviewRequestOp := operation.Body.ReviewCoinsEmissionRequestOp
	if reviewRequestOp.Approve {
		is.Err = is.Ingestion.ApproveEmissionRequestUpdate(
			is.Cursor.Ledger(),
			uint64(reviewRequestOp.Request.RequestId),
		)
	} else {
		is.Err = is.Ingestion.RejectEmissionRequestUpdate(
			is.Cursor.Ledger(),
			uint64(reviewRequestOp.Request.RequestId),
			string(reviewRequestOp.Reason),
		)
	}
	if is.Err != nil {
		return
	}
}

func (is *Session) updateIngestedPaymentRequest(operation xdr.Operation, source xdr.AccountId) {
	if is.Err != nil {
		return
	}
	reviewPaymentOp := operation.Body.MustReviewPaymentRequestOp()
	is.Err = is.Ingestion.UpdatePaymentRequest(
		is.Cursor.Ledger(),
		uint64(reviewPaymentOp.PaymentId),
		source.Address(),
		reviewPaymentOp.Accept,
	)
	if is.Err != nil {
		return
	}
}

func (is *Session) updateIngestedPayment(operation xdr.Operation, source xdr.AccountId, result xdr.OperationResultTr) {
	if is.Err != nil {
		return
	}
	reviewPaymentOp := operation.Body.MustReviewPaymentRequestOp()
	reviewPaymentResponse := result.MustReviewPaymentRequestResult().ReviewPaymentResponse

	if reviewPaymentResponse.RelatedInvoiceId != nil {
		if reviewPaymentOp.Accept {
			is.Ingestion.UpdateInvoice(*reviewPaymentResponse.RelatedInvoiceId,
				history.SUCCESS, nil)
		} else {
			is.Ingestion.UpdateInvoice(*reviewPaymentResponse.RelatedInvoiceId,
				history.FAILED, reviewPaymentOp.RejectReason)
		}
	}

	state := reviewPaymentResponse.State
	if state == xdr.PaymentStatePaymentPending {
		return
	}
	is.Err = is.Ingestion.UpdatePayment(
		reviewPaymentOp.PaymentId,
		state == xdr.PaymentStatePaymentProcessed,
		reviewPaymentOp.RejectReason,
	)
	if is.Err != nil {
		return
	}
}
