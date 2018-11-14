package ingest2

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
)

type ledgerChangesIteration struct {
	LedgerSeq       int32
	LedgerGlobOpSeq int32
	Op              xdr.Operation
	LedgerChange    xdr.LedgerEntryChange
}

// ledgerChangesConsumerFn - is function which processes ledger changes.
// if error is returned it will be wrapped with corresponding sequences and stop.
// Note: DO NOT take pointers of any of fields of it
type ledgerChangesConsumerFn func(it ledgerChangesIteration) (error)

// ledgerChangesIterator - provides methods to iterate over all ledger changes in one ledger
type ledgerChangesIterator struct {
	bundle *ledgerBundle
}

func newLedgerChangesIterator(bundle *ledgerBundle) ledgerChangesIterator{
	return ledgerChangesIterator{
		bundle: bundle,
	}
}

func (it *ledgerChangesIterator) Iterate(consumerFn ledgerChangesConsumerFn) error {
	ledgerGlobalOpSeq := int32(0)
	for txSeq, tx := range it.bundle.Transactions {

		operationsMeta := *tx.ResultMeta.Operations
		for opSeq, op := range tx.Envelope.Tx.Operations {

			for lcSeq, ledgerChange := range operationsMeta[opSeq].Changes {

				err := consumerFn(ledgerChangesIteration{
					LedgerSeq:       it.bundle.Sequence,
					LedgerGlobOpSeq: ledgerGlobalOpSeq,
					Op:              op,
					LedgerChange:    ledgerChange,
				})
				if err != nil {
					return errors.Wrap(err, "failed to process ledger change", logan.F{
						"ledger_seq":        it.bundle.Sequence,
						"tx_seq":            txSeq,
						"op_seq":            opSeq,
						"ledger_change_seq": lcSeq,
					})
				}
			}

			ledgerGlobalOpSeq++
		}
	}

	return nil
}
