package reviewablerequest

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/regources"
)

func PopulateDetails(requestType xdr.ReviewableRequestType, h history.ReviewableRequestDetails) (
	d *regources.ReviewableRequestDetails, err error,
) {
	d = &regources.ReviewableRequestDetails{}
	d.RequestTypeName = requestType.ShortString()
	d.RequestType = int32(requestType)
	switch requestType {
	case xdr.ReviewableRequestTypeAssetCreate:
		d.AssetCreate, err = PopulateAssetCreationRequest(*h.AssetCreation)
		return
	case xdr.ReviewableRequestTypeAssetUpdate:
		d.AssetUpdate, err = PopulateAssetUpdateRequest(*h.AssetUpdate)
		return
	case xdr.ReviewableRequestTypePreIssuanceCreate:
		d.PreIssuanceCreate, err = PopulatePreIssuanceRequest(*h.PreIssuanceCreate)
		return
	case xdr.ReviewableRequestTypeIssuanceCreate:
		d.IssuanceCreate, err = PopulateIssuanceRequest(*h.IssuanceCreate)
		return
	case xdr.ReviewableRequestTypeWithdraw:
		d.Withdraw, err = PopulateWithdrawalRequest(*h.Withdraw)
		return
	case xdr.ReviewableRequestTypeSale:
		d.Sale, err = PopulateSaleCreationRequest(*h.Sale)
		return
	case xdr.ReviewableRequestTypeLimitsUpdate:
		d.LimitsUpdate, err = PopulateLimitsUpdateRequest(*h.LimitsUpdate)
		return
	case xdr.ReviewableRequestTypeTwoStepWithdrawal:
		d.TwoStepWithdraw, err = PopulateWithdrawalRequest(*h.TwoStepWithdraw)
		return
	case xdr.ReviewableRequestTypeAmlAlert:
		d.AMLAlert, err = PopulateAmlAlertRequest(*h.AmlAlert)
		return
	case xdr.ReviewableRequestTypeUpdateKyc:
		d.KYC, err = PopulateUpdateKYCRequest(*h.UpdateKYC)
		return
	case xdr.ReviewableRequestTypeUpdateSaleDetails:
		d.UpdateSaleDetails, err = PopulateUpdateSaleDetailsRequest(*h.UpdateSaleDetails)
		return
	case xdr.ReviewableRequestTypeInvoice:
		d.Invoice, err = PopulateInvoiceRequest(*h.Invoice)
		return
	case xdr.ReviewableRequestTypeUpdateSaleEndTime:
		d.UpdateSaleEndTime, err = PopulateUpdateSaleEndTimeRequest(*h.UpdateSaleEndTimeRequest)
		return
	case xdr.ReviewableRequestTypeUpdatePromotion:
		d.PromotionUpdateRequest, err = PopulatePromotionUpdateRequest(*h.PromotionUpdate)
		return
	case xdr.ReviewableRequestTypeContract:
		d.Contract, err = PopulateContractRequest(*h.Contract)
		return
	default:
		return nil, errors.From(errors.New("unexpected reviewable request type"), map[string]interface{}{
			"request_type": requestType.String(),
		})
	}
}
