package core2

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

// TransactionQ - helper struct to get transactions from db
type TransactionQ struct {
	repo *db2.Repo
	selector sq.SelectBuilder
}

// NewTransactionQ - creates new instance of TransactionQ
func NewTransactionQ(repo *db2.Repo) *TransactionQ {
	return &TransactionQ{
		repo: repo,
		selector: sq.Select(
			"tx.txid",
			"tx.ledgerseq",
			"tx.txindex",
			"tx.txbody",
			"tx.txresult",
			"tx.txmeta",
			).From("txhistory tx"),
	}
}

func (q *TransactionQ) GetByLedger(seq int32) ([]Transaction, error) {
	query := q.selector.Where("ledgerseq = ?", seq)
	return q.performSelect(query)
}

func (q *TransactionQ) GetByLedgerRange(fromSeq int32, toSeq int32) ([]Transaction, error) {
	query := q.selector.Where("ledgerseq >= ? AND ledgerseq <= ?", fromSeq, toSeq)
	return q.performSelect(query)
}

func (q *TransactionQ) performSelect(query sq.SelectBuilder) ([]Transaction, error) {
	var result []Transaction
	err := q.repo.Select(&result, query)
	if err != nil {
		if q.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load transactions")
	}

	return result, nil
}
