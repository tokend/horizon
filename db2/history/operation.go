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

type OperationDetails struct {
	Type                    xdr.OperationType
	CreateAccount           *CreateAccountDetails
	Payment                 *PaymentDetails
	SetOptions              *SetOptionsDetails
	SetFees                 *SetFeesDetails
	ManageAccount           *ManageAccountDetails
	CreateWithdrawalRequest *CreateWithdrawalRequestDetails
	SetLimits               *SetLimitsDetails
	ManageInvoice           *ManageInvoiceDetails
	ManagerOffer            *ManagerOfferDetails
	ManageAssetPair         *ManageAssetPairDetails
	CreateIssuanceRequest   *CreateIssuanceRequestDetails
}

func (o *Operation) Details() OperationDetails {
	result := OperationDetails{
		Type: o.Type,
	}

	switch result.Type {
	case xdr.OperationTypeCreateAccount:
		err := json.Unmarshal([]byte(o.DetailsString.String), &result.CreateAccount)
		if err != nil {
			logrus.WithError(err).Errorf("Error unmarshal operation details")
		}
		return result
	case xdr.OperationTypePayment:
		err := json.Unmarshal([]byte(o.DetailsString.String), &result.Payment)
		if err != nil {
			logrus.WithError(err).Errorf("Error unmarshal operation details")
		}
		return result
	case xdr.OperationTypeSetOptions:
		err := json.Unmarshal([]byte(o.DetailsString.String), &result.SetOptions)
		if err != nil {
			logrus.WithError(err).Errorf("Error unmarshal operation details")
		}
		return result
	case xdr.OperationTypeSetFees:
		result.SetFees = &SetFeesDetails{
			Fee: &FeeDetails{},
		}

		err := json.Unmarshal([]byte(o.DetailsString.String), &result.SetFees)
		if err != nil {
			logrus.WithError(err).Errorf("Error unmarshal operation details")
		}
		return result
	case xdr.OperationTypeManageAccount:
		err := json.Unmarshal([]byte(o.DetailsString.String), &result.ManageAccount)
		if err != nil {
			logrus.WithError(err).Errorf("Error unmarshal operation details")
		}
		return result
	case xdr.OperationTypeCreateWithdrawalRequest:
		err := json.Unmarshal([]byte(o.DetailsString.String), &result.CreateWithdrawalRequest)
		if err != nil {
			logrus.WithError(err).Errorf("Error unmarshal operation details")
		}
		return result
	case xdr.OperationTypeSetLimits:
		err := json.Unmarshal([]byte(o.DetailsString.String), &result.SetLimits)
		if err != nil {
			logrus.WithError(err).Errorf("Error unmarshal operation details")
		}
		return result
	case xdr.OperationTypeManageInvoice:
		err := json.Unmarshal([]byte(o.DetailsString.String), &result.ManageInvoice)
		if err != nil {
			logrus.WithError(err).Errorf("Error unmarshal operation details")
		}
		return result
	case xdr.OperationTypeManageOffer:
		err := json.Unmarshal([]byte(o.DetailsString.String), &result.ManagerOffer)
		if err != nil {
			logrus.WithError(err).Errorf("Error unmarshal operation details")
		}
		return result
	case xdr.OperationTypeManageAssetPair:
		err := json.Unmarshal([]byte(o.DetailsString.String), &result.ManageAssetPair)
		if err != nil {
			logrus.WithError(err).Errorf("Error unmarshal operation details")
		}
		return result
	case xdr.OperationTypeCreateIssuanceRequest:
		err := json.Unmarshal([]byte(o.DetailsString.String), &result.CreateIssuanceRequest)
		if err != nil {
			logrus.WithError(err).Errorf("Error unmarshal operation details")
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

type FeeDetails struct {
	AssetCode   string `json:"asset_code"`
	FixedFee    string `json:"fixed_fee"`
	PercentFee  string `json:"percent_fee"`
	FeeType     int64  `json:"fee_type"`
	AccountID   string `json:"account_id,omitempty"`
	AccountType int64  `json:"account_type"`
	Subtype     int64  `json:"subtype"`
	LowerBound  int64  `json:"lower_bound"`
	UpperBound  int64  `json:"upper_bound"`
}

type SetFeesDetails struct {
	Fee *FeeDetails `json:"fee"`
}

type ManageAccountDetails struct {
	Account              string `json:"account,omitempty"`
	BlockReasonsToAdd    uint32 `json:"block_reasons_to_add,omitempty"`
	BlockReasonsToRemove uint32 `json:"block_reasons_to_remove,omitempty"`
}

type CreateWithdrawalRequestDetails struct {
	Amount          string                 `json:"amount"`
	Balance         string                 `json:"balance"`
	FeeFixed        string                 `json:"fee_fixed"`
	FeePercent      string                 `json:"fee_percent"`
	ExternalDetails map[string]interface{} `json:"external_details"`
	DestAsset       string                 `json:"dest_asset"`
	DestAmount      string                 `json:"dest_amount"`
}

type SetLimitsDetails struct{}

type ManageInvoiceDetails struct {
	Amount          string  `json:"amount"`
	ReceiverBalance string  `json:"receiver_balance,omitempty"`
	Sender          string  `json:"sender,omitempty"`
	InvoiceID       int64   `json:"invoice_id"`
	RejectReason    *string `json:"reject_reason,omitempty"`
	Asset           string  `json:"asset"`
}

type ManagerOfferDetails struct {
	IsBuy     bool   `json:"is_buy"`
	Amount    string `json:"amount"`
	Price     string `json:"price"`
	Fee       string `json:"fee"`
	OfferId   int64  `json:"offer_id"`
	IsDeleted bool   `json:"is_deleted"`
}

type ManageAssetPairDetails struct {
	BaseAsset               string `json:"base_asset"`
	QuoteAsset              string `json:"quote_asset"`
	PhysicalPrice           string `json:"physical_price"`
	PhysicalPriceCorrection string `json:"physical_price_correction"`
	MaxPriceStep            string `json:"max_price_step"`
}

type CreateIssuanceRequestDetails struct {
	Reference       string                 `json:"reference"`
	Amount          string                 `json:"amount"`
	Asset           string                 `json:"asset"`
	FeeFixed        string                 `json:"fee_fixed"`
	FeePercent      string                 `json:"fee_percent"`
	ExternalDetails map[string]interface{} `json:"external_details"`
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
