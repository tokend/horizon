package ingest2

import (
	"context"
	"time"

	"gitlab.com/tokend/horizon/ingest2/changes"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
	"gitlab.com/tokend/horizon/db2/history2"
)

type ledgerStorage interface {
	InsertLedger(ledger history2.Ledger) error
}

type operation struct {
}

type operationConsumer interface {
	Consume(operation) error
}

type ChangesConsumer interface {
	Consume(changes.LedgerChange) error
}

type ledgerConsumer struct {
	log           *logan.Entry
	ledgerStorage ledgerStorage
	lcConsumer    ChangesConsumer
	opConsumer    operationConsumer
}

func (c *ledgerConsumer) Consume(ctx context.Context, bundle ledgerBundle) error {
	ledgerGlobalOpSeq := int32(0)
	fields := logan.F{
		"ledger_seq": bundle.Sequence,
	}
	for txSeq, tx := range bundle.Transactions {
		fields = fields.Add("tx_seq", txSeq)
		operationsMeta := *tx.ResultMeta.Operations
		for opSeq, op := range tx.Envelope.Tx.Operations {
			fields = fields.Add("opSeq", opSeq)

			for lcSeq, lc := range operationsMeta[opSeq].Changes {
				fields = fields.Add("ledger_change_seq", lcSeq)
				err := c.lcConsumer.Consume(
					changes.LedgerChange{
						LedgerChange:    lc,
						LedgerCloseTime: time.Unix(bundle.Header.CloseTime, 0).UTC(),
						LedgerSeq:       bundle.Sequence,
						Operation:       &op,
					})
				if err != nil {
					return errors.Wrap(err, "failed to process ledger change", fields)
				}
			}

			err := c.opConsumer.Consume(operation{})
			if err != nil {
				return errors.Wrap(err, "failed to process operation", fields)
			}

			ledgerGlobalOpSeq++
		}
	}

	err := c.ledgerStorage.InsertLedger(history2.Ledger{
		TotalOrderID: db2.TotalOrderID{
			ID: int64(bundle.Sequence),
		},
		Hash:         bundle.Header.LedgerHash,
		PreviousHash: bundle.Header.PrevHash,
		ClosedAt:     time.Unix(bundle.Header.CloseTime, 0).UTC(),
	})
	if err != nil {
		return errors.Wrap(err, "failed to insert ledger", fields)
	}

	return nil
}
