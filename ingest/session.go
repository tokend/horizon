package ingest

import (
	"time"

	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/ingest/participants"
	"gitlab.com/swarmfund/horizon/resource/operations"
	"gitlab.com/distributed_lab/logan/v3/errors"
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

	// if this is ledger 1 or we are in paranoid mode,
	// create default system accounts and ingest their balances:
	// master, commission, storageFeeManage and operational
	if is.Paranoid || is.Cursor.LedgerSequence() == 1 {
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
					created, is.Err = is.Ingestion.TryIngestBalance(
						balance.BalanceID, balance.Asset, balance.AccountID)
					if is.Err != nil {
						return
					}
					if !created {
						// balance already existed, expecting updates to exist too
						continue
					}

					// we can't always reliably determine balance amount,
					// zero is just safe default (hopefully)
					is.Err = is.Ingestion.TryIngestBalanceUpdate(
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

	err = is.Ingestion.StorePricePoints(priceHistory)
	if err != nil {
		is.Err = errors.Wrap(err, "failed to store price points")
	}
}

func (is *Session) manageOffer(source xdr.AccountId, result xdr.ManageOfferResult) {
	if is.Err != nil {
		return
	}

	is.Err = is.Ingestion.StoreTrades(source, result, is.Cursor.Ledger().CloseTime)
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

func (is *Session) processReviewRequest(op xdr.ReviewRequestOp) {
	if is.Err != nil {
		return
	}

	var err error
	switch op.Action {
	case xdr.ReviewRequestOpActionApprove:
		err = is.Cursor.HistoryQ().ReviewableRequests().Approve(uint64(op.RequestId))
	case xdr.ReviewRequestOpActionPermanentReject:
		err = is.Cursor.HistoryQ().ReviewableRequests().PermanentReject(uint64(op.RequestId), string(op.Reason))
	case xdr.ReviewRequestOpActionReject:
		return
	default:
		err = errors.From(errors.New("Unexpected review request action"), map[string]interface{}{
			"action_type": op.Action,
		})
	}

	if err != nil {
		is.Err = errors.Wrap(err, "failed to process review request", map[string]interface{}{
			"request_id": uint64(op.RequestId),
		})
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
	}
	is.Err = is.Ingestion.InsertPaymentRequest(
		is.Cursor.Ledger(),
		uint64(manageRequestResult.Success.PaymentId),
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
		if invoiceReference.Accept {
			is.Ingestion.UpdateInvoice(invoiceReference.InvoiceId, history.SUCCESS, nil)
		} else if !invoiceReference.Accept {
			is.Ingestion.UpdateInvoice(invoiceReference.InvoiceId, history.REJECTED, nil)
		}
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
	if state == xdr.PaymentStatePending {
		return
	}
	is.Err = is.Ingestion.UpdatePayment(
		reviewPaymentOp.PaymentId,
		state == xdr.PaymentStateProcessed,
		reviewPaymentOp.RejectReason,
	)
	if is.Err != nil {
		return
	}
}
