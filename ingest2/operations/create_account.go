package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/v2"
)

type createAccountOpHandler struct {
	pubKeyProvider IDProvider
}

// Details returns details about create account operation
func (h *createAccountOpHandler) Details(op rawOperation, _ xdr.OperationResultTr,
) (regources.OperationDetails, error) {
	createAccOp := op.Body.MustCreateAccountOp()

	return regources.OperationDetails{
		Type: xdr.OperationTypeCreateAccount,
		CreateAccount: &regources.CreateAccountDetails{
			AccountAddress: createAccOp.Destination.Address(),
			AccountType:    createAccOp.AccountType,
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
		AccountID: h.pubKeyProvider.MustAccountID(createAccountOp.Destination),
	})

	if createAccountOp.Referrer != nil {
		participants = append(participants, history2.ParticipantEffect{
			AccountID: h.pubKeyProvider.MustAccountID(*createAccountOp.Referrer),
		})
	}

	return participants, nil
}
