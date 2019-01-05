package ingest2

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	core "gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/horizon/log"
)

type dbTxManager interface {
	// Begin - starts new db transaction
	Begin() error
	// Rollback - rollbacks db transaction
	Rollback() error
	// Commit - commits db transaction
	Commit() error
}

//Handler - handles ledger and transactions applied for this ledger
type Handler interface {
	// Handle - processes ledger and stores corresponding data to db
	Handle(header *core.LedgerHeader, txs []core.Transaction) error
	// Name - returns name of the Handler
	Name() string
}

// Consumer - consumes ingest data and populate db with it
type Consumer struct {
	log         *log.Entry
	dbTxManager dbTxManager
	dataSource  <-chan LedgerBundle
	handlers    []Handler
}

// NewConsumer - creates new instance of consumer
func NewConsumer(log *log.Entry, dbTxManager dbTxManager, handlers []Handler,
	dataSource <-chan LedgerBundle) *Consumer {
	return &Consumer{
		log:         log.WithField("service", "ingest_consumer"),
		dbTxManager: dbTxManager,
		dataSource:  dataSource,
		handlers:    handlers,
	}
}

//Start - starts consumer in separate goroutine. Must only be used once
func (c *Consumer) Start(ctx context.Context) {
	go c.run(ctx)
}

func (c *Consumer) run(ctx context.Context) {
	for {
		bundles := c.readBatch(ctx)
		if len(bundles) == 0 {
			c.log.Info("Fetched empty ledger bundle batch. It's clear sign that we are going to stop")
			return
		}

		localLog := c.log.WithFields(log.F{
			"batch_len": len(bundles),
			"from":      bundles[0].Header.Sequence,
			"to":        bundles[len(bundles)-1].Header.Sequence,
		})

		localLog.Info("Starting to process new ledger bundles batch")
		running.UntilSuccess(ctx, localLog, "ingest_consumer", func(ctx context.Context) (bool, error) {
			err := c.processBatch(ctx, bundles)
			if err != nil {
				return false, err
			}

			return true, nil
		}, minErrorRecoveryPeriod, maxErrorRecoveryPeriod)
		localLog.Info("Ledger bundles batch processed")
	}

}

func (c *Consumer) readBatch(ctx context.Context) []LedgerBundle {
	const maxBatchSize = 100
	bundles := c.readAtLeastOne(ctx)
	for {
		select {
		case ledgerBundle, ok := <-c.dataSource:
			if !ok {
				return bundles
			}

			bundles = append(bundles, ledgerBundle)
			if len(bundles) >= maxBatchSize {
				return bundles
			}

		case <-ctx.Done():
			return nil
		default:
			return bundles
		}
	}
}

func (c *Consumer) readAtLeastOne(ctx context.Context) []LedgerBundle {
	select {
	case bundle, ok := <-c.dataSource:
		{
			if !ok {
				return nil
			}

			bundles := make([]LedgerBundle, 0, len(c.dataSource)+1)
			bundles = append(bundles, bundle)
			return bundles
		}
	case <-ctx.Done():
		return nil

	}
}

func (c *Consumer) processBatch(ctx context.Context, bundles []LedgerBundle) error {
	err := c.dbTxManager.Begin()
	if err != nil {
		return errors.Wrap(err, "failed to begin db tx")
	}

	defer func() {
		_ = c.dbTxManager.Rollback()
	}()

	for _, bundle := range bundles {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		err = c.processLedgerBundle(ctx, bundle)
		if err != nil {
			return errors.Wrap(err, "failed to process ledger bundle",
				log.F{"ledger_seq": bundle.Header.Sequence})
		}
	}

	err = c.dbTxManager.Commit()
	if err != nil {
		return errors.Wrap(err, "failed to commit db tx")
	}
	return nil
}

func (c *Consumer) processLedgerBundle(ctx context.Context, bundle LedgerBundle) error {
	for _, handler := range c.handlers {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		err := handler.Handle(&bundle.Header, bundle.Transactions)
		if err != nil {
			return errors.Wrap(err, "failed to handle ledger", logan.F{
				"ledger_seq":   bundle.Header.Sequence,
				"handler_name": handler.Name(),
			})
		}
	}

	return nil
}
