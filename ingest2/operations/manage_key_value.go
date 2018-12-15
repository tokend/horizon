package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type manageKeyValueOpHandler struct {
}

// Details returns details about manage key value operation
func (h *manageKeyValueOpHandler) Details(op RawOperation, _ xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	manageKVOp := op.Body.MustManageKeyValueOp()

	var value *xdr.KeyValueEntryValue
	if manageKVOp.Action.Action == xdr.ManageKvActionPut {
		valueForPtr := manageKVOp.Action.MustValue()
		value = &valueForPtr
	}

	return history2.OperationDetails{
		Type: xdr.OperationTypeManageKeyValue,
		ManageKeyValue: &history2.ManageKeyValueDetails{
			Key:    string(manageKVOp.Key),
			Action: manageKVOp.Action.Action,
			Value:  value,
		},
	}, nil
}

func (h *manageKeyValueOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	_ xdr.OperationResultTr, source history2.ParticipantEffect, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{source}, nil
}
