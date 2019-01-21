package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/v2"
)

type manageAccountOpHandler struct {
	pubKeyProvider IDProvider
}

// Details returns details about manage account operation
func (h *manageAccountOpHandler) Details(op rawOperation,
	_ xdr.OperationResultTr,
) (regources.OperationDetails, error) {
	manageAccountOp := op.Body.MustManageAccountOp()

	return regources.OperationDetails{
		Type: xdr.OperationTypeManageAccount,
		ManageAccount: &regources.ManageAccountDetails{
			AccountAddress:       manageAccountOp.Account.Address(),
			BlockReasonsToAdd:    xdr.BlockReasons(manageAccountOp.BlockReasonsToAdd),
			BlockReasonsToRemove: xdr.BlockReasons(manageAccountOp.BlockReasonsToRemove),
		},
	}, nil
}

//ParticipantsEffects - returns slice of participants effected by the operation
func (h *manageAccountOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	_ xdr.OperationResultTr, source history2.ParticipantEffect, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	participants := []history2.ParticipantEffect{source}

	participants = append(participants, history2.ParticipantEffect{
		AccountID: h.pubKeyProvider.MustAccountID(opBody.MustManageAccountOp().Account),
	})

	return participants, nil
}
