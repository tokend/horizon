package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/v2"
)

type manageBalanceOpHandler struct {
	pubKeyProvider IDProvider
}

// Details returns details about manage balance operation
func (h *manageBalanceOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (regources.OperationDetails, error) {
	manageBalanceOp := op.Body.MustManageBalanceOp()
	manageBalanceRes := opRes.MustManageBalanceResult().MustSuccess()

	return regources.OperationDetails{
		Type: xdr.OperationTypeManageBalance,
		ManageBalance: &regources.ManageBalanceDetails{
			DestinationAccount: manageBalanceOp.Destination.Address(),
			Action:             manageBalanceOp.Action,
			Asset:              string(manageBalanceOp.Asset),
			BalanceAddress:     manageBalanceRes.BalanceId.AsString(),
		},
	}, nil
}

//ParticipantsEffects - returns source of the operation and account for which balance was created if they differ
func (h *manageBalanceOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	manageBalanceOp := opBody.MustManageBalanceOp()

	destinationAccount := h.pubKeyProvider.MustAccountID(manageBalanceOp.Destination)

	var participants []history2.ParticipantEffect

	if source.AccountID != destinationAccount {
		participants = []history2.ParticipantEffect{{
			AccountID: destinationAccount,
		}}
	}

	return append(participants, source), nil
}
