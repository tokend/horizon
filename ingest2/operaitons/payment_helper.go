package operaitons

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type paymentHelper struct {
	pubKeyProvider publicKeyProvider
}

func (h *paymentHelper) getParticipantsEffects(op xdr.PaymentOpV2,
	res xdr.PaymentV2Response, source history2.ParticipantEffect, destinationEffectType history2.EffectType,
) []history2.ParticipantEffect {
	sourceBalanceID := h.pubKeyProvider.GetBalanceID(op.SourceBalanceId)
	source.BalanceID = &sourceBalanceID
	source.AssetCode = &res.Asset
	source.Effect = history2.Effect{
		Type: history2.EffectTypeCharged,
		Payment: &history2.PaymentEffect{
			Amount: amount.StringU(uint64(op.Amount)),
		},
	}

	destBalanceID := h.pubKeyProvider.GetBalanceID(res.DestinationBalanceId)

	return []history2.ParticipantEffect{source, {
		AccountID: h.pubKeyProvider.GetAccountID(res.Destination),
		BalanceID: &destBalanceID,
		AssetCode: &res.Asset,
		Effect: history2.Effect{
			Type: destinationEffectType,
			Payment: &history2.PaymentEffect{
				Amount: amount.StringU(uint64(op.Amount)),
			},
		},
	}}
}
