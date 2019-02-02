package core2

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

// TransactionQ - helper struct to get transactions from db
type TransactionQ struct {
	repo *db2.Repo
}

// NewTransactionQ - creates new instance of TransactionQ
func NewTransactionQ(repo *db2.Repo) *TransactionQ {
	return &TransactionQ{
		repo: repo,
	}
}

// GetByLedger returns slice of transaction for given ledger sequence. Returns empty slice, nil if there is no transactions
func (q *TransactionQ) GetByLedger(seq int32) ([]Transaction, error) {
	query := sq.Select("tx.txid, tx.ledgerseq, tx.txindex, tx.txbody, tx.txresult, tx.txmeta").
		From("txhistory tx").Where("ledgerseq = ?", seq)
	var result []Transaction
	err := q.repo.Select(&result, query)
	if err != nil {
		if q.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load transactions for ledger", logan.F{
			"ledger_seq": seq,
		})
	}

	return result, nil
}
