package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type manageBalanceOpHandler struct {
	pubKeyProvider publicKeyProvider
}

// Details returns details about manage balance operation
func (h *manageBalanceOpHandler) Details(op RawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	manageBalanceOp := op.Body.MustManageBalanceOp()
	manageBalanceRes := opRes.MustManageBalanceResult().MustSuccess()

	return history2.OperationDetails{
		Type: xdr.OperationTypeManageBalance,
		ManageBalance: &history2.ManageBalanceDetails{
			DestinationAccount: manageBalanceOp.Destination.Address(),
			Action:             manageBalanceOp.Action,
			Asset:              manageBalanceOp.Asset,
			BalanceAddress:     manageBalanceRes.BalanceId.AsString(),
		},
	}, nil
}

func (h *manageBalanceOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	manageBalanceOp := opBody.MustManageBalanceOp()

	destinationAccount := h.pubKeyProvider.GetAccountID(manageBalanceOp.Destination)

	var participants []history2.ParticipantEffect

	if source.AccountID != destinationAccount {
		participants = []history2.ParticipantEffect{{
			AccountID: destinationAccount,
		}}
	}

	return append(participants, source), nil
}
