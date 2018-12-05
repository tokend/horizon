package operaitons

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type paymentHelper struct {
	pubKeyProvider publicKeyProvider
}

// TODO handle fee after payment response will have normal info about fee
func (h *paymentHelper) getParticipantsEffects(op xdr.PaymentOpV2,
	res xdr.PaymentV2Response, source history2.ParticipantEffect, destinationEffectType history2.EffectType,
) []history2.ParticipantEffect {
	sourceBalanceID := h.pubKeyProvider.GetBalanceID(op.SourceBalanceId)
	source.BalanceID = &sourceBalanceID
	source.AssetCode = &res.Asset
	source.Effect = history2.Effect{
		Type: history2.EffectTypeCharged,
		Charged: &history2.ChargedEffect{
			Amount: amount.StringU(uint64(op.Amount)),
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
		}
	case history2.EffectTypeFundedToLocked:
		destination.Effect.FundedToLocked = &history2.FundedToLockedEffect{
			Amount: amount.StringU(uint64(op.Amount)),
		}
	}

	return []history2.ParticipantEffect{source, destination}
}
