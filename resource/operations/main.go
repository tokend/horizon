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
		d := row.GetDetails().CreateAccount
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
		d := row.GetDetails().Payment
		e := Payment{
			Base: base,
			BasePayment: BasePayment{
				From:                  d.BasePayment.From,
				To:                    d.BasePayment.To,
				FromBalance:           d.BasePayment.FromBalance,
				ToBalance:             d.BasePayment.ToBalance,
				Amount:                d.BasePayment.Amount,
				Asset:                 d.BasePayment.Asset,
				SourcePaymentFee:      d.BasePayment.SourcePaymentFee,
				DestinationPaymentFee: d.BasePayment.DestinationPaymentFee,
				SourceFixedFee:        d.BasePayment.SourceFixedFee,
				DestinationFixedFee:   d.BasePayment.DestinationFixedFee,
				SourcePaysForDest:     d.BasePayment.SourcePaysForDest,
			},
			Subject:    d.Subject,
			Reference:  d.Reference,
			QuoteAsset: d.QuoteAsset,
		}
		if public {
			e.From = ""
			e.To = ""
			e.FromBalance = ""
			e.ToBalance = ""
			e.Subject = ""
			e.Reference = ""
		}
		result = e
	case xdr.OperationTypeSetOptions:
		d := row.GetDetails().SetOptions
		e := SetOptions{
			Base:                            base,
			MasterKeyWeight:                 d.MasterKeyWeight,
			SignerKey:                       d.SignerKey,
			LowThreshold:                    d.LowThreshold,
			MedThreshold:                    d.MedThreshold,
			HighThreshold:                   d.HighThreshold,
			LimitsUpdateRequestDocumentHash: d.LimitsUpdateRequestDocumentHash,
		}
		result = e
	case xdr.OperationTypeSetFees:
		d := row.GetDetails().SetFees
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
		result = e
	case xdr.OperationTypeManageAccount:
		d := row.GetDetails().ManageAccount
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
		d := row.GetDetails().CreateWithdrawalRequest
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

	case xdr.OperationTypeManageBalance:
		d := row.GetDetails().ManageBalance

		e := ManageBalance{
			Base:        base,
			Destination: d.Destination,
			Action:      d.Action,
		}

		if public {
			e.Destination = ""
		}
		result = e

	case xdr.OperationTypeReviewPaymentRequest:
		d := row.GetDetails().ReviewPaymentRequest

		e := ReviewPaymentRequest{
			Base:         base,
			PaymentID:    d.PaymentID,
			Accept:       d.Accept,
			RejectReason: d.RejectReason,
		}

		result = e

	case xdr.OperationTypeSetLimits:
		e := SetLimits{Base: base}
		result = e

	case xdr.OperationTypeDirectDebit:
		d := row.GetDetails().DirectDebit
		e := DirectDebit{
			Base:                  base,
			From:                  d.From,
			To:                    d.To,
			FromBalance:           d.FromBalance,
			ToBalance:             d.ToBalance,
			Amount:                d.Amount,
			SourcePaymentFee:      d.SourcePaymentFee,
			DestinationPaymentFee: d.DestinationPaymentFee,
			SourceFixedFee:        d.SourceFixedFee,
			DestinationFixedFee:   d.DestinationFixedFee,
			SourcePaysForDest:     d.SourcePaysForDest,
			Subject:               d.Subject,
			Reference:             d.Reference,
			AssetCode:             d.AssetCode,
		}

		if public {
			e.From = ""
			e.To = ""
			e.FromBalance = ""
			e.ToBalance = ""
			e.Subject = ""
			e.Reference = ""
		}

		result = e

	case xdr.OperationTypeManageInvoice:
		d := row.GetDetails().ManageInvoice
		e := ManageInvoice{
			Base:            base,
			Amount:          d.Amount,
			ReceiverBalance: d.ReceiverBalance,
			Sender:          d.Sender,
			InvoiceID:       d.InvoiceID,
			Asset:           d.Asset,
		}
		if public {
			e.ReceiverBalance = ""
			e.Sender = ""
		}
		result = e

	case xdr.OperationTypeManageAsset:
		d := row.GetDetails().ManageAsset

		e := ManageAsset{
			Base:      base,
			RequestID: d.RequestID,
			Action:    d.Action,
		}

		result = e

	case xdr.OperationTypeManageOffer:
		d := row.GetDetails().ManagerOffer
		e := ManagerOffer{
			Base:      base,
			IsBuy:     d.IsBuy,
			Amount:    d.Amount,
			Price:     d.Price,
			Fee:       d.Fee,
			OfferId:   d.OfferId,
			IsDeleted: d.IsDeleted,
		}
		result = e
	case xdr.OperationTypeManageAssetPair:
		d := row.GetDetails().ManageAssetPair
		e := ManageAssetPair{
			Base:                    base,
			BaseAsset:               d.BaseAsset,
			QuoteAsset:              d.QuoteAsset,
			PhysicalPrice:           d.PhysicalPrice,
			PhysicalPriceCorrection: d.PhysicalPriceCorrection,
			MaxPriceStep:            d.MaxPriceStep,
			Policies:                d.Policies,
		}
		result = e
	case xdr.OperationTypeCreateIssuanceRequest:
		d := row.GetDetails().CreateIssuanceRequest
		e := CreateIssuanceRequest{
			Base:            base,
			Reference:       d.Reference,
			Amount:          d.Amount,
			Asset:           d.Asset,
			FeeFixed:        d.FeeFixed,
			FeePercent:      d.FeePercent,
			ExternalDetails: d.ExternalDetails,
			BalanceID:       d.BalanceID,
		}
		if public {
			e.ExternalDetails = nil
		}
		result = e
	case xdr.OperationTypeReviewRequest:
		d := row.GetDetails().ReviewRequest

		e := ReviewRequest{
			Base:        base,
			Action:      d.Action,
			Reason:      d.Reason,
			RequestHash: d.RequestHash,
			RequestID:   d.RequestID,
			RequestType: d.RequestType,
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
	Subject    string `json:"subject,omitempty"`
	Reference  string `json:"reference,omitempty"`
	QuoteAsset string `json:"qasset"`
}

