package operaitons

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type paymentHelper struct {
}

func (h *paymentHelper) getParticipantsEffects(op xdr.PaymentOpV2,
	res xdr.PaymentV2Response, source history2.ParticipantEffect,
) []history2.ParticipantEffect {

}
