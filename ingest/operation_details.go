package ingest

import (
	"fmt"

	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
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
		details["account_type"] = int32(op.Details.AccountType)
		if op.Referrer != nil {
			details["referrer"] = (*op.Referrer).Address()
		}
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
	case xdr.OperationTypeManageCoinsEmissionRequest:
		op := c.Operation().Body.MustManageCoinsEmissionRequestOp()
		opResult := c.OperationResult().MustManageCoinsEmissionRequestResult()
		details["request_id"] = opResult.ManageRequestInfo.RequestId
		details["fulfilled"] = opResult.ManageRequestInfo.Fulfilled
		details["amount"] = amount.String(int64(op.Amount))
		details["asset"] = string(op.Asset)
	case xdr.OperationTypeReviewCoinsEmissionRequest:
		op := c.Operation().Body.MustReviewCoinsEmissionRequestOp()
		details["amount"] = amount.String(int64(op.Request.Amount))
		details["issuer"] = op.Request.Issuer.Address()
		details["approved"] = op.Approve
		details["reason"] = string(op.Reason)
		details["asset"] = string(op.Request.Asset)
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

		if op.StorageFeePeriod != nil {
			details["storage_fee_period"] = int64(*op.StorageFeePeriod)
		}

		if op.PayoutsPeriod != nil {
			details["payout_period"] = int64(*op.PayoutsPeriod)
		}

	case xdr.OperationTypeManageAccount:
		op := c.Operation().Body.MustManageAccountOp()
		details["account"] = op.Account.Address()
		details["block_reasons_to_add"] = op.BlockReasonsToAdd
		details["block_reasons_to_remove"] = op.BlockReasonsToRemove
	case xdr.OperationTypeForfeit:
		op := c.Operation().Body.MustForfeitOp()
		details["target"] = op.Balance.AsString()
		details["amount"] = amount.String(int64(op.Amount))
	case xdr.OperationTypeManageForfeitRequest:
		op := c.Operation().Body.MustManageForfeitRequestOp()
		opResult := c.OperationResult().MustManageForfeitRequestResult()
		details["amount"] = amount.String(int64(op.Amount))
		details["asset"] = opResult.ForfeitRequestDetails.Asset
		details["balance"] = op.Balance.AsString()
		details["user_details"] = op.Details
		details["items"] = manageForfeitRequestToForfeitTimes(opResult)
		if opResult.ForfeitRequestDetails.Ext.Fees != nil {
			details["fixed_fee"] = amount.String(int64(opResult.ForfeitRequestDetails.Ext.Fees.FixedFee))
			details["percent_fee"] = amount.String(int64(opResult.ForfeitRequestDetails.Ext.Fees.PercentFee))
		}
	case xdr.OperationTypeRecover:
		op := c.Operation().Body.MustRecoverOp()
		details["account"] = op.Account.Address()
		details["old_signer"] = op.OldSigner
		details["new_signer"] = op.NewSigner
	case xdr.OperationTypeManageBalance:
		op := c.Operation().Body.MustManageBalanceOp()
		details["balance_id"] = op.BalanceId
		details["destination"] = op.Destination
		details["action"] = op.Action
	case xdr.OperationTypeReviewPaymentRequest:
		op := c.Operation().Body.MustReviewPaymentRequestOp()
		details["payment_id"] = op.PaymentId
		details["accept"] = op.Accept
		if op.RejectReason != nil {
			details["reject_reason"] = *op.RejectReason
		}
	case xdr.OperationTypeManageAsset:
		op := c.Operation().Body.MustManageAssetOp()
		details["code"] = op.Code
		details["action"] = op.Action
	case xdr.OperationTypeDemurrage:
		opResult := c.OperationResult().MustDemurrageResult()
		details["quantity"] = len(opResult.DemurrageInfo.Demurrages)
	case xdr.OperationTypeUploadPreemissions:
		op := c.Operation().Body.MustUploadPreemissionsOp()
		details["quantity"] = len(op.PreEmissions)
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
		details["is_direct"] = op.IsDirect
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
	default:
		panic(fmt.Errorf("Unknown operation type: %s", c.OperationType()))
	}
	return details
}
