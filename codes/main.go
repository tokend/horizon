//Package codes is a helper package to help convert to transaction and operation result codes
//to strings used in horizon.
package codes

import (
	"bullioncoin.githost.io/development/go/xdr"
	"github.com/go-errors/errors"
	"gitlab.com/distributed_lab/logan"
)

// ErrUnknownCode is returned when an unexepcted value is provided to `String`
var ErrUnknownCode = errors.New("Unknown result code")

const (
	OpFeeMisMatch = "op_fee_mismatch"
	// OpSuccess is the string code used to specify the operation was successful
	OpSuccess = "op_success"
	// OpMalformed is the string code used to specify the operation was malformed
	// in some way.
	OpMalformed = "op_malformed"
	// OpUnderfunded is the string code used to specify the operation failed
	// due to a lack of funds.
	OpUnderfunded = "op_underfunded"

	// OpLowReserve is the string code used to specify the operation failed
	// because the account in question does not have enough balance to satisfy
	// what their new minimum balance would be.
	OpLowReserve = "op_low_reserve"

	// OpLineFull occurs when a payment would cause a destination account to
	// exceed their declared trust limit for the asset being sent.
	OpLineFull = "op_line_full"

	// OpRequestNotFound occurs when an operation specify request which does not exists
	OpRequestNotFound = "op_request_not_found"

	// OpRequestAlreadyReviewed occurs when request already have been reviewed
	OpRequestAlreadyReviewed = "op_request_already_reviewed"

	// OpNoAccount when no target account
	OpNoAccount = "op_no_account"

	// OpInvalidAmount when invalid amount
	OpInvalidAmount = "op_invalid_amount"

	OpAccountMismatch = "op_account_mismatch"

	OpAlreadyExists = "op_already_exists"

	OpAccountBlocked = "op_account_blocked"

	OpAccountNotFound = "op_account_not_found"

	OpBalanceNotFound = "op_balance_not_found"

	OpAssetNotFound = "op_asset_not_found"

	OpInvalidRequestID = "op_invalid_request_id"

	OpAssetMismatch = "op_asset_mismatch"

	OpBalanceMismatch = "op_balance_mismatch"

	OpReferenceDuplication = "op_reference_duplication"

	OpAccountTypeMismatch = "op_account_type_mismatched"

	OpStatsOverflow = "op_stats_overflow"

	OpLimitsExceeded = "op_limits_exceeded"

	OpDemurrageNotRequired = "op_demurrage_not_required"

	OpExchangePolicyViolated = "op_exchange_policy_violated"
	OpAssetPolicyViolated    = "op_exchange_policy_violated"
)

