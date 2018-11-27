package operaitons

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type manageAccountOpHandler struct {
	pubKeyProvider publicKeyProvider
}

func (h *manageAccountOpHandler) OperationDetails(op rawOperation,
	_ xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	manageAccountOp := op.Body.MustManageAccountOp()

	return history2.OperationDetails{
		Type: xdr.OperationTypeManageAccount,
		ManageAccount: &history2.ManageAccountDetails{
			AccountID:            manageAccountOp.Account.Address(),
			BlockReasonsToAdd:    int32(manageAccountOp.BlockReasonsToAdd),
			BlockReasonsToRemove: int32(manageAccountOp.BlockReasonsToRemove),
		},
	}, nil
}

func (h *manageAccountOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	_ xdr.OperationResultTr, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	participants := []history2.ParticipantEffect{source}

	participants = append(participants, history2.ParticipantEffect{
		AccountID: h.pubKeyProvider.GetAccountID(opBody.MustManageAccountOp().Account),
	})

	return participants, nil
}
