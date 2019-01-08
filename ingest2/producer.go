package ingest2

import (
	"context"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	core "gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ledger"
	"gitlab.com/tokend/horizon/log"
)

const (
	waitForLedgerPeriod    = time.Second
	minErrorRecoveryPeriod = time.Second * 2
	maxErrorRecoveryPeriod = time.Minute * 10
)

// historyLedgerProvider - specifies methods required to get ledger from history db
type historyLedgerProvider interface {
	BySequence(seq int32) (*history2.Ledger, error)
}

// txProvider - specifies methods required to get data from core db needed for ingest
type txProvider interface {
	// TransactionsByLedger returns slice of transaction for given ledger sequence. Returns empty slice,
	// nil if there is no transactions
	TransactionsByLedger(seq int32) ([]core.Transaction, error)
	// LedgerHeaderBySequence returns *core.LedgerHeader by its sequence. Returns nil, nil if ledgerHeader
	// does not exists
	LedgerHeaderBySequence(seq int32) (*core.LedgerHeader, error)
}

// Producer - worker which is responsible for loading sequence of ledgers and transactions and sending them
// to provided channel.
type Producer struct {
	txProvider      txProvider
	hLedgerProvider historyLedgerProvider
	log             *log.Entry

	data          chan LedgerBundle
	currentLedger int32
}

// NewProducer - creates new instance of ingest data produce
func NewProducer(txProvider txProvider, hLedgerProvider historyLedgerProvider, log *log.Entry) *Producer {
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
	// normalPeriod is set to 0 to ensure that we are aggressive during catchup. If we failed to get ledger from core
	// runOnce will handle waiting period
	go running.WithBackOff(ctx, l.log, "ingest_producer", l.runOnce, time.Duration(0),
		minErrorRecoveryPeriod, maxErrorRecoveryPeriod)
	return l.data
}

func (l *Producer) runOnce(ctx context.Context) error {
	isSent, err := l.trySendBlock(ctx, l.currentLedger)
	if err != nil {
		return errors.Wrap(err, "failed to send block", logan.F{"seq": l.currentLedger})
	}

	if isSent {
		l.currentLedger++
		return nil
	}

	// ledger is not available yet in the core, so we can sleep for a sec
	time.Sleep(waitForLedgerPeriod)
	return nil
}

func (l *Producer) trySendBlock(ctx context.Context, seq int32) (bool, error) {
	ledgerHeader, err := l.txProvider.LedgerHeaderBySequence(seq)
	if err != nil {
		return false, errors.Wrap(err, "failed to load ledger header", logan.F{"ledger_seq": seq})
	}

	// ledger still does not exists
	if ledgerHeader == nil {
		return false, nil
	}

	txs, err := l.loadSuccessTxs(seq)
	if err != nil {
		return false, errors.Wrap(err, "failed to load successful transactions for ledger",
			logan.F{"ledger_seq": seq})
	}

	select {
	case <-ctx.Done():
		return false, nil
	case l.data <- newLedgerBundle(*ledgerHeader, txs):
		return true, nil
	}
}

func (l *Producer) loadSuccessTxs(seq int32) ([]core.Transaction, error) {
	txs, err := l.txProvider.TransactionsByLedger(seq)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load transactions for ledger",
			logan.F{"ledger_seq": seq})
	}

	successTxs := make([]core.Transaction, 0, len(txs))
	for i := range txs {
		if !txs[i].IsSuccessful() {
			continue
		}

		successTxs = append(successTxs, txs[i])
	}

	return successTxs, nil
}

// ensureChainIsConsistent - ensures that we'll end up with consistent chain of blocks in both core and horizon.
func (l *Producer) ensureChainIsConsistent(ledgerSeqToIngest int32) error {
	// try to load previously ingested ledger
	prevLedger, err := l.hLedgerProvider.BySequence(ledgerSeqToIngest - 1)
	if err != nil {
		return errors.Wrap(err, "failed to load ledger by sequence", logan.F{
			"ledger_seq": ledgerSeqToIngest - 1,
		})
	}

	// horizon db is empty, so we have the same chain
	if prevLedger == nil {
		return nil
	}

	currentLedger, err := l.txProvider.LedgerHeaderBySequence(ledgerSeqToIngest)
	if err != nil {
		return errors.Wrap(err, "failed to load ledger by seq from core", logan.F{
			"ledger_seq": ledgerSeqToIngest,
		})
	}

	if currentLedger == nil {
		return errors.From(errors.New("failed to load from core ledger to ingest - horizon is ahead of core"), logan.F{
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
