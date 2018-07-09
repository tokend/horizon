package codes

import "gitlab.com/tokend/go/xdr"

type codeProvider func(tr xdr.OperationResultTr) shortStr

var codeProviders = map[xdr.OperationType]codeProvider{
	xdr.OperationTypeCreateAccount:                          func(ir xdr.OperationResultTr) shortStr { return ir.MustCreateAccountResult().Code },
	xdr.OperationTypePayment:                                func(ir xdr.OperationResultTr) shortStr { return ir.MustPaymentResult().Code },
	xdr.OperationTypeSetOptions:                             func(ir xdr.OperationResultTr) shortStr { return ir.MustSetOptionsResult().Code },
	xdr.OperationTypeSetFees:                                func(ir xdr.OperationResultTr) shortStr { return ir.MustSetFeesResult().Code },
	xdr.OperationTypeManageAccount:                          func(ir xdr.OperationResultTr) shortStr { return ir.MustManageAccountResult().Code },
	xdr.OperationTypeCreateWithdrawalRequest:                func(ir xdr.OperationResultTr) shortStr { return ir.MustCreateWithdrawalRequestResult().Code },
	xdr.OperationTypeManageBalance:                          func(ir xdr.OperationResultTr) shortStr { return ir.MustManageBalanceResult().Code },
	xdr.OperationTypeReviewPaymentRequest:                   func(ir xdr.OperationResultTr) shortStr { return ir.MustReviewPaymentRequestResult().Code },
	xdr.OperationTypeManageAsset:                            func(ir xdr.OperationResultTr) shortStr { return ir.MustManageAssetResult().Code },
	xdr.OperationTypeManageLimits:                           func(ir xdr.OperationResultTr) shortStr { return ir.MustManageLimitsResult().Code },
	xdr.OperationTypeDirectDebit:                            func(ir xdr.OperationResultTr) shortStr { return ir.MustDirectDebitResult().Code },
	xdr.OperationTypeManageAssetPair:                        func(ir xdr.OperationResultTr) shortStr { return ir.MustManageAssetPairResult().Code },
	xdr.OperationTypeManageOffer:                            func(ir xdr.OperationResultTr) shortStr { return ir.MustManageOfferResult().Code },
	xdr.OperationTypeManageInvoiceRequest:					 func(ir xdr.OperationResultTr) shortStr { return ir.MustManageInvoiceRequestResult().Code },
	xdr.OperationTypeReviewRequest:                          func(ir xdr.OperationResultTr) shortStr { return ir.MustReviewRequestResult().Code },
	xdr.OperationTypeCreatePreissuanceRequest:               func(ir xdr.OperationResultTr) shortStr { return ir.MustCreatePreIssuanceRequestResult().Code },
	xdr.OperationTypeCreateIssuanceRequest:                  func(ir xdr.OperationResultTr) shortStr { return ir.MustCreateIssuanceRequestResult().Code },
	xdr.OperationTypeCreateSaleRequest:                      func(ir xdr.OperationResultTr) shortStr { return ir.MustCreateSaleCreationRequestResult().Code },
	xdr.OperationTypeCheckSaleState:                         func(ir xdr.OperationResultTr) shortStr { return ir.MustCheckSaleStateResult().Code },
	xdr.OperationTypeManageExternalSystemAccountIdPoolEntry: func(ir xdr.OperationResultTr) shortStr { return ir.MustManageExternalSystemAccountIdPoolEntryResult().Code },
	xdr.OperationTypeBindExternalSystemAccountId:            func(ir xdr.OperationResultTr) shortStr { return ir.MustBindExternalSystemAccountIdResult().Code },
	xdr.OperationTypeCreateAmlAlert:          				 func(ir xdr.OperationResultTr) shortStr { return ir.MustCreateAmlAlertRequestResult().Code },
	xdr.OperationTypeCreateKycRequest:        				 func(ir xdr.OperationResultTr) shortStr { return ir.MustCreateUpdateKycRequestResult().Code },
	xdr.OperationTypePaymentV2:           				     func(ir xdr.OperationResultTr) shortStr { return ir.MustPaymentV2Result().Code },
	xdr.OperationTypeManageSale:               				 func(ir xdr.OperationResultTr) shortStr { return ir.MustManageSaleResult().Code },
	xdr.OperationTypeManageKeyValue: 						 func(ir xdr.OperationResultTr) shortStr { return ir.MustManageKeyValueResult().Code },
	xdr.OperationTypeCreateManageLimitsRequest:				 func(ir xdr.OperationResultTr) shortStr { return ir.MustCreateManageLimitsRequestResult().Code },
	xdr.OperationTypeBillPay:								 func(ir xdr.OperationResultTr) shortStr { return ir.MustBillPayResult().Code },
}
