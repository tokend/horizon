package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type cancelDeferredPaymentCreationRequestOpHandler struct {
	effectsProvider
}

// Details returns details about manage balance operation
func (h *cancelDeferredPaymentCreationRequestOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	cancelDeferredPaymentCreationCreationRequestOp := op.Body.MustCancelDeferredPaymentCreationRequestOp()
	details := history2.OperationDetails{
		Type: xdr.OperationTypeCancelDeferredPaymentCreationRequest,
		CancelDeferredPaymentCreationRequest: &history2.CancelDeferredPaymentCreationRequest{
			RequestID: uint64(cancelDeferredPaymentCreationCreationRequestOp.RequestId),
		},
	}

	return details, nil
}

// ParticipantsEffects returns `funded` and `charged` effects
func (h *cancelDeferredPaymentCreationRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{h.Participant(sourceAccountID)}, nil
}
