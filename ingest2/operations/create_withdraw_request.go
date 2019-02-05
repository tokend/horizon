package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ingest2/internal"
	"gitlab.com/tokend/regources/v2"
)

type createWithdrawRequestOpHandler struct {
	pubKeyProvider IDProvider
}

// CreatorDetails returns details about create withdraw request operation
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
	_ xdr.OperationResultTr, source history2.ParticipantEffect, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	withdrawRequest := opBody.MustCreateWithdrawalRequestOp().Request
	balanceIDInt := h.pubKeyProvider.MustBalanceID(withdrawRequest.Balance)

	source.BalanceID = &balanceIDInt
	source.Effect.Type = history2.EffectTypeLocked
	source.Effect.Locked = &history2.BalanceChangeEffect{
		Amount: regources.Amount(withdrawRequest.Amount),
		Fee:    internal.FeeFromXdr(withdrawRequest.Fee),
	}

	return []history2.ParticipantEffect{source}, nil
}
