package ingest

import (
	"fmt"

	"encoding/hex"
	"encoding/json"

	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/utf8"
	"gitlab.com/tokend/regources"
)

// operationDetails returns the details regarding the current operation, suitable
// for ingestion into a history_operation row
func (is *Session) operationDetails() map[string]interface{} {
	details := map[string]interface{}{}
	c := is.Cursor
	source := c.OperationSourceAccount()
	txFee, ok := c.Transaction().Result.Result.Ext.GetTransactionFee()
	if ok {
		details["operation_fee"] = amount.StringU(getOperationFee(txFee.OperationFees, c.Operation().Body.Type))
		details["operation_fee_asset"] = txFee.AssetCode
	}
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

	case xdr.OperationTypeManageKeyValue:
		op := c.Operation().Body.MustManageKeyValueOp()
		details["source"] = source
		details["key"] = op.Key
		details["action"] = op.Action
		if op.Action.Action != xdr.ManageKvActionRemove {
			details["value"] = op.Action.Value
		}

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
	case xdr.OperationTypeManageLimits:
		op := c.Operation().Body.MustManageLimitsOp()
		if op.Details.Action == xdr.ManageLimitsActionCreate {
			details["account_type"] = op.Details.LimitsCreateDetails.AccountType
			details["account_id"] = op.Details.LimitsCreateDetails.AccountId
			details["stats_op_type"] = op.Details.LimitsCreateDetails.StatsOpType
			details["asset_code"] = op.Details.LimitsCreateDetails.AssetCode
			details["is_convert_needed"] = op.Details.LimitsCreateDetails.IsConvertNeeded
			details["daily_out"] = op.Details.LimitsCreateDetails.DailyOut
			details["weekly_out"] = op.Details.LimitsCreateDetails.WeeklyOut
			details["monthly_out"] = op.Details.LimitsCreateDetails.MonthlyOut
			details["annual_out"] = op.Details.LimitsCreateDetails.AnnualOut
		}

		if op.Details.Action == xdr.ManageLimitsActionRemove {
			details["limit_id"] = op.Details.Id
		} else {
			details["limit_id"] = *c.OperationResult().MustManageLimitsResult().Success.Details.Id
		}
	case xdr.OperationTypeCreateManageLimitsRequest:
		op := c.Operation().Body.MustCreateManageLimitsRequestOp()
		limitsUpdateDetails, ok := op.ManageLimitsRequest.Ext.GetDetails()
		if ok {
			details["limits_manage_request_details"] = string(limitsUpdateDetails)
		}
		requestID, ok := op.Ext.GetRequestId()
		if ok {
			details["request_id"] = uint64(requestID)
		}
		details["limits_manage_request_document_hash"] = hex.EncodeToString(op.ManageLimitsRequest.DeprecatedDocumentHash[:])
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
		op := c.Operation().Body.MustManageOfferOp()
		opResult := c.OperationResult().MustManageOfferResult().MustSuccess()
		isDeleted := opResult.Offer.Effect == xdr.ManageOfferEffectDeleted
		details["is_buy"] = op.IsBuy
		details["amount"] = amount.String(int64(op.Amount))
		details["price"] = amount.String(int64(op.Price))
		details["fee"] = amount.String(int64(op.Fee))
		details["is_deleted"] = isDeleted
		if isDeleted {
			details["offer_id"] = op.OfferId
		} else {
			details["offer_id"] = opResult.Offer.Offer.OfferId
		}
		details["order_book_id"] = op.OrderBookId
		isSaleOffer := op.OrderBookId != 0
		if isSaleOffer {
			details["base_asset"] = getOfferBaseAsset(c.OperationChanges(), op.OrderBookId)
		}
	case xdr.OperationTypeManageInvoiceRequest:
		op := c.Operation().Body.MustManageInvoiceRequestOp()
		opResult := c.OperationResult().MustManageInvoiceRequestResult()
		switch op.Details.Action {
		case xdr.ManageInvoiceRequestActionCreate:
			details["amount"] = amount.String(int64(op.Details.InvoiceRequest.Amount))
			details["sender"] = op.Details.InvoiceRequest.Sender.Address()
			details["request_id"] = opResult.Success.Details.Response.RequestId
			details["asset"] = string(op.Details.InvoiceRequest.Asset)
		case xdr.ManageInvoiceRequestActionRemove:
			details["request_id"] = *op.Details.RequestId
		}
	case xdr.OperationTypeManageContractRequest:
		op := c.Operation().Body.MustManageContractRequestOp()
		opResult := c.OperationResult().MustManageContractRequestResult()
		switch op.Details.Action {
		case xdr.ManageContractRequestActionCreate:
			details["escrow"] = op.Details.ContractRequest.Escrow.Address()
			details["details"] = op.Details.ContractRequest.Details
			details["start_time"] = int64(op.Details.ContractRequest.StartTime)
			details["end_time"] = int64(op.Details.ContractRequest.EndTime)
			details["request_id"] = opResult.Success.Details.Response.RequestId
		case xdr.ManageContractRequestActionRemove:
			details["request_id"] = *op.Details.RequestId
		}
	case xdr.OperationTypeManageContract:
		op := c.Operation().Body.MustManageContractOp()
		details["contract_id"] = int64(op.ContractId)
		switch op.Data.Action {
		case xdr.ManageContractActionAddDetails:
			details["details"] = string(op.Data.MustDetails())
		case xdr.ManageContractActionConfirmCompleted:
			opResult := c.OperationResult().MustManageContractResult()
			details["is_completed"] = opResult.Response.Data.MustIsCompleted()
		case xdr.ManageContractActionStartDispute:
			details["dispute_reason"] = string(op.Data.MustDisputeReason())
		case xdr.ManageContractActionResolveDispute:
			details["is_revert"] = op.Data.MustIsRevert()
		}
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

		opResult := c.OperationResult().MustReviewRequestResult().MustSuccess()
		details["is_fulfilled"] = opResult.Fulfilled

		aSwapExtended, ok := opResult.TypeExt.GetASwapExtended()
		if !ok {
			break
		}
		details["atomic_swap_details"] = getAtomicSwapDetails(aSwapExtended)
	case xdr.OperationTypeManageAsset:
		op := c.Operation().Body.MustManageAssetOp()
		details["request_id"] = uint64(op.RequestId)
		details["action"] = int32(op.Request.Action)
		details["action_string"] = op.Request.Action.ShortString()
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

		allTasks, ok := op.Ext.GetAllTasks()
		if ok && allTasks != nil {
			details["all_tasks"] = *allTasks
		}
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
		details["source_fee_data"] = map[string]interface{}{
			"fixed_fee":                     amount.StringU(uint64(opResult.ActualDestinationPaymentFee.Fixed)),
			"actual_payment_fee":            amount.StringU(uint64(opResult.ActualDestinationPaymentFee.Percent)),
			"actual_payment_fee_asset_code": string(op.FeeData.SourceFee.FeeAsset),
		}
		details["destination_fee_data"] = map[string]interface{}{
			"fixed_fee":                     amount.StringU(uint64(op.FeeData.DestinationFee.FixedFee)),
			"actual_payment_fee":            amount.StringU(uint64(op.FeeData.DestinationFee.MaxPaymentFee)),
			"actual_payment_fee_asset_code": string(op.FeeData.DestinationFee.FeeAsset),
		}
		details["source_pays_for_dest"] = op.FeeData.SourcePaysForDest
		details["subject"] = op.Subject
		details["reference"] = utf8.Scrub(string(op.Reference))
		details["source_sent_universal"] = amount.StringU(uint64(opResult.SourceSentUniversal))
	case xdr.OperationTypeManageSale:
		op := c.Operation().Body.MustManageSaleOp()
		opRes := c.OperationResult().MustManageSaleResult().MustSuccess()
		details["sale_id"] = uint64(op.SaleId)
		details["action"] = op.Data.Action.ShortString()

		fulfilled, ok := opRes.Ext.GetFulfilled()
		if ok {
			details["fulfilled"] = fulfilled
		}
	case xdr.OperationTypeCancelSaleRequest:
		op := c.Operation().Body.MustCancelSaleCreationRequestOp()
		details["request_id"] = uint64(op.RequestId)
	case xdr.OperationTypeCreateAswapBidRequest:
		op := c.Operation().Body.MustCreateASwapBidCreationRequestOp()
		opRes := c.OperationResult().MustCreateASwapBidCreationRequestResult().
			MustSuccess()
		details["base_balance_id"] = op.Request.BaseBalance
		details["amount"] = amount.StringU(uint64(op.Request.Amount))

		var bidDetails map[string]interface{}
		// error is ignored on purpose, we should not block ingest in case of such error
		_ = json.Unmarshal([]byte(op.Request.Details), &bidDetails)
		details["details"] = bidDetails
		details["quote_assets"] = op.Request.QuoteAssets
		details["request_id"] = uint64(opRes.RequestId)
	case xdr.OperationTypeCancelAswapBid:
		op := c.Operation().Body.MustCancelASwapBidOp()

		details["bid_id"] = uint64(op.BidId)
	case xdr.OperationTypeCreateAswapRequest:
		op := c.Operation().Body.MustCreateASwapRequestOp()
		opRes := c.OperationResult().MustCreateASwapRequestResult().
			MustSuccess()
		details["bid_id"] = op.Request.BidId
		details["base_amount"] = amount.StringU(uint64(op.Request.BaseAmount))
		details["quote_asset"] = string(op.Request.QuoteAsset)
		details["request_id"] = opRes.RequestId
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

func getAtomicSwapDetails(atomicSwapExtendedResult xdr.ASwapExtended) map[string]interface{} {
	return map[string]interface{}{
		"bid_id":                          uint64(atomicSwapExtendedResult.BidId),
		"bid_owner_id":                    atomicSwapExtendedResult.BidOwnerId.Address(),
		"bid_owner_base_asset_balance_id": atomicSwapExtendedResult.BidOwnerBaseBalanceId.AsString(),
		"purchaser_id":                    atomicSwapExtendedResult.PurchaserId.Address(),
		"purchaser_base_asset_balance_id": atomicSwapExtendedResult.PurchaserBaseBalanceId.AsString(),
		"base_asset":                      string(atomicSwapExtendedResult.BaseAsset),
		"quote_asset":                     string(atomicSwapExtendedResult.QuoteAsset),
		"base_amount":                     regources.Amount(atomicSwapExtendedResult.BaseAmount),
		"quote_amount":                    regources.Amount(atomicSwapExtendedResult.QuoteAmount),
		"price":                           regources.Amount(atomicSwapExtendedResult.Price),
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

func getOfferBaseAsset(changes xdr.LedgerEntryChanges, saleId xdr.Uint64) xdr.AssetCode {
	for _, change := range changes {
		if change.Type != xdr.LedgerEntryChangeTypeUpdated {
			continue
		}
		data := change.Updated.Data
		if data.Type == xdr.LedgerEntryTypeSale && data.Sale.SaleId == saleId {
			return data.Sale.BaseAsset
		}
	}
	return xdr.AssetCode("")
}

func getOperationFee(opFees []xdr.OperationFee, opType xdr.OperationType) uint64 {
	for _, opFee := range opFees {
		if opFee.OperationType == opType {
			return uint64(opFee.Amount)
		}
	}
	return 0
}
