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
	return []history2.ParticipantEffect{
		h.Participant(sourceID),
		h.Participant(opRes.MustUpdateDataOwnerResult().Success.Owner),
	}, nil
}

func (h *manageUpdateDataOwnerOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	updateDataOwnerOp := op.Body.MustUpdateDataOwnerOp()

	return history2.OperationDetails{
		Type: xdr.OperationTypeUpdateDataOwner,
		UpdateDataOwner: &history2.UpdateDataOwnerDetails{
			ID:       uint64(updateDataOwnerOp.DataId),
			NewOwner: opRes.MustUpdateDataOwnerResult().Success.Owner,
		},
	}, nil
}
