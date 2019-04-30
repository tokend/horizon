package ingest2

import (
	"context"
	"math"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	core "gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ledger"
)

const (
	maxBatchSize = 1000
	waitForLedgerPeriod    = time.Second
	minErrorRecoveryPeriod = time.Second * 2
	maxErrorRecoveryPeriod = time.Minute * 10
)

// historyLedgerProvider - specifies methods required to get ledger from history db
type historyLedgerProvider interface {
	GetBySequence(seq int32) (*history2.Ledger, error)
}

// txProvider - specifies methods required to get data from core db needed for ingest
type txProvider interface {
	GetLatestLedgerSeq() (int32, error)
	GetByLedgerRange(fromSeq int32, toSeq int32) ([]core.Transaction, error)
	GetBySequence(seq int32) (*core.LedgerHeader, error)
	GetBySequenceRange(fromSeq int32, toSeq int32) ([]core.LedgerHeader, error)
}

// Producer - worker which is responsible for loading sequence of ledgers and transactions and sending them
// to provided channel.
type Producer struct {
	txProvider      txProvider
	hLedgerProvider historyLedgerProvider
	log             *logan.Entry

	data          chan LedgerBundle
	currentLedger int32
}

// NewProducer - creates new instance of ingest data produce
func NewProducer(txProvider txProvider, hLedgerProvider historyLedgerProvider, log *logan.Entry) *Producer {
	return &Producer{
		log:             log,
		txProvider:      txProvider,
		hLedgerProvider: hLedgerProvider,
	}
}

// Start - starts the Producer in new goroutine. Panics on initialization step if already started.
func (l *Producer) Start(ctx context.Context, bufferSize int, ledgerState ledger.SystemState) chan LedgerBundle {
	if l.data != nil {
		l.log.Panic("Already started")
	}

	l.data = make(chan LedgerBundle, bufferSize)
	var err error
	l.currentLedger, err = l.getLedgerSeqToStartFrom(ledgerState)
	if err != nil {
		l.log.WithError(err).Panic("Failed to figure out ledger sequence from which to start")
	}

	err = l.ensureChainIsConsistent(l.currentLedger)
	if err != nil {
		l.log.WithError(err).Panic("Failed to ensure consistency of core and horizon chains")
	}
	go running.WithBackOff(ctx, l.log, "ingest_producer", l.checkForNewLedgers, waitForLedgerPeriod,
		minErrorRecoveryPeriod, maxErrorRecoveryPeriod)

	return l.data
}

func (l *Producer) checkForNewLedgers(ctx context.Context) error {
	latestLedger, err := l.txProvider.GetLatestLedgerSeq()
	if err != nil {
		return errors.Wrap(err, "failed to load current ledger sequence")
	}

	if latestLedger > l.currentLedger {
		fromSeq := l.currentLedger + 1
		toSeq := int32(math.Min(float64(latestLedger), float64(l.currentLedger + maxBatchSize)))
		batch, err := l.getBatch(fromSeq, toSeq)
		if err != nil {
			return errors.Wrap(err, "failed to load batch")
		}

		for _, bundle := range batch {
			l.data <- bundle
		}
		l.currentLedger = toSeq
	}

	return nil
}

func (l *Producer) getBatch(fromSeq int32, toSeq int32) ([]LedgerBundle, error) {
	headers, err := l.txProvider.GetBySequenceRange(fromSeq, toSeq)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load ledger headers")
	}
	txs, err := l.txProvider.GetByLedgerRange(fromSeq, toSeq)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load ledger transactions")
	}
	result := make([]LedgerBundle, 0, len(headers))

	for _, header := range headers {
		result = append(result, LedgerBundle{
			Header: header,
			Transactions: make([]core.Transaction, 0, header.Data.MaxTxSetSize),
		})
	}

	for _, tx := range txs {
		if tx.IsSuccessful() {
			bundle := result[tx.LedgerSequence - fromSeq]
			bundle.Transactions = append(bundle.Transactions, tx)
			result[tx.LedgerSequence - fromSeq] = bundle
		}
	}

	return result, nil
}

// ensureChainIsConsistent - ensures that we'll end up with consistent chain of blocks in both core and horizon.
func (l *Producer) ensureChainIsConsistent(ledgerSeqToIngest int32) error {
	// try to load previously ingested ledger
	prevLedger, err := l.hLedgerProvider.GetBySequence(ledgerSeqToIngest - 1)
	if err != nil {
		return errors.Wrap(err, "failed to load ledger by sequence", logan.F{
			"ledger_seq": ledgerSeqToIngest - 1,
		})
	}

	// horizon db is empty, so we have the same chain
	if prevLedger == nil {
		var coreLedger *core.LedgerHeader
		coreLedger, err = l.txProvider.GetBySequence(2)
		if err != nil {
			return errors.Wrap(err, "failed to load ledger from core DB")
		}

		if coreLedger == nil {
			return errors.Wrap(err, "Core does not have full history. Unfortunately this version of horizon does"+
				" not support partial history")
		}
		return nil
	}

	currentLedger, err := l.txProvider.GetBySequence(ledgerSeqToIngest)
	if err != nil {
		return errors.Wrap(err, "failed to load ledger by seq from core", logan.F{
			"ledger_seq": ledgerSeqToIngest,
		})
	}

	if currentLedger == nil {
		return errors.From(errors.New("failed to load from core ledger to ingest - horizon is ahead of core or core is dead"), logan.F{
			"ledger_to_ingest": ledgerSeqToIngest,
		})
	}

	if currentLedger.PrevHash != prevLedger.Hash {
		return errors.From(errors.New("chain in horizon does not match chain in core"), logan.F{
			"current_ledger_prev_hash": currentLedger.PrevHash,
			"prev_ledger_hash":         prevLedger.Hash,
			"prev_ledger_seq":          prevLedger.Sequence,
			"current_ledger_seq":       currentLedger.Sequence,
		})
	}

	return nil

}

func (l *Producer) getLedgerSeqToStartFrom(ledgerState ledger.SystemState) (int32, error) {
	if ledgerState.Core.Latest < ledgerState.History2.Latest {
		return 0, errors.From(errors.New("horizon is ahead of core"), logan.F{
			"core_latest":     ledgerState.Core.Latest,
			"history2_latest": ledgerState.History2.Latest,
		})
	}
	if ledgerState.History2.Latest == 0 {
		return ledgerState.Core.OldestOnStart, nil
	}

	return ledgerState.History2.Latest + 1, nil
}
