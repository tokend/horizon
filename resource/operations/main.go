package operations

import (
	"time"

	"bullioncoin.githost.io/development/go/xdr"
	"gitlab.com/distributed_lab/tokend/horizon/db2/history"
	"gitlab.com/distributed_lab/tokend/horizon/render/hal"
	"gitlab.com/distributed_lab/tokend/horizon/resource/base"
	"golang.org/x/net/context"
)

// TypeNames maps from operation type to the string used to represent that type
// in horizon's JSON responses
var TypeNames = map[xdr.OperationType]string{
	xdr.OperationTypeCreateAccount:              "create_account",
	xdr.OperationTypePayment:                    "payment",
	xdr.OperationTypeSetOptions:                 "set_options",
	xdr.OperationTypeManageCoinsEmissionRequest: "manage_coins_emission_request",
	xdr.OperationTypeReviewCoinsEmissionRequest: "review_coins_emission_request",
	xdr.OperationTypeSetFees:                    "set_fees",
	xdr.OperationTypeManageAccount:              "manage_account",
	xdr.OperationTypeForfeit:                    "forfeit",
	xdr.OperationTypeManageForfeitRequest:       "manage_forfeit_request",
	xdr.OperationTypeRecover:                    "recover",
	xdr.OperationTypeManageBalance:              "manage_balance",
	xdr.OperationTypeReviewPaymentRequest:       "review_payment_request",
	xdr.OperationTypeManageAsset:                "manage_asset",
	xdr.OperationTypeDemurrage:                  "demurrage",
	xdr.OperationTypeUploadPreemissions:         "upload_pre-emissions",
	xdr.OperationTypeSetLimits:                  "set_limits",
	xdr.OperationTypeDirectDebit:                "direct_debit",
	xdr.OperationTypeManageAssetPair:            "manage_asset_pair",
	xdr.OperationTypeManageOffer:                "manage_offer",
	xdr.OperationTypeManageInvoice:              "manage_invoice",
}

// New creates a new operation resource, finding the appropriate type to use
// based upon the row's type.
func New(
	ctx context.Context, row history.Operation, participants []*history.Participant, public bool,
) (result hal.Pageable, err error) {

	base := Base{}
	err = base.Populate(ctx, row, participants, public)
	if err != nil {
		return
	}
	switch row.Type {
	case xdr.OperationTypeCreateAccount:
		e := CreateAccount{Base: base}
		err = row.UnmarshalDetails(&e)
		if public {
			e.Funder = ""
			e.Account = ""
			e.Referrer = nil
		}
		result = e
	case xdr.OperationTypePayment:
		e := Payment{Base: base}
		err = row.UnmarshalDetails(&e)
		if public {
			e.UserDetails = ""
			e.From = ""
			e.To = ""
			e.FromBalance = ""
			e.ToBalance = ""
			e.Subject = ""
			e.Reference = ""
		}
		result = e
	case xdr.OperationTypeSetOptions:
		e := SetOptions{Base: base}
		err = row.UnmarshalDetails(&e)
		if public {
			e.SignerKey = ""
		}
		result = e
	case xdr.OperationTypeManageCoinsEmissionRequest:
		e := ManageCoinsEmissionRequest{Base: base}
		err = row.UnmarshalDetails(&e)
		result = e
	case xdr.OperationTypeReviewCoinsEmissionRequest:
		e := ReviewCoinsEmissionRequest{Base: base}
		err = row.UnmarshalDetails(&e)
		if public {
			e.Reason = ""
			e.Issuer = ""
		}
		result = e
	case xdr.OperationTypeSetFees:
		e := SetFees{Base: base}
		err = row.UnmarshalDetails(&e)
		if public {
			if e.Fee != nil {
				e.Fee.AccountID = ""
			}
		}
		result = e
	case xdr.OperationTypeManageAccount:
		e := ManageAccount{Base: base}
		err = row.UnmarshalDetails(&e)
		if public {
			e.Account = ""
		}
		result = e
	case xdr.OperationTypeForfeit:
		e := Forfeit{Base: base}
		err = row.UnmarshalDetails(&e)
		result = e
	case xdr.OperationTypeManageForfeitRequest:
		e := ManageForfeitRequest{Base: base}
		err = row.UnmarshalDetails(&e)
		if public {
			e.UserDetails = ""
		}
		result = e
	case xdr.OperationTypeDemurrage:
		e := Demurrage{Base: base}
		err = row.UnmarshalDetails(&e)
		result = e
	case xdr.OperationTypeSetLimits:
		e := SetLimits{Base: base}
		err = row.UnmarshalDetails(&e)
		result = e
	case xdr.OperationTypeManageInvoice:
		e := ManageInvoice{Base: base}
		err = row.UnmarshalDetails(&e)
		if public {
			e.ReceiverBalance = ""
			e.Sender = ""
			e.RejectReason = nil
		}
		result = e
	case xdr.OperationTypeManageOffer:
		e := ManagerOffer{Base: base}
		err = row.UnmarshalDetails(&e)
		result = e
	case xdr.OperationTypeManageAssetPair:
		e := ManageAssetPair{Base: base}
		err = row.UnmarshalDetails(&e)
		result = e
	default:
		result = base
	}

	return
}

