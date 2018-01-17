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
			HomeDomain:                      "",
			InflationDest:                   "",
			MasterKeyWeight:                 uint32(*op.MasterWeight),
			SignerKey:                       op.Signer.PubKey.Address(),
			SignerWeight:                    uint32(op.Signer.Weight),
			SignerType:                      uint32(op.Signer.SignerType),
			SignerIdentity:                  uint32(op.Signer.Identity),
			SetFlags:                        nil,
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
				AccountType: int64(*op.Fee.AccountType), //ask about type int 64/32
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
			ExternalDetails: externalDetails, //TODO change to type string
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
		details["is_buy"] = op.IsBuy
		details["amount"] = amount.String(int64(op.Amount))
		details["price"] = amount.String(int64(op.Price))
		details["fee"] = amount.String(int64(op.Fee))
		details["offer_id"] = op.OfferId
		details["is_deleted"] = int64(op.OfferId) != 0
	case xdr.OperationTypeManageInvoice:
		op := c.Operation().Body.MustManageInvoiceOp()
		opResult := c.OperationResult().MustManageInvoiceResult()
		details["amount"] = amount.String(int64(op.Amount))
		details["receiver_balance"] = op.ReceiverBalance.AsString()
		details["sender"] = op.Sender.Address()
		details["invoice_id"] = opResult.Success.InvoiceId
		details["asset"] = string(opResult.Success.Asset)
	case xdr.OperationTypeReviewRequest:
		op := c.Operation().Body.MustReviewRequestOp()
		details["action"] = int32(op.Action)
		details["reason"] = string(op.Reason)
		details["request_hash"] = hex.EncodeToString(op.RequestHash[:])
		details["request_id"] = uint64(op.RequestId)
		details["request_type"] = int32(op.RequestDetails.RequestType)
	case xdr.OperationTypeManageAsset:
		op := c.Operation().Body.MustManageAssetOp()
		details["request_id"] = uint64(op.RequestId)
		details["action"] = int32(op.Request.Action)
	case xdr.OperationTypeCreatePreissuanceRequest:
		// no details needed
	case xdr.OperationTypeCreateIssuanceRequest:
		op := c.Operation().Body.MustCreateIssuanceRequestOp()
		details["fee_fixed"] = amount.StringU(uint64(op.Request.Fee.Fixed))
		details["fee_percent"] = amount.StringU(uint64(op.Request.Fee.Percent))
		details["reference"] = utf8.Scrub(string(op.Reference))
		details["amount"] = amount.StringU(uint64(op.Request.Amount))
		details["asset"] = string(op.Request.Asset)
		details["balance_id"] = op.Request.Receiver.AsString()

		var externalDetails map[string]interface{}
		// error is ignored on purpose, we should not block ingest in case of such error
		_ = json.Unmarshal([]byte(op.Request.ExternalDetails), &externalDetails)
		details["external_details"] = externalDetails
	case xdr.OperationTypeCreateSaleRequest:
		// no details needed
	case xdr.OperationTypeCheckSaleState:
		// no details needed
	default:
		panic(fmt.Errorf("Unknown operation type: %s", c.OperationType()))
	}
	return details
}
