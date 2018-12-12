package operaitons

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type createAccountOpHandler struct {
	pubKeyProvider publicKeyProvider
}

func (h *createAccountOpHandler) OperationDetails(opBody xdr.OperationBody, _ xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	op := opBody.MustCreateAccountOp()

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreateAccount,
		CreateAccount: &history2.CreateAccountDetails{
			AccountID:   h.pubKeyProvider.GetAccountID(op.Destination),
			AccountType: op.AccountType,
		},
	}, nil
}

func (h *createAccountOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	_ xdr.OperationResultTr, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	participants := []history2.ParticipantEffect{source}

	createAccountOp := opBody.MustCreateAccountOp()

	participants = append(participants, history2.ParticipantEffect{
		AccountID: h.pubKeyProvider.GetAccountID(createAccountOp.Destination),
	})

	if createAccountOp.Referrer != nil {
		participants = append(participants, history2.ParticipantEffect{
			AccountID: h.pubKeyProvider.GetAccountID(*createAccountOp.Referrer),
		})
	}

	return participants, nil
}
