package operaitons

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/utf8"
)

type createIssuanceRequestOpHandler struct {
	pubKeyProvider publicKeyProvider
}

// OperationDetails returns details about create issuance request operation
func (h *createIssuanceRequestOpHandler) OperationDetails(op rawOperation,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	createIssuanceRequestOp := op.Body.MustCreateIssuanceRequestOp()
	issuanceRequest := createIssuanceRequestOp.Request

	var allTasks *int64
	rawAllTasks, ok := createIssuanceRequestOp.Ext.GetAllTasks()
	if ok && rawAllTasks != nil {
		allTasksInt := int64(*rawAllTasks)
		allTasks = &allTasksInt
	}

	createIssuanceRequestRes := opRes.MustCreateIssuanceRequestResult().MustSuccess()

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreateIssuanceRequest,
		CreateIssuanceRequest: &history2.CreateIssuanceRequestDetails{
			FixedFee:          amount.StringU(uint64(issuanceRequest.Fee.Fixed)),
			PercentFee:        amount.StringU(uint64(issuanceRequest.Fee.Percent)),
			Reference:         utf8.Scrub(string(createIssuanceRequestOp.Reference)),
			Amount:            amount.StringU(uint64(issuanceRequest.Amount)),
			Asset:             issuanceRequest.Asset,
			ReceiverAccountID: createIssuanceRequestRes.Receiver.Address(),
			ReceiverBalanceID: issuanceRequest.Receiver.AsString(),
			ExternalDetails:   customDetailsUnmarshal([]byte(issuanceRequest.ExternalDetails)),
			AllTasks:          allTasks,
			RequestDetails: history2.RequestDetails{
				IsFulfilled: createIssuanceRequestRes.Fulfilled,
			},
		},
	}, nil
}

// ParticipantsEffects returns source `funded` effect if request was fulfilled
func (h *createIssuanceRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	issuanceRequest := opBody.MustCreateIssuanceRequestOp().Request
	createIssuanceRequestRes := opRes.MustCreateIssuanceRequestResult().MustSuccess()

	if !createIssuanceRequestRes.Fulfilled {
		return []history2.ParticipantEffect{source}, nil
	}

	receiverID := h.pubKeyProvider.GetAccountID(createIssuanceRequestRes.Receiver)
	receiverBalanceID := h.pubKeyProvider.GetBalanceID(issuanceRequest.Receiver)

	effect := history2.Effect{
		Type: history2.EffectTypeFunded,
		Funded: &history2.FundedEffect{
			Amount: amount.String(int64(issuanceRequest.Amount - issuanceRequest.Fee.Fixed - issuanceRequest.Fee.Percent)),
			FeePaid: history2.FeePaid{
				Fixed:             amount.StringU(uint64(issuanceRequest.Fee.Fixed)),
				CalculatedPercent: amount.StringU(uint64(issuanceRequest.Fee.Percent)),
			},
		},
	}

	var participants []history2.ParticipantEffect

	if receiverID == source.AccountID {
		source.BalanceID = &receiverBalanceID
		source.AssetCode = &issuanceRequest.Asset
		source.Effect = effect
	} else {
		participants = append(participants, history2.ParticipantEffect{
			AccountID: receiverID,
			BalanceID: &receiverBalanceID,
			AssetCode: &issuanceRequest.Asset,
			Effect:    effect,
		})
	}

	return append(participants, source), nil
}
