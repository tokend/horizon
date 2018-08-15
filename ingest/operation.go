package ingest

import (
	"fmt"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/tokend/go/xdr"
)

func getStateIdentifier(opType xdr.OperationType, op *xdr.Operation, operationResult *xdr.OperationResultTr) (history.OperationState, uint64) {
	state := history.OperationStateSuccess
	operationIdentifier := uint64(0)
	switch opType {
	case xdr.OperationTypePayment, xdr.OperationTypeDirectDebit:
		var paymentResponse xdr.PaymentResponse
		if opType == xdr.OperationTypePayment {
			paymentResponse = *operationResult.MustPaymentResult().PaymentResponse
		} else {
			paymentResponse = operationResult.MustDirectDebitResult().MustSuccess().PaymentResponse
		}

		operationIdentifier = uint64(paymentResponse.PaymentId)
		return state, operationIdentifier
	case xdr.OperationTypeCreateWithdrawalRequest:
		state = history.OperationStatePending
		manageRequestResult := operationResult.MustCreateWithdrawalRequestResult()
		operationIdentifier = uint64(manageRequestResult.Success.RequestId)
		return state, operationIdentifier
	case xdr.OperationTypeManageInvoice:
		manageInvoiceOp := op.Body.MustManageInvoiceOp()
		if manageInvoiceOp.InvoiceId != 0 {
			return state, operationIdentifier
		}

		state = history.OperationStatePending
		manageInvoiceResult := operationResult.MustManageInvoiceResult()
		operationIdentifier = uint64(manageInvoiceResult.Success.InvoiceId)
		return state, operationIdentifier
	case xdr.OperationTypeCreateIssuanceRequest:
		createIssuanceRequestResult := operationResult.MustCreateIssuanceRequestResult()
		state = history.OperationStatePending
		if createIssuanceRequestResult.Success.Fulfilled {
			state = history.OperationStateSuccess
		}
		return state, uint64(createIssuanceRequestResult.Success.RequestId)
	case xdr.OperationTypeManageOffer:
		manageOfferOp := op.Body.MustManageOfferOp()
		manageOfferResult := operationResult.MustManageOfferResult().MustSuccess()

		switch manageOfferResult.Offer.Effect {
		case xdr.ManageOfferEffectCreated:
			if len(manageOfferResult.OffersClaimed) == 0 {
				return history.OperationStatePending, 0
			}
			return history.OperationStatePartiallyMatched, 0
		case xdr.ManageOfferEffectDeleted:
			if manageOfferOp.Amount != 0 {
				return history.OperationStateFullyMatched, 0
			}
			return history.OperationStateCanceled, 0
		default:
			panic(fmt.Errorf("unknown manage offer op effect: %s", manageOfferResult.Offer.Effect.ShortString()))
		}
	default:
		return state, operationIdentifier
	}
}

