package operations

import (
	"gitlab.com/tokend/go/xdr"
	core "gitlab.com/tokend/horizon/db2/core2"
)

type operation struct {
	tx  core.Transaction
	opI int
}

//LedgerChanges - returns slice of ledger change for this operation
func (op *operation) LedgerChanges() []xdr.LedgerEntryChange {
	return op.tx.ResultMeta.MustOperations()[op.opI].Changes
}

//Operation - returns xdr operation
func (op *operation) Operation() xdr.Operation {
	return op.tx.Envelope.Tx.Operations[op.opI]
}

//Source - returns source of the operation
func (op *operation) Source() xdr.AccountId {
	opSource := op.Operation().SourceAccount
	if opSource != nil {
		return *opSource
	}

	return op.tx.Envelope.Tx.SourceAccount
}

//Result - returns results of the operation
func (op *operation) Result() xdr.OperationResultTr {
	return op.tx.Result.Result.Result.MustResults()[op.opI].MustTr()
}
