package operaitons

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type setFeeOpHandler struct {
	pubKeyProvider publicKeyProvider
}

func (h *setFeeOpHandler) OperationDetails(op rawOperation, _ xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	setFeeOp := op.Body.MustSetFeesOp()

	if setFeeOp.IsDelete || setFeeOp.Fee == nil {
		return history2.OperationDetails{
			Type:   xdr.OperationTypeSetFees,
			SetFee: &history2.SetFeeDetails{},
		}, nil
	}

	fee := *setFeeOp.Fee

	return history2.OperationDetails{
		Type: xdr.OperationTypeSetFees,
		SetFee: &history2.SetFeeDetails{
			AssetCode:   fee.Asset,
			FixedFee:    amount.String(int64(fee.FixedFee)),
			PercentFee:  amount.String(int64(fee.PercentFee)),
			FeeType:     fee.FeeType,
			AccountID:   fee.AccountId.Address(),
			AccountType: fee.AccountType,
			Subtype:     int64(fee.Subtype),
			LowerBound:  amount.String(int64(fee.LowerBound)),
			UpperBound:  amount.String(int64(fee.UpperBound)),
		},
	}, nil
}

func (h *setFeeOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	_ xdr.OperationResultTr, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	participants := []history2.ParticipantEffect{source}

	setFeeOp := opBody.MustSetFeesOp()
	if (setFeeOp.Fee != nil) && (setFeeOp.Fee.AccountId != nil) {
		participants = append(participants, history2.ParticipantEffect{
			AccountID: h.pubKeyProvider.GetAccountID(*setFeeOp.Fee.AccountId),
		})
	}

	return participants, nil
}
