package codes

import "gitlab.com/swarmfund/go/xdr"

type codeProvider func(tr xdr.OperationResultTr) shortStr

var codeProviders = map[xdr.OperationType]codeProvider{
	xdr.OperationTypeCreateAccount:            func(ir xdr.OperationResultTr) shortStr { return ir.MustCreateAccountResult().Code },
	xdr.OperationTypePayment:                  func(ir xdr.OperationResultTr) shortStr { return ir.MustPaymentResult().Code },
	xdr.OperationTypeSetOptions:               func(ir xdr.OperationResultTr) shortStr { return ir.MustSetOptionsResult().Code },
	xdr.OperationTypeSetFees:                  func(ir xdr.OperationResultTr) shortStr { return ir.MustSetFeesResult().Code },
	xdr.OperationTypeManageAccount:            func(ir xdr.OperationResultTr) shortStr { return ir.MustManageAccountResult().Code },
	xdr.OperationTypeCreateWithdrawalRequest:  func(ir xdr.OperationResultTr) shortStr { return ir.MustCreateWithdrawalRequestResult().Code },
	xdr.OperationTypeManageBalance:            func(ir xdr.OperationResultTr) shortStr { return ir.MustManageBalanceResult().Code },
	xdr.OperationTypeReviewPaymentRequest:     func(ir xdr.OperationResultTr) shortStr { return ir.MustReviewPaymentRequestResult().Code },
	xdr.OperationTypeManageAsset:              func(ir xdr.OperationResultTr) shortStr { return ir.MustManageAssetResult().Code },
	xdr.OperationTypeSetLimits:                func(ir xdr.OperationResultTr) shortStr { return ir.MustSetLimitsResult().Code },
	xdr.OperationTypeDirectDebit:              func(ir xdr.OperationResultTr) shortStr { return ir.MustDirectDebitResult().Code },
	xdr.OperationTypeManageAssetPair:          func(ir xdr.OperationResultTr) shortStr { return ir.MustManageAssetPairResult().Code },
	xdr.OperationTypeManageOffer:              func(ir xdr.OperationResultTr) shortStr { return ir.MustManageOfferResult().Code },
	xdr.OperationTypeManageInvoice:            func(ir xdr.OperationResultTr) shortStr { return ir.MustManageInvoiceResult().Code },
	xdr.OperationTypeReviewRequest:            func(ir xdr.OperationResultTr) shortStr { return ir.MustReviewRequestResult().Code },
	xdr.OperationTypeCreatePreissuanceRequest: func(ir xdr.OperationResultTr) shortStr { return ir.MustCreatePreIssuanceRequestResult().Code },
	xdr.OperationTypeCreateIssuanceRequest:    func(ir xdr.OperationResultTr) shortStr { return ir.MustCreateIssuanceRequestResult().Code },
	xdr.OperationTypeCreateSaleRequest:        func(ir xdr.OperationResultTr) shortStr { return ir.MustCreateSaleCreationRequestResult().Code },
	xdr.OperationTypeCheckSaleState:           func(ir xdr.OperationResultTr) shortStr { return ir.MustCheckSaleStateResult().Code },
	xdr.OperationTypeCreateAmlAlert:           func(ir xdr.OperationResultTr) shortStr { return ir.MustCreateAmlAlertRequestResult().Code },
	xdr.OperationTypeCreateKycRequest:         func(ir xdr.OperationResultTr) shortStr { return ir.MustCreateUpdateKycRequestResult().Code },
}
