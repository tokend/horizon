package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type createDeferredPaymentCreationRequestOpHandler struct {
	effectsProvider
}

// Details returns details about manage balance operation
func (h *createDeferredPaymentCreationRequestOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	//oper := op.Body.MustCreateDeferredPaymentCreationRequestOp()
	details := history2.OperationDetails{
		Type: xdr.OperationTypeCreateDeferredPaymentCreationRequest,
		//CreateDeferredPaymentCreationRequest: &history2.CreateDeferredPaymentCreationRequest{
		//	RequestID:               uint64(oper.RequestId),
		//	SourceBalance:           oper.Request.SourceBalance.AsString(),
		//	DestinationAccount:      oper.Request.Destination.Address(),
		//	Amount:                  regources.Amount(oper.Request.Amount),
		//	SourcePayForDestination: oper.Request.FeeData.SourcePaysForDest,
		//	SourceFee: regources.Fee{
		//		CalculatedPercent: regources.Amount(oper.Request.FeeData.SourceFee.Percent),
		//		Fixed:             regources.Amount(oper.Request.FeeData.SourceFee.Fixed),
		//	},
		//	DestinationFee: regources.Fee{
		//		CalculatedPercent: regources.Amount(oper.Request.FeeData.DestinationFee.Percent),
		//		Fixed:             regources.Amount(oper.Request.FeeData.DestinationFee.Fixed),
		//	},
		//	Details:  internal.MarshalCustomDetails(oper.Request.CreatorDetails),
		//	AllTasks: (*uint32)(oper.AllTasks),
		//},
	}

	return details, nil
}

// ParticipantsEffects returns `funded` and `charged` effects
func (h *createDeferredPaymentCreationRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{h.Participant(sourceAccountID)}, nil
}
