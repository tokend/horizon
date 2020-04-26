package storage

import (
	"github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/bridge"
	"gitlab.com/tokend/horizon/db2/history2"
)

//Ledger - handles write operations on db level for ledgers
type Ledger struct {
	repo *bridge.Mediator
}

//NewLedger - creates new instance of ledger
func NewLedger(repo *bridge.Mediator) *Ledger {
	return &Ledger{
		repo: repo,
	}
}

//Insert - inserts Ledger into DB
func (s *Ledger) Insert(ledger *history2.Ledger) error {
	sql := squirrel.Insert("ledgers").Columns("id", "sequence", "hash", "previous_hash", "closed_at",
		"tx_count", "data").Values(ledger.ID, ledger.Sequence, ledger.Hash, ledger.PreviousHash, ledger.ClosedAt,
		ledger.TxCount, ledger.Data)

	err := s.repo.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to insert ledger", logan.F{
			"ledger_seq": ledger.Sequence,
		})
	}

	return nil
}
