package ingest

import (
	"fmt"

	"encoding/hex"
	"encoding/json"

	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/utf8"
)

// operationDetails returns the details regarding the current operation, suitable
// for ingestion into a history_operation row
func (is *Session) operationDetails() interface{} {
	details := map[string]interface{}{}
	c := is.Cursor
	source := c.OperationSourceAccount()

	operationDetails := history.OperationDetails{}
	switch c.OperationType() {
	case xdr.OperationTypeCreateAccount:
		op := c.Operation().Body.MustCreateAccountOp()

		operationDetails.Type = xdr.OperationTypeCreateAccount

		operationDetails.CreateAccount = &history.CreateAccountDetails{
			Funder:      source.Address(),
			Account:     op.Destination.Address(),
			AccountType: int32(op.AccountType),
		}
		return operationDetails
	case xdr.OperationTypePayment:
		op := c.Operation().Body.MustPaymentOp()
		opResult := c.OperationResult().MustPaymentResult()

		operationDetails.Type = xdr.OperationTypePayment

		operationDetails.Payment = &history.PaymentDetails{
			BasePayment: history.BasePayment{
				From:                  source.Address(),
				To:                    opResult.PaymentResponse.Destination.Address(),
				FromBalance:           op.SourceBalanceId.AsString(),
				ToBalance:             op.DestinationBalanceId.AsString(),
				Amount:                amount.String(int64(op.Amount)),
				Asset:                 string(opResult.PaymentResponse.Asset),
				SourcePaymentFee:      amount.String(int64(op.FeeData.SourceFee.PaymentFee)),
				DestinationPaymentFee: amount.String(int64(op.FeeData.DestinationFee.PaymentFee)),
				SourceFixedFee:        amount.String(int64(op.FeeData.SourceFee.FixedFee)),
				DestinationFixedFee:   amount.String(int64(op.FeeData.DestinationFee.FixedFee)),
				SourcePaysForDest:     op.FeeData.SourcePaysForDest,
			},
			Subject:    string(op.Subject),
			Reference:  string(op.Reference),
			QuoteAsset: string(opResult.PaymentResponse.Asset),
		}

		return operationDetails
	case xdr.OperationTypeSetOptions:
		op := c.Operation().Body.MustSetOptionsOp()

		operationDetails.Type = xdr.OperationTypeSetOptions

		operationDetails.SetOptions = &history.SetOptionsDetails{
			HomeDomain:                      "", //TODO Delete or set this fields
			InflationDest:                   "",
			MasterKeyWeight:                 uint32(*op.MasterWeight),
			SignerKey:                       op.Signer.PubKey.Address(),
			SignerWeight:                    uint32(op.Signer.Weight),
			SignerType:                      uint32(op.Signer.SignerType),
			SignerIdentity:                  uint32(op.Signer.Identity),
			SetFlags:                        nil, //TODO Delete or set this fields
			SetFlagsS:                       nil,
			ClearFlags:                      nil,
			ClearFlagsS:                     nil,
			LowThreshold:                    uint32(*op.LowThreshold),
			MedThreshold:                    uint32(*op.MedThreshold),
			HighThreshold:                   uint32(*op.HighThreshold),
			LimitsUpdateRequestDocumentHash: hex.EncodeToString(op.LimitsUpdateRequestData.DocumentHash[:]),
		}

		return operationDetails
	case xdr.OperationTypeSetFees:
		op := c.Operation().Body.MustSetFeesOp()

		operationDetails.Type = xdr.OperationTypeSetFees

		operationDetails.SetFees = &history.SetFeesDetails{
			Fee: &history.FeeDetails{
				AssetCode:   string(op.Fee.Asset),
				FixedFee:    amount.String(int64(op.Fee.FixedFee)),
				PercentFee:  amount.String(int64(op.Fee.PercentFee)),
				FeeType:     int64(op.Fee.FeeType),
				AccountID:   op.Fee.AccountId.Address(),
				AccountType: int32(*op.Fee.AccountType),
				Subtype:     int64(op.Fee.Subtype),
				LowerBound:  int64(op.Fee.LowerBound),
				UpperBound:  int64(op.Fee.UpperBound),
			},
		}

		return operationDetails
	case xdr.OperationTypeManageAccount:
		op := c.Operation().Body.MustManageAccountOp()

		operationDetails.Type = xdr.OperationTypeManageAccount

		operationDetails.ManageAccount = &history.ManageAccountDetails{
			Account:              op.Account.Address(),
			BlockReasonsToAdd:    uint32(op.BlockReasonsToAdd),
			BlockReasonsToRemove: uint32(op.BlockReasonsToRemove),
		}

		return operationDetails
	case xdr.OperationTypeCreateWithdrawalRequest:
		op := c.Operation().Body.MustCreateWithdrawalRequestOp()
		request := op.Request

		operationDetails.Type = xdr.OperationTypeCreateWithdrawalRequest

		var externalDetails map[string]interface{}
		// error is ignored on purpose, we should not block ingest in case of such error
		_ = json.Unmarshal([]byte(request.ExternalDetails), &externalDetails)

		operationDetails.CreateWithdrawalRequest = &history.CreateWithdrawalRequestDetails{
			Amount:          amount.StringU(uint64(request.Amount)),
			Balance:         request.Balance.AsString(),
			FeeFixed:        amount.StringU(uint64(request.Fee.Fixed)),
			FeePercent:      amount.StringU(uint64(request.Fee.Percent)),
			ExternalDetails: externalDetails,
			DestAsset:       string(request.Details.AutoConversion.DestAsset),
			DestAmount:      amount.StringU(uint64(request.Details.AutoConversion.ExpectedAmount)),
		}

		return operationDetails
	case xdr.OperationTypeManageBalance:
		op := c.Operation().Body.MustManageBalanceOp()
		operationDetails.Type = xdr.OperationTypeManageBalance

		//added new struct in resource/main.go and in OperationDetails
		operationDetails.ManageBalance = &history.ManageBalanceDetails{
			Destination: op.Destination.Address(),
			Action:      int32(op.Action),
		}

		return operationDetails
	case xdr.OperationTypeReviewPaymentRequest:
		op := c.Operation().Body.MustReviewPaymentRequestOp()

		operationDetails.Type = xdr.OperationTypeReviewPaymentRequest

		operationDetails.ReviewPaymentRequest = &history.ReviewPaymentRequestDetails{
			PaymentID:    int64(op.PaymentId),
			Accept:       op.Accept,
			RejectReason: string(*op.RejectReason),
		}

		return operationDetails
	case xdr.OperationTypeSetLimits:
		op := c.Operation().Body.MustSetLimitsOp()

		operationDetails.Type = xdr.OperationTypeSetLimits

		details["account_type"] = op.AccountType
		details["account"] = op.Account
	case xdr.OperationTypeDirectDebit:
		op := c.Operation().Body.MustDirectDebitOp().PaymentOp
		opResult := c.OperationResult().MustDirectDebitResult().MustSuccess()

		operationDetails.Type = xdr.OperationTypeDirectDebit

		operationDetails.DirectDebit = &history.DirectDebitDetails{
			From:                  source.Address(),
			To:                    opResult.PaymentResponse.Destination.Address(),
			FromBalance:           op.SourceBalanceId.AsString(),
			ToBalance:             op.DestinationBalanceId.AsString(),
			Amount:                amount.String(int64(op.Amount)),
			SourcePaymentFee:      amount.String(int64(op.FeeData.SourceFee.PaymentFee)),
			DestinationPaymentFee: amount.String(int64(op.FeeData.DestinationFee.PaymentFee)),
			SourceFixedFee:        amount.String(int64(op.FeeData.SourceFee.FixedFee)),
			DestinationFixedFee:   amount.String(int64(op.FeeData.DestinationFee.FixedFee)),
			SourcePaysForDest:     op.FeeData.SourcePaysForDest,
			Subject:               string(op.Subject),
			Reference:             string(op.Reference),
			AssetCode:             string(opResult.PaymentResponse.Asset),
		}

		return operationDetails
	case xdr.OperationTypeManageAssetPair:
		op := c.Operation().Body.MustManageAssetPairOp()

		operationDetails.Type = xdr.OperationTypeManageAssetPair

		operationDetails.ManageAssetPair = &history.ManageAssetPairDetails{
			BaseAsset:               string(op.Base),
			QuoteAsset:              string(op.Quote),
			PhysicalPrice:           amount.String(int64(op.PhysicalPrice)),
			PhysicalPriceCorrection: amount.String(int64(op.PhysicalPriceCorrection)),
			MaxPriceStep:            amount.String(int64(op.MaxPriceStep)),
			Policies:                int32(op.Policies),
		}

		return operationDetails
	case xdr.OperationTypeManageOffer:
		op := c.Operation().Body.ManageOfferOp

		operationDetails.Type = xdr.OperationTypeManageOffer

		operationDetails.ManagerOffer = &history.ManagerOfferDetails{
			IsBuy:     op.IsBuy,
			Amount:    amount.String(int64(op.Amount)),
			Price:     amount.String(int64(op.Price)),
			Fee:       amount.String(int64(op.Fee)),
			OfferId:   uint64(op.OfferId),
			IsDeleted: int64(op.OfferId) != 0,
		}

		return operationDetails
	case xdr.OperationTypeManageInvoice:
		op := c.Operation().Body.MustManageInvoiceOp()
		opResult := c.OperationResult().MustManageInvoiceResult()

		operationDetails.Type = xdr.OperationTypeManageInvoice

		operationDetails.ManageInvoice = &history.ManageInvoiceDetails{
			Amount:          amount.String(int64(op.Amount)),
			ReceiverBalance: op.ReceiverBalance.AsString(),
			Sender:          op.Sender.Address(),
			InvoiceID:       uint64(opResult.Success.InvoiceId),
			RejectReason:    nil, //TODO Delete or set this field
			Asset:           string(opResult.Success.Asset),
		}

		return operationDetails
	case xdr.OperationTypeReviewRequest:
		op := c.Operation().Body.MustReviewRequestOp()

		operationDetails.Type = xdr.OperationTypeReviewRequest

		operationDetails.ReviewRequest = &history.ReviewRequestDetails{
			Action:      int32(op.Action),
			Reason:      string(op.Reason),
			RequestHash: hex.EncodeToString(op.RequestHash[:]),
			RequestID:   uint64(op.RequestId),
			RequestType: int32(op.RequestDetails.RequestType),
		}

		return operationDetails
	case xdr.OperationTypeManageAsset:
		op := c.Operation().Body.MustManageAssetOp()

		operationDetails.Type = xdr.OperationTypeManageAsset

		operationDetails.ManageAsset = &history.ManageAssetDetails{
			RequestID: uint64(op.RequestId),
			Action:    int32(op.Request.Action),
		}

		return operationDetails
	case xdr.OperationTypeCreatePreissuanceRequest:
		// no details needed
	case xdr.OperationTypeCreateIssuanceRequest:
		op := c.Operation().Body.MustCreateIssuanceRequestOp()

		operationDetails.Type = xdr.OperationTypeCreateIssuanceRequest

		var externalDetails map[string]interface{}
		// error is ignored on purpose, we should not block ingest in case of such error
		_ = json.Unmarshal([]byte(op.Request.ExternalDetails), &externalDetails)

		operationDetails.CreateIssuanceRequest = &history.CreateIssuanceRequestDetails{
			Reference:       utf8.Scrub(string(op.Reference)),
			Amount:          amount.StringU(uint64(op.Request.Amount)),
			Asset:           string(op.Request.Asset),
			FeeFixed:        amount.StringU(uint64(op.Request.Fee.Fixed)),
			FeePercent:      amount.StringU(uint64(op.Request.Fee.Percent)),
			BalanceID:       op.Request.Receiver.AsString(),
			ExternalDetails: externalDetails,
		}

		return operationDetails
	case xdr.OperationTypeCreateSaleRequest:
		// no details needed
	case xdr.OperationTypeCheckSaleState:
		// no details needed
	default:
		panic(fmt.Errorf("Unknown operation type: %s", c.OperationType()))
	}
	return details
}
