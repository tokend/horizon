package codes

import "gitlab.com/tokend/go/xdr"

type codeProvider func(tr xdr.OperationResultTr) shortStr

var codeProviders = map[xdr.OperationType]codeProvider{
	xdr.OperationTypeCreateAccount:            func(ir xdr.OperationResultTr) shortStr { return ir.MustCreateAccountResult().Code },
	xdr.OperationTypeSetFees:                  func(ir xdr.OperationResultTr) shortStr { return ir.MustSetFeesResult().Code },
	xdr.OperationTypeCreateWithdrawalRequest:  func(ir xdr.OperationResultTr) shortStr { return ir.MustCreateWithdrawalRequestResult().Code },
	xdr.OperationTypeManageBalance:            func(ir xdr.OperationResultTr) shortStr { return ir.MustManageBalanceResult().Code },
	xdr.OperationTypeManageAsset:              func(ir xdr.OperationResultTr) shortStr { return ir.MustManageAssetResult().Code },
	xdr.OperationTypeManageLimits:             func(ir xdr.OperationResultTr) shortStr { return ir.MustManageLimitsResult().Code },
	xdr.OperationTypeManageAssetPair:          func(ir xdr.OperationResultTr) shortStr { return ir.MustManageAssetPairResult().Code },
	xdr.OperationTypeManageOffer:              func(ir xdr.OperationResultTr) shortStr { return ir.MustManageOfferResult().Code },
	xdr.OperationTypeManageInvoiceRequest:     func(ir xdr.OperationResultTr) shortStr { return ir.MustManageInvoiceRequestResult().Code },
	xdr.OperationTypeReviewRequest:            func(ir xdr.OperationResultTr) shortStr { return ir.MustReviewRequestResult().Code },
	xdr.OperationTypeCreatePreissuanceRequest: func(ir xdr.OperationResultTr) shortStr { return ir.MustCreatePreIssuanceRequestResult().Code },
	xdr.OperationTypeCreateIssuanceRequest:    func(ir xdr.OperationResultTr) shortStr { return ir.MustCreateIssuanceRequestResult().Code },
	xdr.OperationTypeCreateSaleRequest:        func(ir xdr.OperationResultTr) shortStr { return ir.MustCreateSaleCreationRequestResult().Code },
	xdr.OperationTypeCheckSaleState:           func(ir xdr.OperationResultTr) shortStr { return ir.MustCheckSaleStateResult().Code },
	xdr.OperationTypeManageExternalSystemAccountIdPoolEntry: func(ir xdr.OperationResultTr) shortStr {
		return ir.MustManageExternalSystemAccountIdPoolEntryResult().Code
	},
	xdr.OperationTypeBindExternalSystemAccountId: func(ir xdr.OperationResultTr) shortStr { return ir.MustBindExternalSystemAccountIdResult().Code },
	xdr.OperationTypeCreateAmlAlert:              func(ir xdr.OperationResultTr) shortStr { return ir.MustCreateAmlAlertRequestResult().Code },
	xdr.OperationTypeCreateChangeRoleRequest:     func(ir xdr.OperationResultTr) shortStr { return ir.MustCreateChangeRoleRequestResult().Code },
	xdr.OperationTypePayment:                     func(ir xdr.OperationResultTr) shortStr { return ir.MustPaymentResult().Code },
	xdr.OperationTypeManageSale:                  func(ir xdr.OperationResultTr) shortStr { return ir.MustManageSaleResult().Code },
	xdr.OperationTypeManageKeyValue:              func(ir xdr.OperationResultTr) shortStr { return ir.MustManageKeyValueResult().Code },
	xdr.OperationTypeCreateManageLimitsRequest:   func(ir xdr.OperationResultTr) shortStr { return ir.MustCreateManageLimitsRequestResult().Code },
	xdr.OperationTypeManageContractRequest:       func(ir xdr.OperationResultTr) shortStr { return ir.MustManageContractRequestResult().Code },
	xdr.OperationTypeManageContract:              func(ir xdr.OperationResultTr) shortStr { return ir.MustManageContractResult().Code },
	xdr.OperationTypeCancelSaleRequest:           func(ir xdr.OperationResultTr) shortStr { return ir.MustCancelSaleCreationRequestResult().Code },
	xdr.OperationTypePayout:                      func(ir xdr.OperationResultTr) shortStr { return ir.MustPayoutResult().Code },
	xdr.OperationTypeManageAccountRole:           func(ir xdr.OperationResultTr) shortStr { return ir.MustManageAccountRoleResult().Code },
	xdr.OperationTypeManageAccountRule:           func(ir xdr.OperationResultTr) shortStr { return ir.MustManageAccountRuleResult().Code },
	xdr.OperationTypeCreateAtomicSwapBidRequest:  func(ir xdr.OperationResultTr) shortStr { return ir.MustCreateAtomicSwapBidRequestResult().Code },
	xdr.OperationTypeCancelAtomicSwapAsk:         func(ir xdr.OperationResultTr) shortStr { return ir.MustCancelAtomicSwapAskResult().Code },
	xdr.OperationTypeCreateAtomicSwapAskRequest:  func(ir xdr.OperationResultTr) shortStr { return ir.MustCreateAtomicSwapAskRequestResult().Code },
	xdr.OperationTypeManageSigner:                func(ir xdr.OperationResultTr) shortStr { return ir.MustManageSignerResult().Code },
	xdr.OperationTypeManageSignerRole:            func(ir xdr.OperationResultTr) shortStr { return ir.MustManageSignerRoleResult().Code },
	xdr.OperationTypeManageSignerRule:            func(ir xdr.OperationResultTr) shortStr { return ir.MustManageSignerRuleResult().Code },
	xdr.OperationTypeStamp:                       func(ir xdr.OperationResultTr) shortStr { return ir.MustStampResult().Code },
	xdr.OperationTypeLicense:                     func(ir xdr.OperationResultTr) shortStr { return ir.MustLicenseResult().Code },
	xdr.OperationTypeManageCreatePollRequest:     func(ir xdr.OperationResultTr) shortStr { return ir.MustManageCreatePollRequestResult().Code },
	xdr.OperationTypeManagePoll:                  func(ir xdr.OperationResultTr) shortStr { return ir.MustManagePollResult().Code },
	xdr.OperationTypeManageVote:                  func(ir xdr.OperationResultTr) shortStr { return ir.MustManageVoteResult().Code },
	xdr.OperationTypeManageAccountSpecificRule:   func(ir xdr.OperationResultTr) shortStr { return ir.MustManageAccountSpecificRuleResult().Code },
	xdr.OperationTypeCancelChangeRoleRequest:     func(ir xdr.OperationResultTr) shortStr { return ir.MustCancelChangeRoleRequestResult().Code },
	xdr.OperationTypeRemoveAssetPair:             func(ir xdr.OperationResultTr) shortStr { return ir.MustRemoveAssetPairResult().Code },
	xdr.OperationTypeInitiateKycRecovery:         func(ir xdr.OperationResultTr) shortStr { return ir.MustInitiateKycRecoveryResult().Code },
	xdr.OperationTypeCreateKycRecoveryRequest:    func(ir xdr.OperationResultTr) shortStr { return ir.MustCreateKycRecoveryRequestResult().Code },
}
