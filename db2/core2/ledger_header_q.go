package core2

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

// LedgerHeaderQ - helper struct to get ledger headers from db
type LedgerHeaderQ struct {
	repo *db2.Repo
}

// NewLedgerHeaderQ - creates new instance of LedgerHeaderQ
func NewLedgerHeaderQ(repo *db2.Repo) *LedgerHeaderQ {
	return &LedgerHeaderQ{
		repo: repo,
	}
}

// LedgerHeaderBySequence returns *core.LedgerHeader by its sequence. Returns nil, nil if ledgerHeader does not exists
func (q *LedgerHeaderQ) LedgerHeaderBySequence(seq int32) (*LedgerHeader, error) {
	query := sq.Select("l.ledgerhash, l.prevhash, l.bucketlisthash, l.closetime, l.ledgerseq, l.version, l.data").
		From("ledgerheaders l").Where("ledgerseq = ?", seq)
	var header LedgerHeader
	err := q.repo.Get(&header, query)
	if err != nil {
		if q.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load ledger by sequence", logan.F{"seq": seq})
	}

	return &header, nil
}

// LatestLedgerSeq - returns latest ledger sequence available in DB
func (q *LedgerHeaderQ) LatestLedgerSeq() (int32, error) {
	var result int32
	err := q.repo.GetRaw(&result, "SELECT COALESCE(MAX(ledgerseq), 1) FROM ledgerheaders")
	if err != nil {
		return 0, errors.Wrap(err, "failed to get latest ledger seq")
	}

	return result, nil
}

// OldestLedgerSeq - returns oldest ledger sequence (which does not have gaps) available in the core db.
// Due to design of core it will always have 1 ledger in the db. So we try to find min ledger sequence > 1
// In case it's 2 we will return 1.
// Note: this method does not handle and should not handle existence of 2 gaps
func (q *LedgerHeaderQ) OldestLedgerSeq() (int32, error) {
	var result int32
	err := q.repo.GetRaw(&result, "SELECT COALESCE(MIN(ledgerseq), 1) FROM ledgerheaders WHERE ledgerseq > 1")
	if err != nil {
		return 0, errors.Wrap(err, "failed to get oldest ledger seq")
	}

	if result == 2 {
		return 1, nil
	}

	return result, nil
}