// Base represents the common attributes of an operation resource
type Base struct {
	Links struct {
		Self        hal.Link `json:"self"`
		Transaction hal.Link `json:"transaction"`
		Succeeds    hal.Link `json:"succeeds"`
		Precedes    hal.Link `json:"precedes"`
	} `json:"_links"`

	ID              string             `json:"id"`
	PT              string             `json:"paging_token"`
	TransactionID   string             `json:"transaction_id"`
	SourceAccount   string             `json:"source_account,omitempty"`
	Type            string             `json:"type"`
	TypeI           int32              `json:"type_i"`
	State           int32              `json:"state"`
	Identifier      string             `json:"identifier"`
	LedgerCloseTime time.Time          `json:"ledger_close_time"`
	Participants    []base.Participant `json:"participants,omitempty"`
}

// CreateAccount is the json resource representing a single operation whose type
// is CreateAccount.
type CreateAccount struct {
	Base
	Funder      string  `json:"funder,omitempty"`
	Account     string  `json:"account,omitempty"`
	AccountType int32   `json:"account_type"`
	Referrer    *string `json:"referrer,omitempty"`
}

type BasePayment struct {
	From                  string             `json:"from,omitempty"`
	To                    string             `json:"to,omitempty"`
	FromBalance           string             `json:"from_balance,omitempty"`
	ToBalance             string             `json:"to_balance,omitempty"`
	Amount                string             `json:"amount"`
	UserDetails           string             `json:"user_details,omitempty"`
	Asset                 string             `json:"asset"`
	SourcePaymentFee      string             `json:"source_payment_fee"`
	DestinationPaymentFee string             `json:"destination_payment_fee"`
	SourceFixedFee        string             `json:"source_fixed_fee"`
	DestinationFixedFee   string             `json:"destination_fixed_fee"`
	SourcePaysForDest     bool               `json:"source_pays_for_dest"`
	Items                 []base.ForfeitItem `json:"items,omitempty"`
}

// Payment is the json resource representing a single operation whose type is
// Payment.
type Payment struct {
	Base
	BasePayment
	Subject   string `json:"subject,omitempty"`
	Reference string `json:"reference,omitempty"`
	Asset     string `json:"qasset"`
}

// SetOptions is the json resource representing a single operation whose type is
// SetOptions.
type SetOptions struct {
	Base
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

//SetFees is the json resource representing a single operation whose type
// is SetFees.

type Fee struct {
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

type SetFees struct {
	Base
	Fee              *Fee   `json:"fee"`
	StorageFeePeriod *int64 `json:"storage_fee_period"`
	PayoutsPeriod    *int64 `json:"payouts_period"`
}

type ManagerOffer struct {
	Base
	IsBuy     bool   `json:"is_buy"`
	Amount    string `json:"amount"`
	Price     string `json:"price"`
	Fee       string `json:"fee"`
	IsDirect  bool   `json:"is_direct"`
	OfferId   int64  `json:"offer_id"`
	IsDeleted bool   `json:"is_deleted"`
}
