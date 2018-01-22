package history

import (
	"encoding/json"

	"time"

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
	Details          OperationDetails  `db:"details"`
	LedgerCloseTime  time.Time         `db:"ledger_close_time"`
	SourceAccount    string            `db:"source_account"`
	State            OperationState    `db:"state"`
	Identifier       int64             `db:"identifier"`
}

func (o *Operation) GetDetails() OperationDetails {
	result := OperationDetails{
		Type: o.Type,
	}

	details, err := json.Marshal(o.Details)
	if err != nil {
		logrus.WithError(err).Errorf("Error marshal operation details")
	}

	if err := json.Unmarshal(details, &result); err != nil {
		logrus.WithError(err).Errorf("Error unmarshal operation details")
	}
	return result
}
