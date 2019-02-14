package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type manageAccountOpHandler struct {
	effectsProvider
}

// Details returns details about manage account operation
func (h *manageAccountOpHandler) Details(op rawOperation,
	_ xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	manageAccountOp := op.Body.MustManageAccountOp()

	return history2.OperationDetails{
		Type: xdr.OperationTypeManageAccount,
		ManageAccount: &history2.ManageAccountDetails{
			AccountAddress:       manageAccountOp.Account.Address(),
			BlockReasonsToAdd:    xdr.BlockReasons(manageAccountOp.BlockReasonsToAdd),
			BlockReasonsToRemove: xdr.BlockReasons(manageAccountOp.BlockReasonsToRemove),
		},
	}, nil
}

//ParticipantsEffects - returns slice of participants effected by the operation
func (h *manageAccountOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	_ xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	source := h.effectsProvider.Participant(sourceAccountID)
	return []history2.ParticipantEffect{source, h.effectsProvider.Participant(opBody.MustManageAccountOp().Account)}, nil
}
