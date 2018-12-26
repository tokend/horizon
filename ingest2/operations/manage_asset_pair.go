package operations

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type manageAssetPairOpHandler struct {
	manageOfferOpHandler *manageOfferOpHandler
}

// Details returns details about manage asset pair operation
func (h *manageAssetPairOpHandler) Details(op rawOperation,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	manageAssetPairOp := op.Body.MustManageAssetPairOp()

	return history2.OperationDetails{
		Type: xdr.OperationTypeManageAssetPair,
		ManageAssetPair: &history2.ManageAssetPairDetails{
			BaseAsset:               manageAssetPairOp.Base,
			QuoteAsset:              manageAssetPairOp.Quote,
			PhysicalPrice:           amount.String(int64(manageAssetPairOp.PhysicalPrice)),
			PhysicalPriceCorrection: amount.String(int64(manageAssetPairOp.PhysicalPriceCorrection)),
			MaxPriceStep:            amount.String(int64(manageAssetPairOp.MaxPriceStep)),
			PoliciesI:               int32(manageAssetPairOp.Policies),
		},
	}, nil
}

func (h *manageAssetPairOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	_ xdr.OperationResultTr, source history2.ParticipantEffect, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	result := h.manageOfferOpHandler.getDeletedOffersEffect(ledgerChanges)
	return append(result, source), nil
}
