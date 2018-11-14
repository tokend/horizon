package ingest2

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// accountsProcessorI - consumes ledger changes and generates for each account ID unique int64 identifier
type accountsProcessorI interface {
	// panics if not able to provide ID for specified accountID
	GetID(accountID xdr.AccountId) int64
	// Stores all new accounts into persistent storage
	Store() error
	// Consume - processes ledger changes
	Consume(it ledgerChangesIteration) (error)

}

type consumer struct {
	log *logan.Entry

	accountsProcessor *accountsProcessor
}

func (c *consumer) processBatch(db *db2.Repo, bundles []ledgerBundle) error {
	storage, err := newStorage(db)
	if err != nil {
		return errors.Wrap(err, "failed to init storage")
	}

	storage.Rollback()

	for i := range bundles {
		ledgerChanges := newLedgerChangesIterator(&bundles[i])
		err := ledgerChanges.Iterate(c.accountsProcessor.Consume)
		if err != nil {
			return errors.Wrap(err, "failed to process new accounts")
		}

	}

	c.accountsProcessor.Store(storage)
}




