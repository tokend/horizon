package operations

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type createWithdrawRequestOpHandler struct {
	pubKeyProvider publicKeyProvider
}

// Details returns details about create withdraw request operation
func (h *createWithdrawRequestOpHandler) Details(op RawOperation,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	withdrawRequest := op.Body.MustCreateWithdrawalRequestOp().Request

	destinationAsset := xdr.AssetCode("")
	destinationAmount := amount.String(int64(0))
	if autoConversion, ok := withdrawRequest.Details.GetAutoConversion(); ok {
		destinationAsset = autoConversion.DestAsset
		destinationAmount = amount.StringU(uint64(autoConversion.ExpectedAmount))
	}

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreateWithdrawalRequest,
		CreateWithdrawRequest: &history2.CreateWithdrawRequestDetails{
			BalanceAddress:    withdrawRequest.Balance.AsString(),
			Amount:            amount.StringU(uint64(withdrawRequest.Amount)),
			FixedFee:          amount.String(int64(withdrawRequest.Fee.Fixed)),
			PercentFee:        amount.String(int64(withdrawRequest.Fee.Percent)),
			ExternalDetails:   customDetailsUnmarshal([]byte(withdrawRequest.ExternalDetails)),
			DestinationAsset:  destinationAsset,
			DestinationAmount: destinationAmount,
		},
	}, nil
}

// ParticipantsEffects returns source `locked` effect
func (h *createWithdrawRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	_ xdr.OperationResultTr, source history2.ParticipantEffect, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	withdrawRequest := opBody.MustCreateWithdrawalRequestOp().Request
	balanceIDInt := h.pubKeyProvider.GetBalanceID(withdrawRequest.Balance)

	source.BalanceID = &balanceIDInt
	source.Effect.Type = history2.EffectTypeLocked
	source.Effect.Locked = &history2.LockedEffect{
		Amount: amount.StringU(uint64(withdrawRequest.Amount)),
		FeeLocked: history2.FeePaid{
			Fixed:             amount.StringU(uint64(withdrawRequest.Fee.Fixed)),
			CalculatedPercent: amount.StringU(uint64(withdrawRequest.Fee.Percent)),
		},
	}

	return []history2.ParticipantEffect{source}, nil
}
