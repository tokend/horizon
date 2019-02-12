package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ingest2/internal"
	"gitlab.com/tokend/horizon/utf8"
	"gitlab.com/tokend/regources/v2"
)

type createIssuanceRequestOpHandler struct {
	effectsProvider
}

// Details returns details about create issuance request operation
func (h *createIssuanceRequestOpHandler) Details(op rawOperation,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	createIssuanceRequestOp := op.Body.MustCreateIssuanceRequestOp()
	issuanceRequest := createIssuanceRequestOp.Request

	var allTasks *int64
	rawAllTasks := createIssuanceRequestOp.AllTasks
	if rawAllTasks != nil {
		allTasksInt := int64(*rawAllTasks)
		allTasks = &allTasksInt
	}

	createIssuanceRequestRes := opRes.MustCreateIssuanceRequestResult().MustSuccess()

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreateIssuanceRequest,
		CreateIssuanceRequest: &history2.CreateIssuanceRequestDetails{
			Fee:       internal.FeeFromXdr(issuanceRequest.Fee),
			Reference: utf8.Scrub(string(createIssuanceRequestOp.Reference)),
			Amount:    regources.Amount(issuanceRequest.Amount),
			Asset:     string(issuanceRequest.Asset),
			ReceiverAccountAddress: createIssuanceRequestRes.Receiver.Address(),
			ReceiverBalanceAddress: issuanceRequest.Receiver.AsString(),
			ExternalDetails:        internal.MarshalCustomDetails(issuanceRequest.ExternalDetails),
			AllTasks:               allTasks,
			RequestDetails: history2.RequestDetails{
				IsFulfilled: createIssuanceRequestRes.Fulfilled,
			},
		},
	}, nil
}

// ParticipantsEffects returns source `funded` effect if request was fulfilled
func (h *createIssuanceRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	issuanceRequest := opBody.MustCreateIssuanceRequestOp().Request
	createIssuanceRequestRes := opRes.MustCreateIssuanceRequestResult().MustSuccess()

	if !createIssuanceRequestRes.Fulfilled {
		return []history2.ParticipantEffect{h.Participant(sourceAccountID)}, nil
	}

	return h.BalanceEffectWithAccount(sourceAccountID, issuanceRequest.Receiver, &history2.Effect{
		Type: history2.EffectTypeIssued,
		Issued: &history2.BalanceChangeEffect{
			Amount: regources.Amount(issuanceRequest.Amount),
			Fee:    internal.FeeFromXdr(issuanceRequest.Fee),
		},
	}), nil
}
