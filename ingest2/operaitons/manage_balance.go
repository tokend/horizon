package operaitons

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type manageBalanceOpHandler struct {
	pubKeyProvider publicKeyProvider
}

func (h *manageBalanceOpHandler) OperationDetails(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	manageBalanceOp := op.Body.MustManageBalanceOp()

	return history2.OperationDetails{
		Type: xdr.OperationTypeManageBalance,
		ManageBalance: &history2.ManageBalanceDetails{
			DestinationAccount: manageBalanceOp.Destination.Address(),
			Action:             manageBalanceOp.Action,
			Asset:              manageBalanceOp.Asset,
			BalanceID:          opRes.MustManageBalanceResult().MustSuccess().BalanceId.AsString(),
		},
	}, nil
}

func (h *manageBalanceOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	manageBalanceOp := opBody.MustManageBalanceOp()

	destinationAccount := h.pubKeyProvider.GetAccountID(manageBalanceOp.Destination)
	destinationBalance := h.pubKeyProvider.GetBalanceID(opRes.MustManageBalanceResult().MustSuccess().BalanceId)

	var participants []history2.ParticipantEffect

	if source.AccountID != destinationAccount {
		participants = append(participants, history2.ParticipantEffect{
			AccountID: destinationAccount,
			BalanceID: &destinationBalance,
			AssetCode: &manageBalanceOp.Asset,
		})
	} else {
		source.BalanceID = &destinationBalance
		source.AssetCode = &manageBalanceOp.Asset
	}

	return append(participants, source), nil
}
