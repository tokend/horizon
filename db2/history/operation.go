package history

import (
	"encoding/json"

	"time"

	"github.com/guregu/null"
	"github.com/pkg/errors"
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

type OperationDetails struct {
	Type          xdr.OperationType
	CreateAccount *CreateAccountDetails
}

func (o *Operation) Details() OperationDetails {
	result := OperationDetails{
		Type: o.Type,
	}

	switch result.Type {
	case xdr.OperationTypeCreateAccount:
		err := json.Unmarshal([]byte(o.DetailsString.String), &result.CreateAccount)
		if err != nil {
			err = errors.Wrap(err, "Error unmarshal operation details")
		}

		return result
	default:
		panic("Invalid operation type")
	}
}

type CreateAccountDetails struct {
	Funder      string `json:"funder,omitempty"`
	Account     string `json:"account,omitempty"`
	AccountType int32  `json:"account_type"`
}

// UnmarshalDetails unmarshals the details of this operation into `dest`
//DEPRECATED
func (r *Operation) UnmarshalDetails(dest interface{}) error {
	if !r.DetailsString.Valid {
		return nil
	}

	err := json.Unmarshal([]byte(r.DetailsString.String), &dest)
	if err != nil {
		err = errors.Wrap(err, "Error unmarshal operation details")
	}

	return err
}
