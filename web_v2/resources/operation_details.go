package resources

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/v2"
)

//NewOperationDetails - populates operation details into appropriate resource
func NewOperationDetails(op history2.Operation) regources.Resource {
	switch op.Type {
	case xdr.OperationTypeCreateAccount:
		return &regources.CreateAccount{
			Key:        regources.NewKeyInt64(op.ID, regources.TypeCreateAccount),
			Attributes: regources.CreateAccountAttrs(*op.Details.CreateAccount),
		}
	case xdr.OperationTypeSetOptions:
		return regources.NewKeyInt64(op.ID, regources.TypeSetOptions).GetKeyP()
	case xdr.OperationTypeCreateIssuanceRequest:
		return newCreateIssuanceOpDetails(op.ID, *op.Details.CreateIssuanceRequest)
	case xdr.OperationTypeSetFees:
		return &regources.SetFee{
			Key:        regources.NewKeyInt64(op.ID, regources.TypeSetFees),
			Attributes: regources.SetFeeAttrs(*op.Details.SetFee),
		}
	case xdr.OperationTypeManageAccount:
		return &regources.ManageAccount{
			Key:        regources.NewKeyInt64(op.ID, regources.TypeManageAccount),
			Attributes: regources.ManageAccountAttrs(*op.Details.ManageAccount),
		}
	case xdr.OperationTypeCreateWithdrawalRequest:
		return newCreateWithdrawalRequestOp(op.ID, *op.Details.CreateWithdrawRequest)
	case xdr.OperationTypeManageBalance:
		return &regources.ManageBalance{
			Key:        regources.NewKeyInt64(op.ID, regources.TypeManageBalance),
			Attributes: regources.ManageBalanceAttrs(*op.Details.ManageBalance),
		}
	case xdr.OperationTypeManageAsset:
		return &regources.ManageAsset{
			Key:        regources.NewKeyInt64(op.ID, regources.TypeManageAsset),
			Attributes: regources.ManageAssetAttrs(*op.Details.ManageAsset),
		}
	case xdr.OperationTypeCreatePreissuanceRequest:
		return &regources.CreatePreIssuanceRequest{
			Key:        regources.NewKeyInt64(op.ID, regources.TypeCreatePreissuanceRequest),
			Attributes: regources.CreatePreIssuanceRequestAttrs(*op.Details.CreatePreIssuanceRequest),
		}
	case xdr.OperationTypeManageLimits:
		return newManageLimits(op.ID, *op.Details.ManageLimits)
	case xdr.OperationTypeManageAssetPair:
		return &regources.ManageAssetPair{
			Key:        regources.NewKeyInt64(op.ID, regources.TypeManageAssetPair),
			Attributes: regources.ManageAssetPairAttrs(*op.Details.ManageAssetPair),
		}
	case xdr.OperationTypeManageOffer:
		return &regources.ManageOffer{
			Key:        regources.NewKeyInt64(op.ID, regources.TypeManageOffer),
			Attributes: regources.ManageOfferAttrs(*op.Details.ManageOffer),
		}
	case xdr.OperationTypeManageInvoiceRequest:
		return regources.NewKeyInt64(op.ID, regources.TypeManageInvoiceRequest).GetKeyP()
	case xdr.OperationTypeReviewRequest:
		return newReviewRequestOpDetails(op.ID, *op.Details.ReviewRequest)
	case xdr.OperationTypeCreateSaleRequest:
		return &regources.CreateSaleRequest{
			Key:        regources.NewKeyInt64(op.ID, regources.TypeCreateSaleRequest),
			Attributes: regources.CreateSaleRequestAttrs(*op.Details.CreateSaleRequest),
		}
	case xdr.OperationTypeCheckSaleState:
		return &regources.CheckSaleState{
			Key:        regources.NewKeyInt64(op.ID, regources.TypeCheckSaleState),
			Attributes: regources.CheckSaleStateAttrs(*op.Details.CheckSaleState),
		}
	case xdr.OperationTypeCreateAmlAlert:
		return &regources.CreateAMLAlertRequest{
			Key:        regources.NewKeyInt64(op.ID, regources.TypeCreateAmlAlert),
			Attributes: regources.CreateAMLAlertRequestAttrs(*op.Details.CreateAMLAlertRequest),
		}
	case xdr.OperationTypeCreateKycRequest:
		return newKeyRequest(op.ID, *op.Details.CreateKYCRequest)
	case xdr.OperationTypePaymentV2:
		return &regources.Payment{
			Key:        regources.NewKeyInt64(op.ID, regources.TypePaymentV2),
			Attributes: regources.PaymentAttrs(*op.Details.Payment),
		}
	case xdr.OperationTypeManageExternalSystemAccountIdPoolEntry:
		return newManageExternalSystemPool(op.ID, *op.Details.ManageExternalSystemPool)
	case xdr.OperationTypeBindExternalSystemAccountId:
		return &regources.BindExternalSystemAccount{
			Key:        regources.NewKeyInt64(op.ID, regources.TypeBindExternalSystemAccountID),
			Attributes: regources.BindExternalSystemAccountAttrs(*op.Details.BindExternalSystemAccount),
		}
	case xdr.OperationTypeManageSale:
		return &regources.ManageSale{
			Key:        regources.NewKeyInt64(op.ID, regources.TypeManageSale),
			Attributes: regources.ManageSaleAttrs(*op.Details.ManageSale),
		}
	case xdr.OperationTypeManageKeyValue:
		return &regources.ManageKeyValue{
			Key:        regources.NewKeyInt64(op.ID, regources.TypeManageKeyValue),
			Attributes: regources.ManageKeyValueAttrs(*op.Details.ManageKeyValue),
		}
	case xdr.OperationTypeCreateManageLimitsRequest:
		return &regources.CreateManageLimitsRequest{
			Key:        regources.NewKeyInt64(op.ID, regources.TypeCreateManageLimitsRequest),
			Attributes: regources.CreateManageLimitsRequestAttrs(*op.Details.CreateManageLimitsRequest),
		}
	case xdr.OperationTypeManageContractRequest:
		return regources.NewKeyInt64(op.ID, regources.TypeManageContractRequest).GetKeyP()
	case xdr.OperationTypeManageContract:
		return regources.NewKeyInt64(op.ID, regources.TypeManageContract).GetKeyP()
	case xdr.OperationTypeCancelSaleRequest:
		return regources.NewKeyInt64(op.ID, regources.TypeCancelSaleRequest).GetKeyP()
	case xdr.OperationTypePayout:
		return &regources.Payout{
			Key:        regources.NewKeyInt64(op.ID, regources.TypePayout),
			Attributes: regources.PayoutAttrs(*op.Details.Payout),
		}
	case xdr.OperationTypeManageAccountRole:
		return regources.NewKeyInt64(op.ID, regources.TypeManageAccountRole).GetKeyP()
	case xdr.OperationTypeManageAccountRolePermission:
		return regources.NewKeyInt64(op.ID, regources.TypeManageAccountRolePermission).GetKeyP()
	case xdr.OperationTypeCreateAswapBidRequest:
		return regources.NewKeyInt64(op.ID, regources.TypeCreateAswapBidRequest).GetKeyP()
	case xdr.OperationTypeCancelAswapBid:
		return regources.NewKeyInt64(op.ID, regources.TypeCancelAswapBid).GetKeyP()
	case xdr.OperationTypeCreateAswapRequest:
		return regources.NewKeyInt64(op.ID, regources.TypeCreateAswapBidRequest).GetKeyP()
	default:
		panic(errors.From(errors.New("unexpected operation type"), logan.F{
			"type": op.Type,
		}))
	}
}
