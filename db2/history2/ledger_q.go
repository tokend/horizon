package history2

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

// LedgerQ - is a helper struct to construct requests to ledger table
type LedgerQ struct {
	repo *db2.Repo
}

// NewLedgerQ - creates new instance of LedgerQ
func NewLedgerQ(repo *db2.Repo) *LedgerQ {
	return &LedgerQ{
		repo: repo,
	}
}

// LatestLedgerSeq - returns latest ledger sequence available in DB
func (q *LedgerQ) LatestLedgerSeq() (int32, error) {
	var result int32
	err := q.repo.GetRaw(&result, "SELECT COALESCE(MAX(sequence), 0) FROM ledgers")
	if err != nil {
		return 0, errors.Wrap(err, "failed to get latest ledger seq")
	}

	return result, nil
}

// OldestLedgerSeq - returns oldest ledger sequence
func (q *LedgerQ) OldestLedgerSeq() (int32, error) {
	var result int32
	err := q.repo.GetRaw(&result, "SELECT COALESCE(MIN(sequence), 0) FROM ledgers")
	if err != nil {
		return 0, errors.Wrap(err, "failed to get oldest ledger seq")
	}

	return result, nil
}

//BySequence - returns ledger, if ledger with specified seq does not exists - returns nil, nil
func (q *LedgerQ) BySequence(seq int32) (*Ledger, error) {
	var result Ledger
	err := q.repo.Get(&result, sq.Select("l.id, l.sequence, l.hash, l.previous_hash", "l.closed_at", "l.tx_count", "l.data").
		From("ledgers l").Where("l.id = ?", seq))
	if err != nil {
		if q.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load ledger by sequence", logan.F{"sequence": seq})
	}

	return &result, nil
}
