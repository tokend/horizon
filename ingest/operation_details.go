package ingest

import (
	"fmt"

	"encoding/hex"

	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon/utf8"
)

// operationDetails returns the details regarding the current operation, suitable
// for ingestion into a history_operation row
func (is *Session) operationDetails() map[string]interface{} {
	details := map[string]interface{}{}
	c := is.Cursor
	source := c.OperationSourceAccount()
	switch c.OperationType() {
	case xdr.OperationTypeCreateAccount:
		op := c.Operation().Body.MustCreateAccountOp()
		details["funder"] = source.Address()
		details["account"] = op.Destination.Address()
		details["account_type"] = int32(op.AccountType)
	case xdr.OperationTypePayment:
		op := c.Operation().Body.MustPaymentOp()
		opResult := c.OperationResult().MustPaymentResult()
		details["from"] = source.Address()
		details["to"] = opResult.PaymentResponse.Destination.Address()
		details["from_balance"] = op.SourceBalanceId.AsString()
		details["to_balance"] = op.DestinationBalanceId.AsString()
		details["amount"] = amount.String(int64(op.Amount))
		details["source_payment_fee"] = amount.String(int64(op.FeeData.SourceFee.PaymentFee))
		details["destination_payment_fee"] = amount.String(int64(op.FeeData.DestinationFee.PaymentFee))
		details["source_fixed_fee"] = amount.String(int64(op.FeeData.SourceFee.FixedFee))
		details["destination_fixed_fee"] = amount.String(int64(op.FeeData.DestinationFee.FixedFee))
		details["source_pays_for_dest"] = op.FeeData.SourcePaysForDest
		details["subject"] = op.Subject
		details["reference"] = op.Reference
		details["asset"] = opResult.PaymentResponse.Asset
	case xdr.OperationTypeSetOptions:
		op := c.Operation().Body.MustSetOptionsOp()

		if op.MasterWeight != nil {
			details["master_key_weight"] = *op.MasterWeight
		}

		if op.LowThreshold != nil {
			details["low_threshold"] = *op.LowThreshold
		}

		if op.MedThreshold != nil {
			details["med_threshold"] = *op.MedThreshold
		}

		if op.HighThreshold != nil {
			details["high_threshold"] = *op.HighThreshold
		}

		if op.Signer != nil {
			details["signer_key"] = op.Signer.PubKey.Address()
			details["signer_weight"] = op.Signer.Weight
			details["signer_type"] = op.Signer.SignerType
			details["signer_identity"] = op.Signer.Identity
		}
	case xdr.OperationTypeSetFees:
		op := c.Operation().Body.MustSetFeesOp()
		if op.Fee != nil {
			accountID := ""
			if op.Fee.AccountId != nil {
				accountID = op.Fee.AccountId.Address()
			}
			accountType := op.Fee.AccountType
			details["fee"] = map[string]interface{}{
				"asset_code":   string(op.Fee.Asset),
				"fixed_fee":    amount.String(int64(op.Fee.FixedFee)),
				"percent_fee":  amount.String(int64(op.Fee.PercentFee)),
				"fee_type":     int64(op.Fee.FeeType),
				"account_id":   accountID,
				"account_type": accountType,
				"subtype":      int64(op.Fee.Subtype),
				"lower_bound":  int64(op.Fee.LowerBound),
				"upper_bound":  int64(op.Fee.UpperBound),
			}
		}

	case xdr.OperationTypeManageAccount:
		op := c.Operation().Body.MustManageAccountOp()
		details["account"] = op.Account.Address()
		details["block_reasons_to_add"] = op.BlockReasonsToAdd
		details["block_reasons_to_remove"] = op.BlockReasonsToRemove
	case xdr.OperationTypeCreateWithdrawalRequest:
		op := c.Operation().Body.MustCreateWithdrawalRequestOp()
		request := op.Request
		details["amount"] = amount.StringU(uint64(request.Amount))
		details["balance"] = request.Balance.AsString()
		details["fee_fixed"] = amount.StringU(uint64(request.Fee.Fixed))
		details["fee_percent"] = amount.StringU(uint64(request.Fee.Percent))
		details["external_details"] = request.ExternalDetails
		details["dest_asset"] = request.Details.AutoConversion.DestAsset
		details["dest_amount"] = amount.StringU(uint64(request.Details.AutoConversion.ExpectedAmount))
	case xdr.OperationTypeRecover:
		op := c.Operation().Body.MustRecoverOp()
		details["account"] = op.Account.Address()
		details["old_signer"] = op.OldSigner
		details["new_signer"] = op.NewSigner
	case xdr.OperationTypeManageBalance:
		op := c.Operation().Body.MustManageBalanceOp()
		details["destination"] = op.Destination
		details["action"] = op.Action
	case xdr.OperationTypeReviewPaymentRequest:
		op := c.Operation().Body.MustReviewPaymentRequestOp()
		details["payment_id"] = op.PaymentId
		details["accept"] = op.Accept
		if op.RejectReason != nil {
			details["reject_reason"] = *op.RejectReason
		}
	case xdr.OperationTypeSetLimits:
		op := c.Operation().Body.MustSetLimitsOp()
		details["account_type"] = op.AccountType
		details["account"] = op.Account
	case xdr.OperationTypeDirectDebit:
		op := c.Operation().Body.MustDirectDebitOp().PaymentOp
		opResult := c.OperationResult().MustDirectDebitResult().MustSuccess()
		details["from"] = source.Address()
		details["to"] = opResult.PaymentResponse.Destination.Address()
		details["from_balance"] = op.SourceBalanceId.AsString()
		details["to_balance"] = op.DestinationBalanceId.AsString()
		details["amount"] = amount.String(int64(op.Amount))
		details["source_payment_fee"] = amount.String(int64(op.FeeData.SourceFee.PaymentFee))
		details["destination_payment_fee"] = amount.String(int64(op.FeeData.DestinationFee.PaymentFee))
		details["source_fixed_fee"] = amount.String(int64(op.FeeData.SourceFee.FixedFee))
		details["destination_fixed_fee"] = amount.String(int64(op.FeeData.DestinationFee.FixedFee))
		details["source_pays_for_dest"] = op.FeeData.SourcePaysForDest
		details["subject"] = op.Subject
		details["reference"] = op.Reference
		details["asset"] = opResult.PaymentResponse.Asset
	case xdr.OperationTypeManageAssetPair:
		op := c.Operation().Body.MustManageAssetPairOp()
		details["base_asset"] = op.Base
		details["quote_asset"] = op.Quote
		details["physical_price"] = amount.String(int64(op.PhysicalPrice))
		details["physical_price_correction"] = amount.String(int64(op.PhysicalPriceCorrection))
		details["max_price_step"] = amount.String(int64(op.MaxPriceStep))
		details["policies_i"] = int32(op.Policies)
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
		details["reference"] = utf8.Scrub(string(op.Reference))
		details["amount"] = amount.StringU(uint64(op.Request.Amount))
		details["asset"] = string(op.Request.Asset)
		details["balance_id"] = op.Request.Receiver.AsString()
		details["external_details"] = op.Request.ExternalDetails
	case xdr.OperationTypeCreateSaleRequest:
		// no details needed
	default:
		panic(fmt.Errorf("Unknown operation type: %s", c.OperationType()))
	}
	return details
}