func (is *Session) operation() error {

	err := is.operationChanges(is.Cursor.OperationChanges())
	if err != nil {
		return errors.Wrap(err, "failed to process operation changes")
	}

	state, operationIdentifier := getStateIdentifier(is.Cursor.OperationType(), is.Cursor.Operation(), is.Cursor.OperationResult())
	err = is.Ingestion.Operation(
		is.Cursor.OperationID(),
		is.Cursor.TransactionID(),
		is.Cursor.OperationOrder(),
		is.Cursor.OperationSourceAccount(),
		is.Cursor.OperationType(),
		is.operationDetails(),
		is.Cursor.Ledger().CloseTime,
		operationIdentifier,
		state,
	)
	if err != nil {
		return errors.Wrap(err, "failed to ingest operation")
	}

	err = is.ingestOperationParticipants()
	if err != nil {
		return errors.Wrap(err, "failed to ingest operation participants")
	}
	switch is.Cursor.OperationType() {
	case xdr.OperationTypePayment:
		err = is.processPayment(is.Cursor.Operation().Body.MustPaymentOp(), is.Cursor.OperationSourceAccount(),
			*is.Cursor.OperationResult().MustPaymentResult().PaymentResponse)
		if err != nil {
			return errors.Wrap(err, "failed to process payment")
		}

	case xdr.OperationTypeReviewPaymentRequest:
		err = is.updateIngestedPaymentRequest(*is.Cursor.Operation(), is.Cursor.OperationSourceAccount())
		if err != nil {
			return errors.Wrap(err, "failed to update ingested payment request")
		}

		err = is.updateIngestedPayment(*is.Cursor.Operation(), is.Cursor.OperationSourceAccount(), *is.Cursor.OperationResult())
		if err != nil {
			return errors.Wrap(err, "failed to update ingested payment")
		}

	case xdr.OperationTypeDirectDebit:
		opDirectDebit := is.Cursor.Operation().Body.MustDirectDebitOp()
		err = is.processPayment(opDirectDebit.PaymentOp,
			opDirectDebit.From,
			is.Cursor.OperationResult().MustDirectDebitResult().MustSuccess().PaymentResponse)
		if err != nil {
			return errors.Wrap(err, "failed to process payment")
		}

	case xdr.OperationTypeManageOffer:
		op := is.Cursor.Operation().Body.MustManageOfferOp()
		opResult := is.Cursor.OperationResult().MustManageOfferResult().MustSuccess()
		err = is.storeTrades(uint64(is.Cursor.Operation().Body.MustManageOfferOp().OrderBookId), opResult)
		if err != nil {
			return errors.Wrap(err, "failed to store trades")
		}

		offerIsCancelled := op.OfferId != 0 && op.Amount == 0
		if offerIsCancelled {
			err = is.updateOfferState(uint64(op.OfferId), uint64(history.OperationStateCanceled))
			if err != nil {
				return errors.Wrap(err, "failed to update offer state")
			}
			return nil
		}

		err = is.processManageOfferLedgerChanges(uint64(is.Cursor.Operation().Body.MustManageOfferOp().OfferId))
		if err != nil {
			return errors.Wrap(err, "failed to process manage offer ledger changes")
		}

	case xdr.OperationTypeManageInvoice:
		err = is.processManageInvoice(is.Cursor.Operation().Body.MustManageInvoiceOp(),
			is.Cursor.OperationResult().MustManageInvoiceResult())
		if err != nil {
			return errors.Wrap(err, "failed to process manage invoice operation")
		}

	case xdr.OperationTypeReviewRequest:
		err = is.processReviewRequest(
			is.Cursor.Operation().Body.MustReviewRequestOp(),
			is.Cursor.OperationResult().ReviewRequestResult.MustSuccess(),
			is.Cursor.OperationChanges(),
		)
		if err != nil {
			return errors.Wrap(err, "failed to process review request")
		}
	case xdr.OperationTypeManageAsset:
		err = is.processManageAsset(is.Cursor.Operation().Body.ManageAssetOp)
		if err != nil {
			return errors.Wrap(err, "failed to process manage asset operation")
		}

	case xdr.OperationTypeCheckSaleState:
		success := *is.Cursor.OperationResult().MustCheckSaleStateResult().Success
		err = is.handleCheckSaleState(success)
		if err != nil {
			return errors.Wrap(err, "failed to handle check sale state")
		}

		if success.Effect.Effect == xdr.CheckSaleStateEffectClosed {
			closed := success.Effect.SaleClosed
			for i := range closed.Results {
				err = is.storeTrades(uint64(success.SaleId), closed.Results[i].SaleDetails)
				if err != nil {
					errors.Wrap(err, "failed to insert sale into db", logan.F{
						"sale":         success.SaleId,
						"sale results": closed.Results[i],
					})
				}
			}
		}
	case xdr.OperationTypeManageSale:
		opManageSale := is.Cursor.Operation().Body.MustManageSaleOp()
		err = is.handleManageSale(&opManageSale)
		if err != nil {
			return errors.Wrap(err, "failed to handle manage sale")
		}
	}
	return nil
}
