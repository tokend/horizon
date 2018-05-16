package reviewablerequest

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/swarmfund/horizon/db2/history"
)

// Details - provides specific for request type details.
// Note: json key of specific request must be equal to xdr.ReviewableRequestType.ShortString result
type Details struct {
	RequestType
	AssetCreation     *AssetCreationRequest     `json:"asset_create,omitempty"`
	AssetUpdate       *AssetUpdateRequest       `json:"asset_update,omitempty"`
	PreIssuanceCreate *PreIssuanceRequest       `json:"pre_issuance_create,omitempty"`
	IssuanceCreate    *IssuanceRequest          `json:"issuance_create,omitempty"`
	Withdrawal        *WithdrawalRequest        `json:"withdraw,omitempty"`
	TwoStepWithdrawal *WithdrawalRequest        `json:"two_step_withdrawal"`
	Sale              *SaleCreationRequest      `json:"sale,omitempty"`
	LimitsUpdate      *LimitsUpdateRequest      `json:"limits_update"`
	AmlAlert          *AmlAlertRequest          `json:"aml_alert"`
	UpdateKYC         *UpdateKYCRequest         `json:"update_kyc,omitempty"`
	UpdateSaleDetails *UpdateSaleDetailsRequest `json:"update_sale_details"`
}

func (d *Details) Populate(requestType xdr.ReviewableRequestType, h history.ReviewableRequestDetails) error {
	d.RequestType.Populate(requestType)
	switch requestType {
	case xdr.ReviewableRequestTypeAssetCreate:
		d.AssetCreation = new(AssetCreationRequest)
		return d.AssetCreation.Populate(*h.AssetCreation)
	case xdr.ReviewableRequestTypeAssetUpdate:
		d.AssetUpdate = new(AssetUpdateRequest)
		return d.AssetUpdate.Populate(*h.AssetUpdate)
	case xdr.ReviewableRequestTypePreIssuanceCreate:
		d.PreIssuanceCreate = new(PreIssuanceRequest)
		return d.PreIssuanceCreate.Populate(*h.PreIssuanceCreate)
	case xdr.ReviewableRequestTypeIssuanceCreate:
		d.IssuanceCreate = new(IssuanceRequest)
		return d.IssuanceCreate.Populate(*h.IssuanceCreate)
	case xdr.ReviewableRequestTypeWithdraw:
		d.Withdrawal = new(WithdrawalRequest)
		return d.Withdrawal.Populate(*h.Withdrawal)
	case xdr.ReviewableRequestTypeSale:
		d.Sale = new(SaleCreationRequest)
		return d.Sale.Populate(*h.Sale)
	case xdr.ReviewableRequestTypeLimitsUpdate:
		d.LimitsUpdate = new(LimitsUpdateRequest)
		return d.LimitsUpdate.Populate(*h.LimitsUpdate)
	case xdr.ReviewableRequestTypeTwoStepWithdrawal:
		d.TwoStepWithdrawal = new(WithdrawalRequest)
		return d.TwoStepWithdrawal.Populate(*h.TwoStepWithdrawal)
	case xdr.ReviewableRequestTypeAmlAlert:
		d.AmlAlert = new(AmlAlertRequest)
		return d.AmlAlert.Populate(*h.AmlAlert)
	case xdr.ReviewableRequestTypeUpdateKyc:
		d.UpdateKYC = new(UpdateKYCRequest)
		return d.UpdateKYC.Populate(*h.UpdateKYC)
	case xdr.ReviewableRequestTypeUpdateSaleDetails:
		d.UpdateSaleDetails = new(UpdateSaleDetailsRequest)
		return d.UpdateSaleDetails.Populate(*h.UpdateSaleDetails)
	default:
		return errors.From(errors.New("unexpected reviewable request type"), map[string]interface{}{
			"request_type": requestType.String(),
		})
	}
}
