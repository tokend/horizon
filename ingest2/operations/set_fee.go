package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/v2"
)

type setFeeOpHandler struct {
	pubKeyProvider IDProvider
}

// Details returns details about set fee operation
func (h *setFeeOpHandler) Details(op rawOperation, _ xdr.OperationResultTr,
) (regources.OperationDetails, error) {
	setFeeOp := op.Body.MustSetFeesOp()

	if setFeeOp.IsDelete || setFeeOp.Fee == nil {
		return regources.OperationDetails{
			Type:   xdr.OperationTypeSetFees,
			SetFee: &regources.SetFeeDetails{},
		}, nil
	}

	fee := *setFeeOp.Fee

	var feeAccountAddress *string
	if fee.AccountId != nil {
		feeAccountAddress = new(string)
		*feeAccountAddress = fee.AccountId.Address()
	}

	return regources.OperationDetails{
		Type: xdr.OperationTypeSetFees,
		SetFee: &regources.SetFeeDetails{
			AssetCode:      string(fee.Asset),
			FixedFee:       regources.Amount(fee.FixedFee),
			PercentFee:     regources.Amount(fee.PercentFee),
			FeeType:        fee.FeeType,
			AccountAddress: feeAccountAddress,
			AccountType:    fee.AccountType,
			Subtype:        int64(fee.Subtype),
			LowerBound:     regources.Amount(fee.LowerBound),
			UpperBound:     regources.Amount(fee.UpperBound),
		},
	}, nil
}

//ParticipantsEffects - returns source participant and counterparty for which fee has been set if one explicitly
// specified
func (h *setFeeOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	_ xdr.OperationResultTr, source history2.ParticipantEffect, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	participants := []history2.ParticipantEffect{source}

	setFeeOp := opBody.MustSetFeesOp()
	if (setFeeOp.Fee != nil) && (setFeeOp.Fee.AccountId != nil) {
		participants = append(participants, history2.ParticipantEffect{
			AccountID: h.pubKeyProvider.MustAccountID(*setFeeOp.Fee.AccountId),
		})
	}

	return participants, nil
}
