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
	Payment       *PaymentDetails
	SetOptions    *SetOptionsDetails
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
	case xdr.OperationTypePayment:
		err := json.Unmarshal([]byte(o.DetailsString.String), &result.Payment)
		if err != nil {
			err = errors.Wrap(err, "Error unmarshal operation details")
		}
		return result
	case xdr.OperationTypeSetOptions:
		err := json.Unmarshal([]byte(o.DetailsString.String), &result.SetOptions)
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

type BasePayment struct {
	From                  string `json:"from,omitempty"`
	To                    string `json:"to,omitempty"`
	FromBalance           string `json:"from_balance,omitempty"`
	ToBalance             string `json:"to_balance,omitempty"`
	Amount                string `json:"amount"`
	UserDetails           string `json:"user_details,omitempty"`
	Asset                 string `json:"asset"`
	SourcePaymentFee      string `json:"source_payment_fee"`
	DestinationPaymentFee string `json:"destination_payment_fee"`
	SourceFixedFee        string `json:"source_fixed_fee"`
	DestinationFixedFee   string `json:"destination_fixed_fee"`
	SourcePaysForDest     bool   `json:"source_pays_for_dest"`
}

// Payment is the json resource representing a single operation whose type is
// Payment.
type PaymentDetails struct {
	BasePayment
	Subject   string `json:"subject,omitempty"`
	Reference string `json:"reference,omitempty"`
	Asset     string `json:"qasset"`
}
type SetOptionsDetails struct {
	HomeDomain    string `json:"home_domain,omitempty"`
	InflationDest string `json:"inflation_dest,omitempty"`

	MasterKeyWeight *int   `json:"master_key_weight,omitempty"`
	SignerKey       string `json:"signer_key,omitempty"`
	SignerWeight    *int   `json:"signer_weight,omitempty"`
	SignerType      *int   `json:"signer_type,omitempty"`
	SignerIdentity  *int   `json:"signer_identity,omitempty"`

	SetFlags    []int    `json:"set_flags,omitempty"`
	SetFlagsS   []string `json:"set_flags_s,omitempty"`
	ClearFlags  []int    `json:"clear_flags,omitempty"`
	ClearFlagsS []string `json:"clear_flags_s,omitempty"`

	LowThreshold  *int `json:"low_threshold,omitempty"`
	MedThreshold  *int `json:"med_threshold,omitempty"`
	HighThreshold *int `json:"high_threshold,omitempty"`
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
