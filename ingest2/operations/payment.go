package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ingest2/internal"
	"gitlab.com/tokend/horizon/utf8"
	"gitlab.com/tokend/regources/rgenerated"
)

type paymentOpHandler struct {
	effectsProvider
}

// Details returns details about payment operation
func (h *paymentOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	paymentOp := op.Body.MustPaymentOp()
	paymentRes := opRes.MustPaymentResult().MustPaymentResponse()

	return history2.OperationDetails{
		Type: xdr.OperationTypePayment,
		Payment: &history2.PaymentDetails{
			AccountFrom:             op.Source.Address(),
			AccountTo:               paymentRes.Destination.Address(),
			BalanceFrom:             paymentOp.SourceBalanceId.AsString(),
			BalanceTo:               paymentRes.DestinationBalanceId.AsString(),
			Amount:                  rgenerated.Amount(paymentOp.Amount),
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
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	op := opBody.MustPaymentOp()
	res := opRes.MustPaymentResult().MustPaymentResponse()

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

	source := h.BalanceEffect(op.SourceBalanceId, &history2.Effect{
		Type: history2.EffectTypeCharged,
		Charged: &history2.BalanceChangeEffect{
			Amount: rgenerated.Amount(op.Amount),
			Fee: rgenerated.Fee{
				Fixed:             rgenerated.Amount(sourceFixedFee),
				CalculatedPercent: rgenerated.Amount(sourcePercentFee),
			},
		},
	})

	destination := h.BalanceEffect(res.DestinationBalanceId, &history2.Effect{
		Type: history2.EffectTypeFunded,
		Funded: &history2.BalanceChangeEffect{
			Amount: rgenerated.Amount(op.Amount),
			Fee: rgenerated.Fee{
				Fixed:             rgenerated.Amount(destFixedFee),
				CalculatedPercent: rgenerated.Amount(destPercentFee),
			},
		},
	})
	return []history2.ParticipantEffect{source, destination}, nil
}
