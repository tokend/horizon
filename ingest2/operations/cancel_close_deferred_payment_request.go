package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type cancelCloseDeferredPaymentRequestOpHandler struct {
	effectsProvider
}

// Details returns details about manage balance operation
func (h *cancelCloseDeferredPaymentRequestOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	cancelCloseDeferredPaymentCreationRequestOp := op.Body.MustCancelCloseDeferredPaymentRequestOp()
	details := history2.OperationDetails{
		Type: xdr.OperationTypeCancelCloseDeferredPaymentRequest,
		CancelCloseDeferredPaymentRequest: &history2.CancelCloseDeferredPaymentRequest{
			RequestID: uint64(cancelCloseDeferredPaymentCreationRequestOp.RequestId),
		},
	}

	return details, nil
}

// ParticipantsEffects returns `funded` and `charged` effects
func (h *cancelCloseDeferredPaymentRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{h.Participant(sourceAccountID)}, nil
}
