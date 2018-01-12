package reviewablerequest

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdr"
)

// Details - provides specific for request type details.
// Note: json key of specific request must be equal to xdr.ReviewableRequestType.ShortString result
type Details struct {
	RequestType
	AssetCreation     *AssetCreationRequest `json:"asset_create,omitempty"`
	AssetUpdate       *AssetUpdateRequest   `json:"asset_update,omitempty"`
	PreIssuanceCreate *PreIssuanceRequest   `json:"pre_issuance_create,omitempty"`
	IssuanceCreate    *IssuanceRequest      `json:"issuance_create,omitempty"`
	Withdrawal        *WithdrawalRequest    `json:"withdraw,omitempty"`
	Sale              *SaleCreationRequest  `json:"sale,omitempty"`
	UpdateKYC		  *UpdateKYCRequest		`json:"updateKYC,omitempty"`
}

func (d *Details) PopulateFromRawJSON(requestType xdr.ReviewableRequestType, rawJSON []byte) error {
	d.RequestType.Populate(requestType)
	err := d.PopulateSpecificRequest(requestType, rawJSON)
	if err != nil {
		return errors.Wrap(err, "failed to populate reviewable request details")
	}

	return nil
}

func (d *Details) PopulateSpecificRequest(requestType xdr.ReviewableRequestType, rawJSON []byte) error {
	switch requestType {
	case xdr.ReviewableRequestTypeAssetCreate:
		d.AssetCreation = new(AssetCreationRequest)
		return d.AssetCreation.PopulateFromRawJsonHistory(rawJSON)
	case xdr.ReviewableRequestTypeAssetUpdate:
		d.AssetUpdate = new(AssetUpdateRequest)
		return d.AssetUpdate.PopulateFromRawJsonHistory(rawJSON)
	case xdr.ReviewableRequestTypePreIssuanceCreate:
		d.PreIssuanceCreate = new(PreIssuanceRequest)
		return d.PreIssuanceCreate.PopulateFromRawJsonHistory(rawJSON)
	case xdr.ReviewableRequestTypeIssuanceCreate:
		d.IssuanceCreate = new(IssuanceRequest)
		return d.IssuanceCreate.PopulateFromRawJsonHistory(rawJSON)
	case xdr.ReviewableRequestTypeWithdraw:
		d.Withdrawal = new(WithdrawalRequest)
		return d.Withdrawal.PopulateFromRawJsonHistory(rawJSON)
	case xdr.ReviewableRequestTypeSale:
		d.Sale = new(SaleCreationRequest)
		return d.Sale.PopulateFromRawJsonHistory(rawJSON)
	case xdr.ReviewableRequestTypeUpdateKyc:
		d.UpdateKYC = new(UpdateKYCRequest)
		return d.UpdateKYC.PopulateFromRawJsonHistory(rawJSON)
	default:
		return errors.From(errors.New("unexpected reviewable request type"), map[string]interface{}{
			"request_type": requestType.String(),
		})
	}
}
