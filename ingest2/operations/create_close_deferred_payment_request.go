package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ingest2/internal"
	regources "gitlab.com/tokend/regources/generated"
)

type createCloseDeferredPaymentRequestOpHandler struct {
	effectsProvider
}

// Details returns details about manage balance operation
func (h *createCloseDeferredPaymentRequestOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	oper := op.Body.MustCreateCloseDeferredPaymentRequestOp()
	operRes := opRes.MustCreateCloseDeferredPaymentRequestResult().MustSuccess()
	details := history2.OperationDetails{
		Type: xdr.OperationTypeCreateCloseDeferredPaymentRequest,
		CreateCloseDeferredPaymentRequest: &history2.CreateCloseDeferredPaymentRequest{
			RequestID:          uint64(oper.RequestId),
			DestinationBalance: operRes.ExtendedResult.DestinationBalance.AsString(),
			DestinationAccount: operRes.ExtendedResult.Destination.Address(),
			Amount:             regources.Amount(oper.Request.Amount),
			Details:            internal.MarshalCustomDetails(oper.Request.CreatorDetails),
			AllTasks:           (*uint32)(oper.AllTasks),
		},
	}

	return details, nil
}

// ParticipantsEffects returns `funded` and `charged` effects
func (h *createCloseDeferredPaymentRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{h.Participant(sourceAccountID)}, nil
}
