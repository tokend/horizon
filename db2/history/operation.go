package history

import (
	"encoding/json"

	"time"

	"bullioncoin.githost.io/development/go/xdr"
	"gitlab.com/distributed_lab/tokend/horizon/db2"
	"github.com/go-errors/errors"
	"github.com/guregu/null"
)

const (
	PENDING  = 1 + iota
	SUCCESS  = 1 + iota
	REJECTED = 1 + iota
	CANCELED = 1 + iota
	FAILED   = 1 + iota
)

// Operation is a row of data from the `history_operations` table
type Operation struct {
	db2.TotalOrderID
	TransactionID    int64             `db:"transaction_id"`
	TransactionHash  string            `db:"transaction_hash"`
	ApplicationOrder int32             `db:"application_order"`
	Type             xdr.OperationType `db:"type"`
	DetailsString    null.String       `db:"details"`
	LedgerCloseTime  time.Time         `db:"ledger_close_time"`
	SourceAccount    string            `db:"source_account"`
	State            int32             `db:"state"`
	Identifier       int64             `db:"identifier"`
}

// UnmarshalDetails unmarshals the details of this operation into `dest`
func (r *Operation) UnmarshalDetails(dest interface{}) error {
	if !r.DetailsString.Valid {
		return nil
	}

	err := json.Unmarshal([]byte(r.DetailsString.String), &dest)
	if err != nil {
		err = errors.Wrap(err, 1)
	}

	return err
}
