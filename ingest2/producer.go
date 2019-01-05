package ingest2

import (
	"context"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	core "gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/horizon/ledger"
	"gitlab.com/tokend/horizon/log"
)

const (
	waitForLedgerPeriod    = time.Second
	minErrorRecoveryPeriod = time.Second * 2
	maxErrorRecoveryPeriod = time.Minute * 10
)

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
	txProvider txProvider
	log        *log.Entry

	data          chan LedgerBundle
	currentLedger int32
}

// NewProducer - creates new instance of ingest data produce
func NewProducer(txProvider txProvider, log *log.Entry) *Producer {
	return &Producer{
		log:        log.WithField("service", "ingest_producer"),
		txProvider: txProvider,
	}
}

// Start - starts the Producer in new goroutine. Panics on initialization step if already started.
func (l *Producer) Start(ctx context.Context, bufferSize int, ledgerState ledger.SystemState) chan LedgerBundle {
	if l.data != nil {
		l.log.Panic("Already started")
	}

	l.data = make(chan LedgerBundle, bufferSize)
	l.currentLedger = l.getLedgerSeqToStartFrom(ledgerState)
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

func (l *Producer) getLedgerSeqToStartFrom(ledgerState ledger.SystemState) int32 {
	if ledgerState.History2.Latest == 0 {
		return ledgerState.Core.OldestOnStart
	}

	return ledgerState.History2.Latest + 1
}
