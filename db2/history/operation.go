package history

import (
	"encoding/json"

	"time"

	"github.com/guregu/null"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon/db2"
)

// Operation is a row of data from the `history_operations` table
type Operation struct {
	db2.TotalOrderID
	TransactionID    int64             `db:"transaction_id"`
	TransactionHash  string            `db:"transaction_hash"`
	ApplicationOrder int32             `db:"application_order"`
	Type             xdr.OperationType `db:"type"`
	// DEPRECATED
	DetailsString   null.String    `db:"details"`
	LedgerCloseTime time.Time      `db:"ledger_close_time"`
	SourceAccount   string         `db:"source_account"`
	State           OperationState `db:"state"`
	Identifier      int64          `db:"identifier"`
}

func (o *Operation) Details() OperationDetails {
	result := OperationDetails{
		Type: o.Type,
	}

	if err := json.Unmarshal([]byte(o.DetailsString.String), &result); err != nil {
		logrus.WithError(err).Errorf("Error unmarshal operation details")
	}
	return result
}

// UnmarshalDetails unmarshals the details of this operation into `dest`
//DEPRECATED
func (r *Operation) UnmarshalDetails(dest interface{}) error {
	if !r.DetailsString.Valid {
		return nil
	}
	err := json.Unmarshal([]byte(r.DetailsString.String), &dest)
	if err != nil {
		return errors.Wrap(err, "Error unmarshal operation details")
	}

	return nil
}
