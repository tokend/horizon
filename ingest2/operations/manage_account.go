package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type manageAccountOpHandler struct {
	pubKeyProvider publicKeyProvider
}

// Details returns details about manage account operation
func (h *manageAccountOpHandler) Details(op RawOperation,
	_ xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	manageAccountOp := op.Body.MustManageAccountOp()

	return history2.OperationDetails{
		Type: xdr.OperationTypeManageAccount,
		ManageAccount: &history2.ManageAccountDetails{
			AccountAddress:       manageAccountOp.Account.Address(),
			BlockReasonsToAdd:    int32(manageAccountOp.BlockReasonsToAdd),
			BlockReasonsToRemove: int32(manageAccountOp.BlockReasonsToRemove),
		},
	}, nil
}

func (h *manageAccountOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	_ xdr.OperationResultTr, source history2.ParticipantEffect, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	participants := []history2.ParticipantEffect{source}

	participants = append(participants, history2.ParticipantEffect{
		AccountID: h.pubKeyProvider.GetAccountID(opBody.MustManageAccountOp().Account),
	})

	return participants, nil
}
