package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/rgenerated"
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
			BaseAsset:               string(manageAssetPairOp.Base),
			QuoteAsset:              string(manageAssetPairOp.Quote),
			PhysicalPrice:           rgenerated.Amount(manageAssetPairOp.PhysicalPrice),
			PhysicalPriceCorrection: rgenerated.Amount(manageAssetPairOp.PhysicalPriceCorrection),
			MaxPriceStep:            rgenerated.Amount(manageAssetPairOp.MaxPriceStep),
			Policies:                xdr.AssetPairPolicy(manageAssetPairOp.Policies),
		},
	}, nil
}

//ParticipantsEffects - returns source of the operation
func (h *manageAssetPairOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	_ xdr.OperationResultTr, sourceAccountID xdr.AccountId, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	result := h.manageOfferOpHandler.getDeletedOffersEffect(ledgerChanges)
	return append(result, h.manageOfferOpHandler.Participant(sourceAccountID)), nil
}
