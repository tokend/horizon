package operations

import (
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/render/hal"
	"golang.org/x/net/context"
)

// New creates a new operation resource, finding the appropriate type to use
// based upon the row's type.
func New(
	ctx context.Context, row history.Operation,
	participants []*history.Participant, public bool,
) (result hal.Pageable, err error) {

	base := Base{}
	err = base.Populate(ctx, row, participants, public)
	if err != nil {
		return
	}

	switch row.Type {
	case xdr.OperationTypeCreateAccount:
		d := row.Details().CreateAccount
		e := CreateAccount{
			Base:        base,
			Funder:      d.Funder,
			Account:     d.Account,
			AccountType: d.AccountType,
		}
		if public {
			e.Funder = ""
			e.Account = ""
		}
		result = e
	case xdr.OperationTypePayment:
		d := row.Details().Payment
		e := Payment{
			Base: base,
			BasePayment: BasePayment{
				From:                  d.BasePayment.From,
				To:                    d.BasePayment.To,
				FromBalance:           d.BasePayment.FromBalance,
				ToBalance:             d.BasePayment.ToBalance,
				Amount:                d.BasePayment.Amount,
				UserDetails:           d.BasePayment.UserDetails,
				Asset:                 d.BasePayment.Asset,
				SourcePaymentFee:      d.BasePayment.SourcePaymentFee,
				DestinationPaymentFee: d.BasePayment.DestinationPaymentFee,
				SourceFixedFee:        d.BasePayment.SourceFixedFee,
				DestinationFixedFee:   d.BasePayment.DestinationFixedFee,
				SourcePaysForDest:     d.BasePayment.SourcePaysForDest,
			},
			Subject:   d.Subject,
			Reference: d.Reference,
			Asset:     d.Asset,
		}
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
		d := row.Details().SetOptions
		e := SetOptions{
			Base:            base,
			HomeDomain:      d.HomeDomain,
			InflationDest:   d.InflationDest,
			MasterKeyWeight: d.MasterKeyWeight,
			SignerKey:       d.SignerKey,
			SignerWeight:    d.SignerWeight,
			SignerType:      d.SignerType,
			SignerIdentity:  d.SignerIdentity,
			SetFlags:        d.SetFlags,
			SetFlagsS:       d.SetFlagsS,
			ClearFlags:      d.ClearFlags,
			ClearFlagsS:     d.ClearFlagsS,
			LowThreshold:    d.LowThreshold,
			MedThreshold:    d.MedThreshold,
			HighThreshold:   d.HighThreshold,
		}
		if public {
			e.SignerKey = ""
		}
		result = e
	case xdr.OperationTypeSetFees:
		d := row.Details().SetFees
		e := SetFees{
			Base: base,
			Fee: &Fee{
				AssetCode:   d.Fee.AssetCode,
				FixedFee:    d.Fee.FixedFee,
				PercentFee:  d.Fee.PercentFee,
				FeeType:     d.Fee.FeeType,
				AccountID:   d.Fee.AccountID,
				AccountType: d.Fee.AccountType,
				Subtype:     d.Fee.Subtype,
				LowerBound:  d.Fee.LowerBound,
				UpperBound:  d.Fee.UpperBound,
			},
		}
		if public {
			if e.Fee != nil {
				e.Fee.AccountID = ""
			}
		}
		result = e
	case xdr.OperationTypeManageAccount:
		d := row.Details().ManageAccount
		e := ManageAccount{
			Base:                 base,
			Account:              d.Account,
			BlockReasonsToAdd:    d.BlockReasonsToAdd,
			BlockReasonsToRemove: d.BlockReasonsToRemove,
		}
		if public {
			e.Account = ""
		}
		result = e
	case xdr.OperationTypeCreateWithdrawalRequest:
		d := row.Details().CreateWithdrawalRequest
		e := CreateWithdrawalRequest{
			Base:            base,
			Amount:          d.Amount,
			Balance:         d.Balance,
			FeeFixed:        d.FeeFixed,
			FeePercent:      d.FeePercent,
			ExternalDetails: d.ExternalDetails,
			DestAsset:       d.DestAsset,
			DestAmount:      d.DestAmount,
		}
		if public {
			e.ExternalDetails = nil
		}
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
	case xdr.OperationTypeCreateIssuanceRequest:
		e := CreateIssuanceRequest{Base: base}
		err = row.UnmarshalDetails(&e)
		if public {
			e.ExternalDetails = nil
		}
		result = e
	default:
		result = base
	}

	return
}

// CreateAccount is the json resource representing a single operation whose type
// is CreateAccount.
type CreateAccount struct {
	Base
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
	Fee *Fee `json:"fee"`
}

type ManagerOffer struct {
	Base
	IsBuy     bool   `json:"is_buy"`
	Amount    string `json:"amount"`
	Price     string `json:"price"`
	Fee       string `json:"fee"`
	OfferId   int64  `json:"offer_id"`
	IsDeleted bool   `json:"is_deleted"`
}
