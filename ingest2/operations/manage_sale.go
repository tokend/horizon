package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	regources "gitlab.com/tokend/regources/v2"
)

type manageSaleHandler struct {
	manageOfferOpHandler *manageOfferOpHandler
}

// Details returns details about payout operation
func (h *manageSaleHandler) Details(op rawOperation, res xdr.OperationResultTr,
) (regources.OperationDetails, error) {
	manageSale := op.Body.MustManageSaleOp()

	return regources.OperationDetails{
		Type: xdr.OperationTypeManageSale,
		ManageSale: &regources.ManageSaleDetails{
			SaleID: uint64(manageSale.SaleId),
			Action: manageSale.Data.Action,
		},
	}, nil
}

// ParticipantsEffects returns `charged` and `funded` effects
func (h *manageSaleHandler) ParticipantsEffects(opBody xdr.OperationBody,
	res xdr.OperationResultTr, source history2.ParticipantEffect, changes []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	manageSale := opBody.MustManageSaleOp().Data
	if manageSale.Action != xdr.ManageSaleActionCancel {
		return []history2.ParticipantEffect{source}, nil
	}

	return h.manageOfferOpHandler.getDeletedOffersEffect(changes), nil
}
