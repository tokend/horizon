package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ingest2/internal"
	"gitlab.com/tokend/regources/v2"
)

type createWithdrawRequestOpHandler struct {
	effectsProvider
}

// Details returns details about create withdraw request operation
func (h *createWithdrawRequestOpHandler) Details(op rawOperation,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	withdrawRequest := op.Body.MustCreateWithdrawalRequestOp().Request

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreateWithdrawalRequest,
		CreateWithdrawRequest: &history2.CreateWithdrawRequestDetails{
			BalanceAddress:  withdrawRequest.Balance.AsString(),
			Amount:          regources.Amount(withdrawRequest.Amount),
			Fee:             internal.FeeFromXdr(withdrawRequest.Fee),
			ExternalDetails: internal.MarshalCustomDetails(withdrawRequest.CreatorDetails),
		},
	}, nil
}

// ParticipantsEffects returns source `locked` effect
func (h *createWithdrawRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	_ xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	withdrawRequest := opBody.MustCreateWithdrawalRequestOp().Request

	return []history2.ParticipantEffect{h.BalanceEffect(withdrawRequest.Balance, &history2.Effect{
		Type: history2.EffectTypeLocked,
		Locked: &history2.BalanceChangeEffect{
			Amount: regources.Amount(withdrawRequest.Amount),
			Fee:    internal.FeeFromXdr(withdrawRequest.Fee),
		},
	})}, nil
}
