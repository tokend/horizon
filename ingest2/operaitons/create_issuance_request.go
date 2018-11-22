package operaitons

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/utf8"
)

type createIssuanceRequestOpHandler struct {
	pubKeyProvider  publicKeyProvider
	balanceProvider balanceProvider
}

func (h *createIssuanceRequestOpHandler) OperationDetails(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	createIssuanceRequestOp := opBody.MustCreateIssuanceRequestOp()
	issuanceRequest := createIssuanceRequestOp.Request

	var allTasks *int64
	rawAllTasks, ok := createIssuanceRequestOp.Ext.GetAllTasks()
	if ok && rawAllTasks != nil {
		allTasksInt := int64(*rawAllTasks)
		allTasks = &allTasksInt
	}

	balanceID := h.pubKeyProvider.GetBalanceID(issuanceRequest.Receiver)

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreateIssuanceRequest,
		CreateIssuanceRequest: &history2.CreateIssuanceRequestDetails{
			FixedFee:          amount.StringU(uint64(issuanceRequest.Fee.Fixed)),
			PercentFee:        amount.StringU(uint64(issuanceRequest.Fee.Percent)),
			Reference:         utf8.Scrub(string(createIssuanceRequestOp.Reference)),
			Amount:            amount.StringU(uint64(issuanceRequest.Amount)),
			Asset:             issuanceRequest.Asset,
			ReceiverAccountID: h.balanceProvider.GetBalanceByID(balanceID).AccountID,
			ReceiverBalanceID: balanceID,
			ExternalDetails:   customDetailsUnmarshal([]byte(issuanceRequest.ExternalDetails)),
			AllTasks:          allTasks,
		},
	}, nil
}

func (h *createIssuanceRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	issuanceRequest := opBody.MustCreateIssuanceRequestOp().Request

	balanceID := h.pubKeyProvider.GetBalanceID(issuanceRequest.Receiver)
	balance := h.balanceProvider.GetBalanceByID(balanceID)

	issuanceAmount := int64(issuanceRequest.Amount)

	effect := history2.Effect{
		Type:           history2.EffectTypeIssuance,
		IssuanceAmount: &issuanceAmount,
	}

	var participants []history2.ParticipantEffect

	if balance.AccountID == source.AccountID {
		source.BalanceID = &balanceID
		source.AssetCode = &issuanceRequest.Asset
		source.Effect = effect
	} else {
		participants = append(participants, history2.ParticipantEffect{
			AccountID: balance.AccountID,
			BalanceID: &balanceID,
			AssetCode: &issuanceRequest.Asset,
			Effect:    effect,
		})
	}

	participants = append(participants, source)

	return participants, nil
}
