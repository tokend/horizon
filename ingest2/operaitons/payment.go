package operaitons

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/utf8"
)

type paymentOpHandler struct {
	pubKeyProvider publicKeyProvider
}

func (h *paymentOpHandler) OperationDetails(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	paymentOp := op.Body.MustPaymentOpV2()
	paymentRes := opRes.MustPaymentV2Result().MustPaymentV2Response()

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreateAccount,
		Payment: &history2.PaymentDetails{
			AccountFrom: op.Source.Address(),
			AccountTo:   paymentRes.Destination.Address(),
			BalanceFrom: paymentOp.SourceBalanceId.AsString(),
			BalanceTo:   paymentRes.DestinationBalanceId.AsString(),
			Amount:      amount.StringU(uint64(paymentOp.Amount)),
			Asset:       paymentRes.Asset,
			SourceFeeData: history2.FeeData{
				FixedFee:  amount.StringU(uint64(paymentOp.FeeData.SourceFee.FixedFee)),
				ActualFee: amount.StringU(uint64(paymentRes.ActualSourcePaymentFee)),
			},
			DestinationFeeData: history2.FeeData{
				FixedFee:  amount.StringU(uint64(paymentOp.FeeData.DestinationFee.FixedFee)),
				ActualFee: amount.StringU(uint64(paymentRes.ActualDestinationPaymentFee)),
			},
			SourcePayForDestination: paymentOp.FeeData.SourcePaysForDest,
			Subject:                 string(paymentOp.Subject),
			Reference:               utf8.Scrub(string(paymentOp.Reference)),
			UniversalAmount:         amount.StringU(uint64(paymentRes.SourceSentUniversal)),
		},
	}, nil
}

func (h *paymentOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	paymentOp := opBody.MustPaymentOpV2()
	paymentRes := opRes.MustPaymentV2Result().PaymentV2Response

	sourceBalanceID :=  h.pubKeyProvider.GetBalanceID(paymentOp.SourceBalanceId)
	source.BalanceID = &sourceBalanceID
	source.AssetCode = &paymentRes.Asset
	source.Effect = history2.Effect{
		Type: history2.EffectTypeCharged,
		Payment: &history2.PaymentEffect{
			Amount: amount.StringU(uint64(paymentOp.Amount)),
		},
	}

	destBalanceID := h.pubKeyProvider.GetBalanceID(paymentRes.DestinationBalanceId),

	return []history2.ParticipantEffect{source, {
		AccountID: h.pubKeyProvider.GetAccountID(paymentRes.Destination),
		BalanceID: &destBalanceID,
		AssetCode: &paymentRes.Asset,
		Effect: history2.Effect{
			Type: history2.EffectTypeFunded,
			Payment: &history2.PaymentEffect{
				Amount: amount.StringU(uint64(paymentOp.Amount)),
			},
		},
	}}, nil
}
