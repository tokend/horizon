package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type manageUpdateDataOwnerOpHandler struct {
	effectsProvider
}

func (h *manageUpdateDataOwnerOpHandler) ParticipantsEffects(opBody xdr.OperationBody, opRes xdr.OperationResultTr,
	sourceID xdr.AccountId, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{h.Participant(sourceID)}, nil
}

func (h *manageUpdateDataOwnerOpHandler) Details(op rawOperation, _ xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	updateDataOwnerOp := op.Body.MustUpdateDataOwnerOp()

	return history2.OperationDetails{
		Type: xdr.OperationTypeUpdateData,
		UpdateDataOwner: &history2.UpdateDataOwnerDetails{
			NewOwner: updateDataOwnerOp.NewOwner,
			ID:       uint64(updateDataOwnerOp.DataId),
		},
	}, nil
}
