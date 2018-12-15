package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type createAccountOpHandler struct {
	pubKeyProvider publicKeyProvider
}

// Details returns details about create account operation
func (h *createAccountOpHandler) Details(op RawOperation, _ xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	operation := op.Body.MustCreateAccountOp()

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreateAccount,
		CreateAccount: &history2.CreateAccountDetails{
			AccountAddress: operation.Destination.Address(),
			AccountType:    operation.AccountType,
		},
	}, nil
}

// ParticipantsEffects returns counterparties without effects
func (h *createAccountOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	_ xdr.OperationResultTr, source history2.ParticipantEffect, _ []xdr.LedgerEntryChange,
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
