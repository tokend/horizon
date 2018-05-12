package ingest

import (
	"fmt"

	"encoding/hex"
	"encoding/json"

	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
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
		details["reference"] = utf8.Scrub(string(op.Reference))
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

		if op.LimitsUpdateRequestData != nil {
			details["limits_update_request_document_hash"] = hex.EncodeToString(op.LimitsUpdateRequestData.DocumentHash[:])
		}

	case xdr.OperationTypeSetFees:
		op := c.Operation().Body.MustSetFeesOp()
		if op.Fee != nil {
			accountID := ""
			if op.Fee.AccountId != nil {
				accountID = op.Fee.AccountId.Address()
			}

			feeAssetString := ""
			feeAsset, ok := op.Fee.Ext.GetFeeAsset()
			if ok {
				feeAssetString = string(feeAsset)
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
				"fee_asset":    feeAssetString,
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

		var externalDetails map[string]interface{}
		// error is ignored on purpose, we should not block ingest in case of such error
		_ = json.Unmarshal([]byte(request.ExternalDetails), &externalDetails)
		details["external_details"] = externalDetails

		details["dest_asset"] = request.Details.AutoConversion.DestAsset
		details["dest_amount"] = amount.StringU(uint64(request.Details.AutoConversion.ExpectedAmount))
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
		details["reference"] = utf8.Scrub(string(op.Reference))
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
		details["action"] = op.Action.ShortString()
		details["reason"] = string(op.Reason)
		details["request_hash"] = hex.EncodeToString(op.RequestHash[:])
		details["request_id"] = uint64(op.RequestId)
		details["request_type"] = op.RequestDetails.RequestType.ShortString()
		if op.Action == xdr.ReviewRequestOpActionApprove {
			details["is_fulfilled"] = hasDeletedReviewableRequest(c.OperationChanges())
		}
		details["details"] = getReviewRequestOpDetails(op.RequestDetails)
	case xdr.OperationTypeManageAsset:
		op := c.Operation().Body.MustManageAssetOp()
		details["request_id"] = uint64(op.RequestId)
		details["action"] = int32(op.Request.Action)
	case xdr.OperationTypeCreatePreissuanceRequest:
		// no details needed
	case xdr.OperationTypeCreateIssuanceRequest:
		op := c.Operation().Body.MustCreateIssuanceRequestOp()
		opResult := c.OperationResult().MustCreateIssuanceRequestResult().MustSuccess()
		details["fee_fixed"] = amount.StringU(uint64(opResult.Fee.Fixed))
		details["fee_percent"] = amount.StringU(uint64(opResult.Fee.Percent))
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
		op := c.Operation().Body.MustCheckSaleStateOp()
		opResult := c.OperationResult().MustCheckSaleStateResult().MustSuccess()
		details["sale_id"] = uint64(op.SaleId)
		details["effect"] = opResult.Effect.Effect.String()
		// no details needed
	case xdr.OperationTypeManageExternalSystemAccountIdPoolEntry:
		// no details needed
	case xdr.OperationTypeBindExternalSystemAccountId:
		// no details needed
	case xdr.OperationTypeCreateAmlAlert:
		op := c.Operation().Body.MustCreateAmlAlertRequestOp()
		details["amount"] = amount.StringU(uint64(op.AmlAlertRequest.Amount))
		details["balance_id"] = op.AmlAlertRequest.BalanceId.AsString()
		details["reason"] = op.AmlAlertRequest.Reason
		details["reference"] = op.Reference
	case xdr.OperationTypeCreateKycRequest:
		op := c.Operation().Body.MustCreateUpdateKycRequestOp()
		opResult := c.OperationResult().MustCreateUpdateKycRequestResult().MustSuccess()
		details["request_id"] = uint64(opResult.RequestId)
		details["account_to_update_kyc"] = op.UpdateKycRequestData.AccountToUpdateKyc.Address()
		details["account_type_to_set"] = int32(op.UpdateKycRequestData.AccountTypeToSet)
		details["kyc_level_to_set"] = uint32(op.UpdateKycRequestData.KycLevelToSet)

		var kycData map[string]interface{}
		// error is ignored on purpose, we should not block ingest in case of such error
		_ = json.Unmarshal([]byte(op.UpdateKycRequestData.KycData), &kycData)
		details["kyc_data"] = kycData

		if op.UpdateKycRequestData.AllTasks != nil {
			details["all_tasks"] = *op.UpdateKycRequestData.AllTasks
		}
	case xdr.OperationTypePaymentV2:
		op := c.Operation().Body.MustPaymentOpV2()
		opResult := c.OperationResult().MustPaymentV2Result().MustPaymentV2Response()
		details["payment_id"] = uint64(opResult.PaymentId)
		details["from"] = source.Address()
		details["to"] = opResult.Destination.Address()
		details["from_balance"] = op.SourceBalanceId.AsString()
		details["to_balance"] = opResult.DestinationBalanceId.AsString()
		details["amount"] = amount.StringU(uint64(op.Amount))
		details["asset"] = string(opResult.Asset)
		details["source_fee_data"] = map[string]interface{} {
			"fixed_fee": amount.StringU(uint64(op.FeeData.SourceFee.FixedFee)),
			"actual_payment_fee": amount.StringU(uint64(opResult.ActualSourcePaymentFee)),
			"actual_payment_fee_asset_code": string(op.FeeData.SourceFee.FeeAsset),
		}
		details["destination_fee_data"] = map[string]interface{} {
			"fixed_fee": amount.StringU(uint64(op.FeeData.DestinationFee.FixedFee)),
			"actual_payment_fee": amount.StringU(uint64(opResult.ActualDestinationPaymentFee)),
			"actual_payment_fee_asset_code": string(op.FeeData.DestinationFee.FeeAsset),
		}
		details["source_pays_for_dest"] = op.FeeData.SourcePaysForDest
		details["subject"] = op.Subject
		details["reference"] = utf8.Scrub(string(op.Reference))
		details["source_sent_universal"] = amount.StringU(uint64(opResult.SourceSentUniversal))
	default:
		panic(fmt.Errorf("Unknown operation type: %s", c.OperationType()))
	}
	return details
}

func getReviewRequestOpDetails(requestDetails xdr.ReviewRequestOpRequestDetails) map[string]interface{} {
	return map[string]interface{}{
		"request_type": requestDetails.RequestType.ShortString(),
		"update_kyc":   getUpdateKYCDetails(requestDetails.UpdateKyc),
	}
}

func getUpdateKYCDetails(details *xdr.UpdateKycDetails) map[string]interface{} {
	if details == nil {
		return nil
	}

	var externalDetails map[string]interface{}
	_ = json.Unmarshal([]byte(details.ExternalDetails), externalDetails)
	return map[string]interface{}{
		"external_details": externalDetails,
		"tasks_to_add":     uint32(details.TasksToAdd),
		"tasks_to_remove":  uint32(details.TasksToRemove),
	}
}