//String returns the appropriate string representation of the provided result code
func String(code interface{}) (string, error) {
	switch code := code.(type) {
	case xdr.TransactionResultCode:
		switch code {
		case xdr.TransactionResultCodeTxSuccess:
			return "tx_success", nil
		case xdr.TransactionResultCodeTxFailed:
			return "tx_failed", nil
		case xdr.TransactionResultCodeTxTooEarly:
			return "tx_too_early", nil
		case xdr.TransactionResultCodeTxTooLate:
			return "tx_too_late", nil
		case xdr.TransactionResultCodeTxMissingOperation:
			return "tx_missing_operation", nil
		case xdr.TransactionResultCodeTxBadAuth:
			return "tx_bad_auth", nil
		case xdr.TransactionResultCodeTxNoAccount:
			return "tx_no_source_account", nil
		case xdr.TransactionResultCodeTxBadAuthExtra:
			return "tx_bad_auth_extra", nil
		case xdr.TransactionResultCodeTxInternalError:
			return "tx_internal_error", nil
		case xdr.TransactionResultCodeTxAccountBlocked:
			return "account_blocked", nil
		case xdr.TransactionResultCodeTxDuplication:
			return "tx_duplication", nil
		}
	case xdr.OperationResultCode:
		switch code {
		case xdr.OperationResultCodeOpInner:
			return "op_inner", nil
		case xdr.OperationResultCodeOpBadAuth:
			return "op_bad_auth", nil
		case xdr.OperationResultCodeOpNoAccount:
			return "op_no_source_account", nil
		case xdr.OperationResultCodeOpNotAllowed:
			return "not_allowed", nil
		case xdr.OperationResultCodeOpAccountBlocked:
			return OpAccountBlocked, nil
		case xdr.OperationResultCodeOpNoCounterparty:
			return "op_no_counterparty", nil
		case xdr.OperationResultCodeOpCounterpartyBlocked:
			return "op_counterparty_blocked", nil
		case xdr.OperationResultCodeOpCounterpartyWrongType:
			return "op_counterparty_wrong_type", nil
		case xdr.OperationResultCodeOpBadAuthExtra:
			return "op_bad_auth_extra", nil

		}
	case xdr.CreateAccountResultCode:
		switch code {
		case xdr.CreateAccountResultCodeCreateAccountSuccess:
			return OpSuccess, nil
		case xdr.CreateAccountResultCodeCreateAccountMalformed:
			return OpMalformed, nil
		case xdr.CreateAccountResultCodeCreateAccountAccountTypeMismatched:
			return "op_account_type_mismatched", nil
		case xdr.CreateAccountResultCodeCreateAccountTypeNotAllowed:
			return "op_type_not_allowed", nil
		case xdr.CreateAccountResultCodeCreateAccountNameDuplication:
			return "op_name_duplication", nil
		case xdr.CreateAccountResultCodeCreateAccountReferrerNotFound:
			return "op_referrer_not_found", nil
		}
	case xdr.PaymentResultCode:
		switch code {
		case xdr.PaymentResultCodePaymentSuccess:
			return OpSuccess, nil
		case xdr.PaymentResultCodePaymentMalformed:
			return OpMalformed, nil
		case xdr.PaymentResultCodePaymentUnderfunded:
			return OpUnderfunded, nil
		case xdr.PaymentResultCodePaymentLineFull:
			return OpLineFull, nil
		case xdr.PaymentResultCodePaymentFeeMismatched:
			return "op_fee_mismatched", nil
		case xdr.PaymentResultCodePaymentBalanceNotFound:
			return OpBalanceNotFound, nil
		case xdr.PaymentResultCodePaymentBalanceAccountMismatched:
			return "op_balance_account_mismatched", nil
		case xdr.PaymentResultCodePaymentBalanceAssetsMismatched:
			return "balance_assets_mismatched", nil
		case xdr.PaymentResultCodePaymentSrcBalanceNotFound:
			return "src_balance_not_found", nil
		case xdr.PaymentResultCodePaymentReferenceDuplication:
			return OpReferenceDuplication, nil
		case xdr.PaymentResultCodePaymentStatsOverflow:
			return OpStatsOverflow, nil
		case xdr.PaymentResultCodePaymentLimitsExceeded:
			return OpLimitsExceeded, nil
		case xdr.PaymentResultCodePaymentNotAllowedByAssetPolicy:
			return OpAssetPolicyViolated, nil
		case xdr.PaymentResultCodePaymentInvoiceNotFound:
			return "op_invoice_not_found", nil
		case xdr.PaymentResultCodePaymentInvoiceWrongAmount:
			return "op_invoice_wrong_amount", nil
		case xdr.PaymentResultCodePaymentInvoiceBalanceMismatch:
			return "op_invoice_balance_mismatch", nil
		case xdr.PaymentResultCodePaymentInvoiceAccountMismatch:
			return "op_invoice_account_mismatch", nil
		case xdr.PaymentResultCodePaymentInvoiceAlreadyPaid:
			return "op_invoice_already_paid", nil
		}
	case xdr.SetOptionsResultCode:
		switch code {
		case xdr.SetOptionsResultCodeSetOptionsSuccess:
			return OpSuccess, nil
		case xdr.SetOptionsResultCodeSetOptionsTooManySigners:
			return "op_too_many_signers", nil
		case xdr.SetOptionsResultCodeSetOptionsThresholdOutOfRange:
			return "op_threshold_out_of_range", nil
		case xdr.SetOptionsResultCodeSetOptionsBadSigner:
			return "op_bad_signer", nil
		case xdr.SetOptionsResultCodeSetOptionsBalanceNotFound:
			return OpBalanceNotFound, nil
		case xdr.SetOptionsResultCodeSetOptionsTrustMalformed:
			return OpMalformed, nil
		case xdr.SetOptionsResultCodeSetOptionsTrustTooMany:
			return "op_too_many_trust_lines", nil

		}

	case xdr.ManageCoinsEmissionRequestResultCode:
		switch code {
		case xdr.ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestSuccess:
			return OpSuccess, nil
		case xdr.ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestInvalidAmount:
			return OpInvalidAmount, nil
		case xdr.ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestInvalidRequestId:
			return OpInvalidRequestID, nil
		case xdr.ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestNotFound:
			return OpRequestNotFound, nil
		case xdr.ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestAlreadyReviewed:
			return OpRequestAlreadyReviewed, nil
		case xdr.ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestAssetNotFound:
			return OpAssetNotFound, nil
		case xdr.ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestBalanceNotFound:
			return OpBalanceNotFound, nil
		case xdr.ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestAssetMismatch:
			return OpAssetMismatch, nil
		case xdr.ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestInvalidAsset:
			return "op_invalid_asset", nil
		case xdr.ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestReferenceDuplication:
			return "reference_duplication", nil
		case xdr.ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestLineFull:
			return OpLineFull, nil
		case xdr.ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestInvalidReference:
			return "invalid_reference", nil
		case xdr.ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestNotAllowedByExchangePolicy:
			return OpExchangePolicyViolated, nil
		}

	case xdr.ReviewCoinsEmissionRequestResultCode:
		switch code {
		case xdr.ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestSuccess:
			return OpSuccess, nil
		case xdr.ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestInvalidReason:
			return "op_invalid_reason", nil
		case xdr.ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestNotFound:
			return OpRequestNotFound, nil
		case xdr.ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestNotEqual:
			return "op_request_not_equal", nil
		case xdr.ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestAlreadyReviewed:
			return OpRequestAlreadyReviewed, nil
		case xdr.ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestMalformed:
			return OpMalformed, nil
		case xdr.ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestNotEnoughPreemissions:
			return "not_enough_preemissions", nil
		case xdr.ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestLineFull:
			return OpLineFull, nil
		case xdr.ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestAssetNotFound:
			return OpAssetNotFound, nil
		case xdr.ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestBalanceNotFound:
			return OpBalanceNotFound, nil
		case xdr.ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestReferenceDuplication:
			return "op_reference_duplication", nil
		}

	case xdr.SetFeesResultCode:
		switch code {
		case xdr.SetFeesResultCodeSetFeesSuccess:
			return OpSuccess, nil
		case xdr.SetFeesResultCodeSetFeesInvalidAmount:
			return OpInvalidAmount, nil
		case xdr.SetFeesResultCodeSetFeesInvalidFeeType:
			return "invalid_fee_type", nil
		case xdr.SetFeesResultCodeSetFeesAssetNotFound:
			return OpAssetNotFound, nil
		case xdr.SetFeesResultCodeSetFeesInvalidAsset:
			return "invalid_asset", nil
		case xdr.SetFeesResultCodeSetFeesMalformed:
			return OpMalformed, nil
		case xdr.SetFeesResultCodeSetFeesMalformedRange:
			return "op_malformed_range", nil
		case xdr.SetFeesResultCodeSetFeesRangeOverlap:
			return "op_range_overlap", nil
		case xdr.SetFeesResultCodeSetFeesNotFound:
			return "op_not_found", nil
		case xdr.SetFeesResultCodeSetFeesSubTypeNotExist:
			return "op_sub_type_not_exist", nil
		}
	case xdr.ManageAccountResultCode:
		switch code {
		case xdr.ManageAccountResultCodeManageAccountSuccess:
			return OpSuccess, nil
		case xdr.ManageAccountResultCodeManageAccountNotFound:
			return OpAccountNotFound, nil
		case xdr.ManageAccountResultCodeManageAccountMalformed:
			return OpMalformed, nil
		case xdr.ManageAccountResultCodeManageAccountNotAllowed:
			return "not_allowed", nil
		case xdr.ManageAccountResultCodeManageAccountTypeMismatch:
			return "type_mismatch", nil
		}
	case xdr.ForfeitResultCode:
		switch code {
		case xdr.ForfeitResultCodeForfeitSuccess:
			return OpSuccess, nil
		case xdr.ForfeitResultCodeForfeitMalformed:
			return OpMalformed, nil
		case xdr.ForfeitResultCodeForfeitBalanceNotFound:
			return OpBalanceNotFound, nil
		case xdr.ForfeitResultCodeForfeitUnderfunded:
			return OpUnderfunded, nil
		case xdr.ForfeitResultCodeForfeitStatsOverflow:
			return OpStatsOverflow, nil
		case xdr.ForfeitResultCodeForfeitLimitsExceeded:
			return OpLimitsExceeded, nil
		}
	case xdr.ManageForfeitRequestResultCode:
		switch code {
		case xdr.ManageForfeitRequestResultCodeManageForfeitRequestSuccess:
			return OpSuccess, nil
		case xdr.ManageForfeitRequestResultCodeManageForfeitRequestUnderfunded:
			return OpUnderfunded, nil
		case xdr.ManageForfeitRequestResultCodeManageForfeitRequestInvalidAmount:
			return OpInvalidAmount, nil
		case xdr.ManageForfeitRequestResultCodeManageForfeitRequestLineFull:
			return OpLineFull, nil
		case xdr.ManageForfeitRequestResultCodeManageForfeitRequestBalanceMismatch:
			return OpBalanceMismatch, nil
		case xdr.ManageForfeitRequestResultCodeManageForfeitRequestStatsOverflow:
			return OpStatsOverflow, nil
		case xdr.ManageForfeitRequestResultCodeManageForfeitRequestLimitsExceeded:
			return OpLimitsExceeded, nil
		case xdr.ManageForfeitRequestResultCodeManageForfeitRequestReviewerNotFound:
			return "op_reviewer_not_found", nil
		case xdr.ManageForfeitRequestResultCodeManageForfeitRequestInvalidDetails:
			return "op_invalid_details", nil
		case xdr.ManageForfeitRequestResultCodeManageForfeitRequestBalanceRequiresReview:
			return "op_exchange_requires_review_request_with_reviewer_not_available", nil
		}

	case xdr.RecoverResultCode:
		switch code {
		case xdr.RecoverResultCodeRecoverSuccess:
			return OpSuccess, nil
		case xdr.RecoverResultCodeRecoverMalformed:
			return OpMalformed, nil
		case xdr.RecoverResultCodeRecoverOldSignerNotFound:
			return "op_old_signer_not_found", nil
		case xdr.RecoverResultCodeRecoverSignerAlreadyExists:
			return "op_signer_already_exists", nil

		}
	case xdr.ManageBalanceResultCode:
		switch code {
		case xdr.ManageBalanceResultCodeManageBalanceSuccess:
			return OpSuccess, nil
		case xdr.ManageBalanceResultCodeManageBalanceMalformed:
			return OpMalformed, nil
		case xdr.ManageBalanceResultCodeManageBalanceNotFound:
			return OpBalanceNotFound, nil
		case xdr.ManageBalanceResultCodeManageBalanceDestinationNotFound:
			return "op_destination_not_found", nil
		case xdr.ManageBalanceResultCodeManageBalanceAlreadyExists:
			return OpAlreadyExists, nil
		case xdr.ManageBalanceResultCodeManageBalanceAnotherExchange:
			return "op_another_exchange", nil
		case xdr.ManageBalanceResultCodeManageBalanceAssetNotFound:
			return OpAssetNotFound, nil
		case xdr.ManageBalanceResultCodeManageBalanceInvalidAsset:
			return "op_invalid_asset", nil
		case xdr.ManageBalanceResultCodeManageBalanceNotAllowedByExchangePolicy:
			return OpExchangePolicyViolated, nil

		}
	case xdr.ReviewPaymentRequestResultCode:
		switch code {
		case xdr.ReviewPaymentRequestResultCodeReviewPaymentRequestSuccess:
			return OpSuccess, nil
		case xdr.ReviewPaymentRequestResultCodeReviewPaymentRequestNotFound:
			return OpRequestNotFound, nil
		case xdr.ReviewPaymentRequestResultCodeReviewPaymentRequestLineFull:
			return OpLineFull, nil
		case xdr.ReviewPaymentRequestResultCodeReviewPaymentDemurrageRejectionNotAllowed:
			return "demurrage_rejection_not_allowed", nil
		}

	case xdr.ManageAssetResultCode:
		switch code {
		case xdr.ManageAssetResultCodeManageAssetSuccess:
			return OpSuccess, nil
		case xdr.ManageAssetResultCodeManageAssetNotFound:
			return OpAssetNotFound, nil
		case xdr.ManageAssetResultCodeManageAssetAlreadyExists:
			return OpAlreadyExists, nil
		case xdr.ManageAssetResultCodeManageAssetMalformed:
			return OpMalformed, nil
		case xdr.ManageAssetResultCodeManageAssetTokenAlreadyExists:
			return "op_token_already_exists", nil
		case xdr.ManageAssetResultCodeManageAssetTokenTokenAlredySet:
			return "op_asset_token_already_set", nil
		}

	case xdr.DemurrageResultCode:
		switch code {
		case xdr.DemurrageResultCodeDemurrageSuccess:
			return OpSuccess, nil
		case xdr.DemurrageResultCodeDemurrageAssetNotFound:
			return OpAssetNotFound, nil
		case xdr.DemurrageResultCodeDemurrageInvalidAsset:
			return "op_invalid_asset", nil
		case xdr.DemurrageResultCodeDemurrageNotRequired:
			return OpDemurrageNotRequired, nil
		case xdr.DemurrageResultCodeDemurrageStatsOverflow:
			return OpStatsOverflow, nil
		case xdr.DemurrageResultCodeDemurrageLimitsExceeded:
			return OpLimitsExceeded, nil
		}

	case xdr.UploadPreemissionsResultCode:
		switch code {
		case xdr.UploadPreemissionsResultCodeUploadPreemissionsSuccess:
			return OpSuccess, nil
		case xdr.UploadPreemissionsResultCodeUploadPreemissionsMalformed:
			return OpMalformed, nil
		case xdr.UploadPreemissionsResultCodeUploadPreemissionsSerialDuplication:
			return "serial_duplication", nil
		case xdr.UploadPreemissionsResultCodeUploadPreemissionsMalformedPreemissions:
			return "malformed_preemissions", nil
		case xdr.UploadPreemissionsResultCodeUploadPreemissionsAssetNotFound:
			return OpAssetNotFound, nil
		case xdr.UploadPreemissionsResultCodeUploadPreemissionsLineFull:
			return OpLineFull, nil
		}
	case xdr.SetLimitsResultCode:
		switch code {
		case xdr.SetLimitsResultCodeSetLimitsSuccess:
			return OpSuccess, nil
		case xdr.SetLimitsResultCodeSetLimitsMalformed:
			return OpMalformed, nil
		}
	case xdr.DirectDebitResultCode:
		switch code {
		case xdr.DirectDebitResultCodeDirectDebitSuccess:
			return OpSuccess, nil
		case xdr.DirectDebitResultCodeDirectDebitMalformed:
			return OpMalformed, nil
		case xdr.DirectDebitResultCodeDirectDebitUnderfunded:
			return OpUnderfunded, nil
		case xdr.DirectDebitResultCodeDirectDebitLineFull:
			return OpLineFull, nil
		case xdr.DirectDebitResultCodeDirectDebitFeeMismatched:
			return "op_fee_mismatched", nil
		case xdr.DirectDebitResultCodeDirectDebitBalanceNotFound:
			return OpBalanceNotFound, nil
		case xdr.DirectDebitResultCodeDirectDebitBalanceAssetsMismatched:
			return "balance_assets_mismatched", nil
		case xdr.DirectDebitResultCodeDirectDebitBalanceAccountMismatched:
			return "op_balance_account_mismatched", nil
		case xdr.DirectDebitResultCodeDirectDebitSrcBalanceNotFound:
			return "src_balance_not_found", nil
		case xdr.DirectDebitResultCodeDirectDebitReferenceDuplication:
			return OpReferenceDuplication, nil
		case xdr.DirectDebitResultCodeDirectDebitStatsOverflow:
			return OpStatsOverflow, nil
		case xdr.DirectDebitResultCodeDirectDebitLimitsExceeded:
			return OpLimitsExceeded, nil
		case xdr.DirectDebitResultCodeDirectDebitNotAllowedByAssetPolicy:
			return OpAssetPolicyViolated, nil
		case xdr.DirectDebitResultCodeDirectDebitNoTrust:
			return "op_no_trust", nil
		}
	case xdr.ManageAssetPairResultCode:
		switch code {
		case xdr.ManageAssetPairResultCodeManageAssetPairSuccess:
			return OpSuccess, nil
		case xdr.ManageAssetPairResultCodeManageAssetPairNotFound:
			return "op_asset_pair_not_found", nil
		case xdr.ManageAssetPairResultCodeManageAssetPairAlreadyExists:
			return OpAlreadyExists, nil
		case xdr.ManageAssetPairResultCodeManageAssetPairMalformed:
			return OpMalformed, nil
		case xdr.ManageAssetPairResultCodeManageAssetPairInvalidAsset:
			return "op_invalid_asset", nil
		case xdr.ManageAssetPairResultCodeManageAssetPairInvalidAction:
			return "op_invalid_action", nil
		case xdr.ManageAssetPairResultCodeManageAssetPairInvalidPolicies:
			return "op_invalid_policies", nil
		case xdr.ManageAssetPairResultCodeManageAssetPairAssetNotFound:
			return "op_asset_not_dound", nil
		}
	case xdr.ManageOfferResultCode:
		switch code {
		case xdr.ManageOfferResultCodeManageOfferSuccess:
			return OpSuccess, nil
		case xdr.ManageOfferResultCodeManageOfferMalformed:
			return OpMalformed, nil
		case xdr.ManageOfferResultCodeManageOfferPairNotTraded:
			return "op_pair_not_traded", nil
		case xdr.ManageOfferResultCodeManageOfferBalanceNotFound:
			return "balance_not_found", nil
		case xdr.ManageOfferResultCodeManageOfferRequiresReview:
			return "op_requiers_review", nil
		case xdr.ManageOfferResultCodeManageOfferUnderfunded:
			return OpUnderfunded, nil
		case xdr.ManageOfferResultCodeManageOfferCrossSelf:
			return "op_offer_cross_self", nil
		case xdr.ManageOfferResultCodeManageOfferOverflow:
			return "op_overflow", nil
		case xdr.ManageOfferResultCodeManageOfferAssetPairNotTradable:
			return "op_pair_not_traded", nil
		case xdr.ManageOfferResultCodeManageOfferPhysicalPriceRestriction:
			return "physical_price_restriction", nil
		case xdr.ManageOfferResultCodeMaangeOfferCurrentPriceRestriction:
			return "op_current_price_restriction", nil
		case xdr.ManageOfferResultCodeManageOfferNotFound:
			return "op_not_found", nil
		case xdr.ManageOfferResultCodeManageOfferInvalidPercentFee:
			return "op_invalid_percent_fee", nil
		case xdr.ManageOfferResultCodeManageOfferDirectBuyNotAllowed:
			return "op_offer_direct_buy_not_allowed", nil
		case xdr.ManageOfferResultCodeManageOfferInsuffisientPrice:
			return "op_offer_insuffisient_price", nil
		}
	case xdr.ManageInvoiceResultCode:
		switch code {
		case xdr.ManageInvoiceResultCodeManageInvoiceSuccess:
			return OpSuccess, nil
		case xdr.ManageInvoiceResultCodeManageInvoiceMalformed:
			return OpMalformed, nil
		case xdr.ManageInvoiceResultCodeManageInvoiceBalanceNotFound:
			return OpBalanceNotFound, nil
		case xdr.ManageInvoiceResultCodeManageInvoiceOverflow:
			return "overflow", nil
		case xdr.ManageInvoiceResultCodeManageInvoiceNotFound:
			return "op_invoice_not_found", nil
		case xdr.ManageInvoiceResultCodeManageInvoiceTooManyInvoices:
			return "op_too_many_invoices", nil
		case xdr.ManageInvoiceResultCodeManageInvoiceCanNotDeleteInProgress:
			return "op_can_not_delete_in_progress", nil
		}
	}

	return "", errors.New(ErrUnknownCode)
}

