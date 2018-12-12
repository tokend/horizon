package operaitons

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type manageBalanceOpHandler struct {
	pubKeyProvider publicKeyProvider
}

func (h *manageBalanceOpHandler) OperationDetails(opBody xdr.OperationBody, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	manageBalanceOp := opBody.MustManageBalanceOp()

	return history2.OperationDetails{
		Type: xdr.OperationTypeManageBalance,
		ManageBalance: &history2.ManageBalanceDetails{
			DestinationAccount: h.pubKeyProvider.GetAccountID(manageBalanceOp.Destination),
			Action:             manageBalanceOp.Action,
			Asset:              manageBalanceOp.Asset,
			BalanceID:          h.pubKeyProvider.GetBalanceID(opRes.MustManageBalanceResult().MustSuccess().BalanceId),
		},
	}, nil
}

func (h *manageBalanceOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	participants := []history2.ParticipantEffect{source}

	manageBalanceOp := opBody.MustManageBalanceOp()

	destinationAccount := h.pubKeyProvider.GetAccountID(manageBalanceOp.Destination)
	destinationBalance := h.pubKeyProvider.GetBalanceID(opRes.MustManageBalanceResult().MustSuccess().BalanceId)

	if source.AccountID != destinationAccount {
		participants = append(participants, history2.ParticipantEffect{
			AccountID: destinationAccount,
			BalanceID: &destinationBalance,
			AssetCode: &manageBalanceOp.Asset,
		})
	} else {
		participants[0].BalanceID = &destinationBalance
		participants[0].AssetCode = &manageBalanceOp.Asset
	}

	return participants, nil
}
