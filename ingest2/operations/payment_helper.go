package operations

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

// PaymentHelper used to get info about payment operation changes - payment and review invoice
type PaymentHelper struct {
	pubKeyProvider publicKeyProvider
}

// NewPaymentHelper returns new PaymentHelper
func NewPaymentHelper(provider publicKeyProvider) PaymentHelper {
	return PaymentHelper{
		pubKeyProvider: provider,
	}
}

// GetParticipantsEffects returns participants effects of payment transfer
func (h *PaymentHelper) GetParticipantsEffects(op xdr.PaymentOpV2,
	res xdr.PaymentV2Response, source history2.ParticipantEffect, destinationEffectType history2.EffectType,
) ([]history2.ParticipantEffect, error) {
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

	sourceBalanceID := h.pubKeyProvider.GetBalanceID(op.SourceBalanceId)
	source.BalanceID = &sourceBalanceID
	source.AssetCode = &res.Asset
	source.Effect = history2.Effect{
		Type: history2.EffectTypeCharged,
		Charged: &history2.ChargedEffect{
			Amount: amount.StringU(uint64(op.Amount)),
			FeePaid: history2.FeePaid{
				Fixed:             amount.StringU(uint64(sourceFixedFee)),
				CalculatedPercent: amount.StringU(uint64(sourcePercentFee)),
			},
		},
	}

	destBalanceID := h.pubKeyProvider.GetBalanceID(res.DestinationBalanceId)
	destination := history2.ParticipantEffect{
		AccountID: h.pubKeyProvider.GetAccountID(res.Destination),
		BalanceID: &destBalanceID,
		AssetCode: &res.Asset,
		Effect: history2.Effect{
			Type: destinationEffectType,
		},
	}

	switch destination.Effect.Type {
	case history2.EffectTypeFunded:
		destination.Effect.Funded = &history2.FundedEffect{
			Amount: amount.StringU(uint64(op.Amount)),
			FeePaid: history2.FeePaid{
				Fixed:             amount.StringU(uint64(destFixedFee)),
				CalculatedPercent: amount.StringU(uint64(destPercentFee)),
			},
		}
	case history2.EffectTypeFundedToLocked:
		if (destFixedFee != 0) || (destPercentFee != 0) {
			return nil, errors.New("unexpected state, expected zero fee if destination amount locked")
		}

		destination.Effect.FundedToLocked = &history2.FundedToLockedEffect{
			Amount: amount.StringU(uint64(op.Amount)),
		}
	}

	return []history2.ParticipantEffect{source, destination}, nil
}
