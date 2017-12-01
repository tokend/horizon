package ingest

import (
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon/db2/history"
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
	case xdr.OperationTypeManageForfeitRequest:
		state = history.OperationStatePending
		manageRequestResult := operationResult.MustManageForfeitRequestResult()
		operationIdentifier = uint64(manageRequestResult.Success.PaymentId)
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
	case xdr.OperationTypeManageForfeitRequest:
		is.processManageForfeitRequest(*is.Cursor.Operation(), is.Cursor.OperationSourceAccount(), *is.Cursor.OperationResult())
	case xdr.OperationTypePayment:
		is.processPayment(is.Cursor.Operation().Body.MustPaymentOp(), is.Cursor.OperationSourceAccount(),
			*is.Cursor.OperationResult().MustPaymentResult().PaymentResponse)
	case xdr.OperationTypeReviewPaymentRequest:
		is.updateIngestedPaymentRequest(*is.Cursor.Operation(), is.Cursor.OperationSourceAccount())
		is.updateIngestedPayment(*is.Cursor.Operation(), is.Cursor.OperationSourceAccount(), *is.Cursor.OperationResult())
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
	case xdr.OperationTypeReviewRequest:
		is.processReviewRequest(is.Cursor.Operation().Body.MustReviewRequestOp())
	case xdr.OperationTypeManageAsset:
		is.processManageAsset(is.Cursor.Operation().Body.ManageAssetOp)
	}
}