// ForOperationResult returns the strong represtation used by horizon for the
// error code `opr`
func ForOperationResult(opr xdr.OperationResult) (string, error) {
	if opr.Code != xdr.OperationResultCodeOpInner {
		return String(opr.Code)
	}

	ir := opr.MustTr()
	var ic interface{}

	switch ir.Type {
	case xdr.OperationTypeCreateAccount:
		ic = ir.MustCreateAccountResult().Code
	case xdr.OperationTypePayment:
		ic = ir.MustPaymentResult().Code
	case xdr.OperationTypeSetOptions:
		ic = ir.MustSetOptionsResult().Code
	case xdr.OperationTypeManageCoinsEmissionRequest:
		ic = ir.MustManageCoinsEmissionRequestResult().Code
	case xdr.OperationTypeReviewCoinsEmissionRequest:
		ic = ir.MustReviewCoinsEmissionRequestResult().Code
	case xdr.OperationTypeSetFees:
		ic = ir.MustSetFeesResult().Code
	case xdr.OperationTypeManageAccount:
		ic = ir.MustManageAccountResult().Code
	case xdr.OperationTypeForfeit:
		ic = ir.MustForfeitResult().Code
	case xdr.OperationTypeManageForfeitRequest:
		ic = ir.MustManageForfeitRequestResult().Code
	case xdr.OperationTypeRecover:
		ic = ir.MustRecoverResult().Code
	case xdr.OperationTypeManageBalance:
		ic = ir.MustManageBalanceResult().Code
	case xdr.OperationTypeReviewPaymentRequest:
		ic = ir.MustReviewPaymentRequestResult().Code
	case xdr.OperationTypeManageAsset:
		ic = ir.MustManageAssetResult().Code
	case xdr.OperationTypeDemurrage:
		ic = ir.MustDemurrageResult().Code
	case xdr.OperationTypeUploadPreemissions:
		ic = ir.MustUploadPreemissionsResult().Code
	case xdr.OperationTypeSetLimits:
		ic = ir.MustSetLimitsResult().Code
	case xdr.OperationTypeDirectDebit:
		ic = ir.MustDirectDebitResult().Code
	case xdr.OperationTypeManageAssetPair:
		ic = ir.MustManageAssetPairResult().Code
	case xdr.OperationTypeManageOffer:
		ic = ir.MustManageOfferResult().Code
	case xdr.OperationTypeManageInvoice:
		ic = ir.MustManageInvoiceResult().Code
	}

	return String(ic)
}

func ForTxResult(txResult xdr.TransactionResult) (txResultCode string, opResultCodes []string, err error) {
	txResultCode, err = String(txResult.Result.Code)
	if err != nil {
		err = logan.Wrap(err, "Failed to convert to string tx result code")
		return
	}

	if txResult.Result.Results == nil {
		return
	}

	opResults := txResult.Result.MustResults()
	opResultCodes = make([]string, len(opResults))
	for i := range opResults {
		opResultCodes[i], err = ForOperationResult(opResults[i])
		if err != nil {
			err = logan.Wrap(err, "Failed to convert to string op result")
			return
		}
	}

	return
}
