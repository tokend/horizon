package ingest

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history"
)

func getStateIdentifier(opType xdr.OperationType, op *xdr.Operation, operationResult *xdr.OperationResultTr) (int, uint64) {
	state := history.SUCCESS
	operationIdentifier := uint64(0)
	switch opType {
	case xdr.OperationTypePayment, xdr.OperationTypeDirectDebit:
		var paymentResponse xdr.PaymentResponse
		if opType == xdr.OperationTypePayment {
			paymentResponse = *operationResult.MustPaymentResult().PaymentResponse
		} else {
			paymentResponse = operationResult.MustDirectDebitResult().MustSuccess().PaymentResponse
		}

		if len(paymentResponse.Exchanges) > 0 {
			state = history.PENDING
		}
		operationIdentifier = uint64(paymentResponse.PaymentId)
		return state, operationIdentifier
	case xdr.OperationTypeManageForfeitRequest:
		state = history.PENDING
		manageRequestResult := operationResult.MustManageForfeitRequestResult()
		operationIdentifier = uint64(manageRequestResult.ForfeitRequestDetails.PaymentId)
		return state, operationIdentifier
	case xdr.OperationTypeManageInvoice:
		manageInvoiceOp := op.Body.MustManageInvoiceOp()
		if manageInvoiceOp.InvoiceId != 0 {
			return state, operationIdentifier
		}

		state = history.PENDING
		manageInvoiceResult := operationResult.MustManageInvoiceResult()
		operationIdentifier = uint64(manageInvoiceResult.Success.InvoiceId)
		return state, operationIdentifier
	default:
		return state, operationIdentifier
	}
}

func (is *Session) operation() {
	if is.Err != nil {
		return
	}

	err := is.operationChanges(is.Cursor.OperationChanges())
	if err != nil {
		is.log.WithError(err).Error("Failed to process operation changes")
		is.Err = err
		return
	}

	state, operationIdentifier := getStateIdentifier(is.Cursor.OperationType(), is.Cursor.Operation(), is.Cursor.OperationResult())
	is.Err = is.Ingestion.Operation(
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
	if is.Err != nil {
		return
	}

	is.ingestOperationParticipants()
	switch is.Cursor.OperationType() {
	case xdr.OperationTypeReviewCoinsEmissionRequest:
		is.processReviewEmissionRequest(*is.Cursor.Operation(), *is.Cursor.OperationResult())
	case xdr.OperationTypeForfeit:
		is.processForfeit(*is.Cursor.Operation(), is.Cursor.OperationSourceAccount(), *is.Cursor.OperationResult())
	case xdr.OperationTypeManageForfeitRequest:
		is.processManageForfeitRequest(*is.Cursor.Operation(), is.Cursor.OperationSourceAccount(), *is.Cursor.OperationResult())
	case xdr.OperationTypePayment:
		is.processPayment(is.Cursor.Operation().Body.MustPaymentOp(), is.Cursor.OperationSourceAccount(),
			*is.Cursor.OperationResult().MustPaymentResult().PaymentResponse)
	case xdr.OperationTypeReviewPaymentRequest:
		is.updateIngestedPaymentRequest(*is.Cursor.Operation(), is.Cursor.OperationSourceAccount())
		is.updateIngestedPayment(*is.Cursor.Operation(), is.Cursor.OperationSourceAccount(), *is.Cursor.OperationResult())
	case xdr.OperationTypeDemurrage:
		is.processDemurrage(*is.Cursor.OperationResult())
	case xdr.OperationTypeDirectDebit:
		opDirectDebit := is.Cursor.Operation().Body.MustDirectDebitOp()
		is.processPayment(opDirectDebit.PaymentOp,
			opDirectDebit.From,
			is.Cursor.OperationResult().MustDirectDebitResult().MustSuccess().PaymentResponse)
	case xdr.OperationTypeManageOffer:
		is.manageOffer(is.Cursor.OperationSourceAccount(), is.Cursor.OperationResult().MustManageOfferResult())
	case xdr.OperationTypeManageInvoice:
		is.processManageInvoice(is.Cursor.Operation().Body.MustManageInvoiceOp(),
			is.Cursor.OperationResult().MustManageInvoiceResult())

	}
}
