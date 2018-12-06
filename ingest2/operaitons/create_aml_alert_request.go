package operaitons

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type createAMLAlertReqeustOpHandler struct {
	balanceProvider balanceProvider
}

func (h *createAMLAlertReqeustOpHandler) OperationDetails(op rawOperation,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	amlAlertRequest := op.Body.MustCreateAmlAlertRequestOp().AmlAlertRequest

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreateAmlAlert,
		CreateAMLAlertRequest: &history2.CreateAMLAlertRequestDetails{
			Amount:    amount.StringU(uint64(amlAlertRequest.Amount)),
			BalanceID: amlAlertRequest.BalanceId.AsString(),
			Reason:    string(amlAlertRequest.Reason),
		},
	}, nil
}

func (h *createAMLAlertReqeustOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	amlAlertRequest := opBody.MustCreateAmlAlertRequestOp().AmlAlertRequest

	balance := h.balanceProvider.GetBalanceByID(amlAlertRequest.BalanceId)

	assetCode := xdr.AssetCode(balance.AssetCode)

	effect := history2.Effect{
		Type: history2.EffectTypeLocked,
		Locked: &history2.LockedEffect{
			Amount: amount.String(int64(amlAlertRequest.Amount)),
		},
	}

	var participants []history2.ParticipantEffect

	if balance.AccountID == source.AccountID {
		source.BalanceID = &balance.ID
		source.AssetCode = &assetCode
		source.Effect = effect
	} else {
		participants = append(participants, history2.ParticipantEffect{
			AccountID: balance.AccountID,
			BalanceID: &balance.ID,
			AssetCode: &assetCode,
			Effect:    effect,
		})
	}

	return append(participants, source), nil
}
