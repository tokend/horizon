package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/v2"
)

type createAMLAlertReqeustOpHandler struct {
	balanceProvider balanceProvider
}

// Details returns details about create AML alert request operation
func (h *createAMLAlertReqeustOpHandler) Details(op rawOperation,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	amlAlertRequest := op.Body.MustCreateAmlAlertRequestOp().AmlAlertRequest

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreateAmlAlert,
		CreateAMLAlertRequest: &history2.CreateAMLAlertRequestDetails{
			Amount:         regources.Amount(amlAlertRequest.Amount),
			BalanceAddress: amlAlertRequest.BalanceId.AsString(),
			Reason:         string(amlAlertRequest.Reason),
		},
	}, nil
}

// ParticipantsEffects returns `locked` effect for account
// which is suspected in illegal obtaining of tokens
func (h *createAMLAlertReqeustOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	amlAlertRequest := opBody.MustCreateAmlAlertRequestOp().AmlAlertRequest

	effect := history2.Effect{
		Type: history2.EffectTypeLocked,
		Locked: &history2.BalanceChangeEffect{
			Amount: regources.Amount(amlAlertRequest.Amount),
		},
	}

	balance := h.balanceProvider.MustBalance(amlAlertRequest.BalanceId)
	return populateEffects(balance, effect, source), nil
}
