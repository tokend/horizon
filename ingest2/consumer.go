package ingest2

import (
	"context"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	"time"
)

type dbTxManager interface {
	// Begin - starts new db transaction
	Begin() error
	// Rollback - rollbacks db transaction
	Rollback() error
	// Commit - commits db transaction
	Commit() error
}

type consumer struct {
	log        *logan.Entry
	dbTxManager dbTxManager
	dataSource <-chan ledgerBundle
}

func (c *consumer) Start(ctx context.Context) {
	// normalPeriod is set to 0 to ensure that we are aggressive during catchup. If we failed to get ledger from core
	// runOnce will handle waiting period
	go running.WithBackOff(ctx, c.log, "ingest_consumer", c.runOnce, time.Duration(0), minErrorRecoveryPeriod, maxErrorRecoveryPeriod)
}

func (c *consumer) runOnce(ctx context.Context) error {
	bundles := c.readBatch(ctx)
	if len(bundles) == 0 {
		c.log.Info("Fetched empty ledger bundle batch. It's clear sign that we are going to stop")
		return nil
	}

	c.log.WithFields(logan.F{
		"batch_len": len(bundles),
		"from": bundles[0].Sequence,
		"to": bundles[len(bundles)-1].Sequence,
	}).Info("Starting to process new ledger bundles batch")

	err := c.dbTxManager.Begin()
	if err != nil {
		return errors.Wrap(err, "failed to begin db tx")
	}

	defer c.dbTxManager.Rollback()
	err = c.processBatch(ctx, bundles)
	if err != nil {
		return errors.Wrap(err, "failed to process batch of ledger bundles")
	}

	err = c.dbTxManager.Commit()
	if err != nil {
		return errors.Wrap(err, "failed to commit db tx")
	}
	return nil
}

func (c *consumer) readBatch(ctx context.Context) []ledgerBundle {
	bundles := c.readAtLeastOne(ctx)
	for {
		select {
		case ledgerBundle, ok := <-c.dataSource:
			if !ok {
				return bundles
			}

			bundles = append(bundles, ledgerBundle)
		case <-ctx.Done():
				return nil
		default:
			return bundles
		}
	}
}

func (c *consumer) readAtLeastOne(ctx context.Context) []ledgerBundle{
	select {
	case bundle, ok := <-c.dataSource:
		{
			if !ok {
				return nil
			}

			bundles := make([]ledgerBundle, len(c.dataSource)+1)
			bundles = append(bundles, bundle)
			return bundles
		}
		case <- ctx.Done():
			return nil

	}
}

func (c *consumer) processBatch(ctx context.Context, bundles []ledgerBundle) error {
	for _, bundle := range bundles {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		err := c.processLedgerBundle(ctx, bundle)
		if err != nil {
			return errors.Wrap(err, "failed to process ledger bundle", logan.Field("ledger_seq", bundle.Sequence))
		}
	}

	return nil
}