// SetOptions is the json resource representing a single operation whose type is
// SetOptions.
type SetOptions struct {
	Base
	MasterKeyWeight                 *uint32 `json:"master_key_weight"`
	SignerKey                       string  `json:"signer_key,omitempty"`
	LowThreshold                    *uint32 `json:"low_threshold,omitempty"`
	MedThreshold                    *uint32 `json:"med_threshold,omitempty"`
	HighThreshold                   *uint32 `json:"high_threshold,omitempty"`
	LimitsUpdateRequestDocumentHash string  `json:"limits_update_request_document_hash,omitempty"`
}

//SetFees is the json resource representing a single operation whose type
// is SetFees.

type Fee struct {
	AssetCode   string  `json:"asset_code"`
	FixedFee    string  `json:"fixed_fee"`
	PercentFee  string  `json:"percent_fee"`
	FeeType     int64   `json:"fee_type"`
	AccountID   *string `json:"account_id,omitempty"`
	AccountType *int32  `json:"account_type"`
	Subtype     int64   `json:"subtype"`
	LowerBound  int64   `json:"lower_bound"`
	UpperBound  int64   `json:"upper_bound"`
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
	OfferId   uint64 `json:"offer_id"`
	IsDeleted bool   `json:"is_deleted"`
}

type ManageBalance struct {
	Base
	Destination string `json:"destination"`
	Action      int32  `json:"action"`
}

type ReviewPaymentRequest struct {
	Base
	PaymentID    int64   `json:"payment_id"`
	Accept       bool    `json:"accept"`
	RejectReason *string `json:"reject_reason"`
}

type DirectDebit struct {
	Base
	From                  string `json:"from"`
	To                    string `json:"to"`
	FromBalance           string `json:"from_balance"`
	ToBalance             string `json:"to_balance"`
	Amount                string `json:"amount"`
	SourcePaymentFee      string `json:"source_payment_fee"`
	DestinationPaymentFee string `json:"destination_payment_fee"`
	SourceFixedFee        string `json:"source_fixed_fee"`
	DestinationFixedFee   string `json:"destination_fixed_fee"`
	SourcePaysForDest     bool   `json:"source_pays_for_dest"`
	Subject               string `json:"subject"`
	Reference             string `json:"reference"`
	AssetCode             string `json:"asset"`
}

type ManageAsset struct {
	Base
	RequestID uint64 `json:"request_id"`
	Action    int32  `json:"action"`
}

type ReviewRequest struct {
	Base
	Action      int32  `json:"action"`
	Reason      string `json:"reason"`
	RequestHash string `json:"request_hash"`
	RequestID   uint64 `json:"request_id"`
	RequestType int32  `json:"request_type"`
}
