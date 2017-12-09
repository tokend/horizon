package codes

import "gitlab.com/swarmfund/go/xdr"

type codeProvider func(tr xdr.OperationResultTr) interface{}

var codeProviders = map[xdr.OperationType]codeProvider{
	xdr.OperationTypeCreateAccount:            func(ir xdr.OperationResultTr) interface{} { return ir.MustCreateAccountResult().Code },
	xdr.OperationTypePayment:                  func(ir xdr.OperationResultTr) interface{} { return ir.MustPaymentResult().Code },
	xdr.OperationTypeSetOptions:               func(ir xdr.OperationResultTr) interface{} { return ir.MustSetOptionsResult().Code },
	xdr.OperationTypeSetFees:                  func(ir xdr.OperationResultTr) interface{} { return ir.MustSetFeesResult().Code },
	xdr.OperationTypeManageAccount:            func(ir xdr.OperationResultTr) interface{} { return ir.MustManageAccountResult().Code },
	xdr.OperationTypeManageForfeitRequest:     func(ir xdr.OperationResultTr) interface{} { return ir.MustManageForfeitRequestResult().Code },
	xdr.OperationTypeRecover:                  func(ir xdr.OperationResultTr) interface{} { return ir.MustRecoverResult().Code },
	xdr.OperationTypeManageBalance:            func(ir xdr.OperationResultTr) interface{} { return ir.MustManageBalanceResult().Code },
	xdr.OperationTypeReviewPaymentRequest:     func(ir xdr.OperationResultTr) interface{} { return ir.MustReviewPaymentRequestResult().Code },
	xdr.OperationTypeManageAsset:              func(ir xdr.OperationResultTr) interface{} { return ir.MustManageAssetResult().Code },
	xdr.OperationTypeSetLimits:                func(ir xdr.OperationResultTr) interface{} { return ir.MustSetLimitsResult().Code },
	xdr.OperationTypeDirectDebit:              func(ir xdr.OperationResultTr) interface{} { return ir.MustDirectDebitResult().Code },
	xdr.OperationTypeManageAssetPair:          func(ir xdr.OperationResultTr) interface{} { return ir.MustManageAssetPairResult().Code },
	xdr.OperationTypeManageOffer:              func(ir xdr.OperationResultTr) interface{} { return ir.MustManageOfferResult().Code },
	xdr.OperationTypeManageInvoice:            func(ir xdr.OperationResultTr) interface{} { return ir.MustManageInvoiceResult().Code },
	xdr.OperationTypeReviewRequest:            func(ir xdr.OperationResultTr) interface{} { return ir.MustReviewRequestResult().Code },
	xdr.OperationTypeCreatePreissuanceRequest: func(ir xdr.OperationResultTr) interface{} { return ir.MustCreatePreIssuanceRequestResult().Code },
	xdr.OperationTypeCreateIssuanceRequest:    func(ir xdr.OperationResultTr) interface{} { return ir.MustCreateIssuanceRequestResult().Code },
}
