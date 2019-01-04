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
