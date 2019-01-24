package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ingest2/internal"
	"gitlab.com/tokend/horizon/utf8"
	"gitlab.com/tokend/regources/v2"
)

type paymentOpHandler struct {
	pubKeyProvider IDProvider
}

// Details returns details about payment operation
func (h *paymentOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (regources.OperationDetails, error) {
	paymentOp := op.Body.MustPaymentOpV2()
	paymentRes := opRes.MustPaymentV2Result().MustPaymentV2Response()

	return regources.OperationDetails{
		Type: xdr.OperationTypePaymentV2,
		Payment: &regources.PaymentDetails{
			AccountFrom:             op.Source.Address(),
			AccountTo:               paymentRes.Destination.Address(),
			BalanceFrom:             paymentOp.SourceBalanceId.AsString(),
			BalanceTo:               paymentRes.DestinationBalanceId.AsString(),
			Amount:                  regources.Amount(paymentOp.Amount),
			Asset:                   string(paymentRes.Asset),
			SourceFee:               internal.FeeFromXdr(paymentRes.ActualSourcePaymentFee),
			DestinationFee:          internal.FeeFromXdr(paymentRes.ActualDestinationPaymentFee),
			SourcePayForDestination: paymentOp.FeeData.SourcePaysForDest,
			Subject:                 string(paymentOp.Subject),
			Reference:               utf8.Scrub(string(paymentOp.Reference)),
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
	source.Effect = &regources.Effect{
		Type: regources.EffectTypeCharged,
		Charged: &regources.BalanceChangeEffect{
			Amount: regources.Amount(op.Amount),
			Fee: regources.Fee{
				Fixed:             regources.Amount(sourceFixedFee),
				CalculatedPercent: regources.Amount(sourcePercentFee),
			},
		},
	}

	destBalanceID := h.pubKeyProvider.MustBalanceID(res.DestinationBalanceId)
	destination := history2.ParticipantEffect{
		AccountID: h.pubKeyProvider.MustAccountID(res.Destination),
		BalanceID: &destBalanceID,
		AssetCode: source.AssetCode,
		Effect: &regources.Effect{
			Type: regources.EffectTypeFunded,
			Funded: &regources.BalanceChangeEffect{
				Amount: regources.Amount(op.Amount),
				Fee: regources.Fee{
					Fixed:             regources.Amount(destFixedFee),
					CalculatedPercent: regources.Amount(destPercentFee),
				},
			},
		},
	}

	return []history2.ParticipantEffect{source, destination}, nil
}
