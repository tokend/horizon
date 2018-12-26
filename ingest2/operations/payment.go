package operations

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/utf8"
)

type paymentOpHandler struct {
	pubKeyProvider IDProvider
}

// Details returns details about payment operation
func (h *paymentOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	paymentOp := op.Body.MustPaymentOpV2()
	paymentRes := opRes.MustPaymentV2Result().MustPaymentV2Response()

	return history2.OperationDetails{
		Type: xdr.OperationTypePaymentV2,
		Payment: &history2.PaymentDetails{
			AccountFrom: op.Source.Address(),
			AccountTo:   paymentRes.Destination.Address(),
			BalanceFrom: paymentOp.SourceBalanceId.AsString(),
			BalanceTo:   paymentRes.DestinationBalanceId.AsString(),
			Amount:      amount.StringU(uint64(paymentOp.Amount)),
			Asset:       paymentRes.Asset,
			SourceFeeData: history2.FeeData{
				FixedFee:  amount.StringU(uint64(paymentRes.ActualSourcePaymentFee.Fixed)),
				ActualFee: amount.StringU(uint64(paymentRes.ActualSourcePaymentFee.Percent)),
			},
			DestinationFeeData: history2.FeeData{
				FixedFee:  amount.StringU(uint64(paymentRes.ActualDestinationPaymentFee.Fixed)),
				ActualFee: amount.StringU(uint64(paymentRes.ActualDestinationPaymentFee.Percent)),
			},
			SourcePayForDestination: paymentOp.FeeData.SourcePaysForDest,
			Subject:                 string(paymentOp.Subject),
			Reference:               utf8.Scrub(string(paymentOp.Reference)),
			UniversalAmount:         amount.StringU(uint64(paymentRes.SourceSentUniversal)),
		},
	}, nil
}

// ParticipantsEffects returns `funded` and `charged` effects
func (h *paymentOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	op := opBody.MustPaymentOpV2()
	res := opRes.MustPaymentV2Result().MustPaymentV2Response()

	sourceFixedFee := res.ActualSourcePaymentFee.Fixed
	sourcePercentFee := res.ActualSourcePaymentFee.Percent
	destFixedFee := res.ActualDestinationPaymentFee.Fixed
	destPercentFee := res.ActualDestinationPaymentFee.Percent
	if op.FeeData.SourcePaysForDest {
		sourceFixedFee += destFixedFee
		destFixedFee = 0
		sourcePercentFee += destPercentFee
		destPercentFee = 0
	}

	sourceBalanceID := h.pubKeyProvider.MustBalanceID(op.SourceBalanceId)
	source.BalanceID = &sourceBalanceID
	source.AssetCode = new(string)
	*source.AssetCode = string(res.Asset)
	source.Effect = history2.Effect{
		Type: history2.EffectTypeCharged,
		Charged: &history2.BalanceChangeEffect{
			Amount: amount.StringU(uint64(op.Amount)),
			Fee: history2.Fee{
				Fixed:             amount.StringU(uint64(sourceFixedFee)),
				CalculatedPercent: amount.StringU(uint64(sourcePercentFee)),
			},
		},
	}

	destBalanceID := h.pubKeyProvider.MustBalanceID(res.DestinationBalanceId)
	destination := history2.ParticipantEffect{
		AccountID: h.pubKeyProvider.MustAccountID(res.Destination),
		BalanceID: &destBalanceID,
		AssetCode: source.AssetCode,
		Effect: history2.Effect{
			Type: history2.EffectTypeFunded,
			Funded: &history2.BalanceChangeEffect{
				Amount: amount.StringU(uint64(op.Amount)),
				Fee: history2.Fee{
					Fixed:             amount.StringU(uint64(destFixedFee)),
					CalculatedPercent: amount.StringU(uint64(destPercentFee)),
				},
			},
		},
	}

	return []history2.ParticipantEffect{source, destination}, nil
}
