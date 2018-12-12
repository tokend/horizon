package reviewrequest

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ingest2/operations"
)

type invoiceHandler struct {
	paymentHelper operations.PaymentHelper
}

func (h *invoiceHandler) specificParticipantsEffects(op xdr.ReviewRequestOp,
	res xdr.ReviewRequestSuccessResult, request xdr.ReviewableRequestEntry,
	source history2.ParticipantEffect, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	paymentOp := op.RequestDetails.MustBillPay().PaymentDetails
	paymentRes := res.TypeExt.MustInvoiceExtended().PaymentV2Response

	effect := history2.EffectTypeFunded
	if request.Body.MustInvoiceRequest().ContractId != nil {
		effect = history2.EffectTypeFundedToLocked
	}

	return h.paymentHelper.GetParticipantsEffects(paymentOp, paymentRes, source, effect)
}
