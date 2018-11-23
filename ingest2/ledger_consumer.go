package ingest2

import (
	"context"
	"time"

	"gitlab.com/tokend/go/xdr"

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

type ledgerChange struct {
	changes.LedgerChange
}

func (c *ledgerChange) Consume(ledgerChange) error {
	switch c.LedgerChange.LedgerChange.Type {
	case xdr.LedgerEntryChangeTypeCreated:
		return c.LedgerChange.Created()
	case xdr.LedgerEntryChangeTypeUpdated:
		return c.LedgerChange.Updated()
	case xdr.LedgerEntryChangeTypeRemoved:
		return c.LedgerChange.Deleted()
	case xdr.LedgerEntryChangeTypeState:
		return nil
	default:
		return errors.Errorf("Unrecognized ledger entry change type: %d", c.LedgerChange.LedgerChange.Type)
	}
}

type ledgerChangesConsumer interface {
	Consume(ledgerChange) error
}

type ledgerConsumer struct {
	log           *logan.Entry
	ledgerStorage ledgerStorage
	lcConsumer    ledgerChangesConsumer
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
				err := c.lcConsumer.Consume(ledgerChange{
					changes.LedgerChange{
						LedgerChange:    lc,
						LedgerCloseTime: time.Unix(bundle.Header.CloseTime, 0),
						LedgerSeq:       bundle.Sequence,
						Operation:       &op,
					},
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
