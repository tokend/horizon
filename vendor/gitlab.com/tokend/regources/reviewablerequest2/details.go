package reviewablerequest2

// Details - provides specific for request type details.
// Note: json key of specific request must be equal to xdr.ReviewableRequestType.ShortString result
type Details struct {
	RequestType
	AssetCreation          *AssetCreationRequest     `json:"asset_create,omitempty"`
	AssetUpdate            *AssetUpdateRequest       `json:"asset_update,omitempty"`
	PreIssuanceCreate      *PreIssuanceRequest       `json:"pre_issuance_create,omitempty"`
	IssuanceCreate         *IssuanceRequest          `json:"issuance_create,omitempty"`
	Withdrawal             *WithdrawalRequest        `json:"withdraw,omitempty"`
	TwoStepWithdrawal      *WithdrawalRequest        `json:"two_step_withdrawal"`
	Sale                   *SaleCreationRequest      `json:"sale,omitempty"`
	LimitsUpdate           *LimitsUpdateRequest      `json:"limits_update,omitempty"`
	AmlAlert               *AmlAlertRequest          `json:"aml_alert"`
	UpdateKYC              *UpdateKYCRequest         `json:"update_kyc,omitempty"`
	UpdateSaleDetails      *UpdateSaleDetailsRequest `json:"update_sale_details"`
	UpdateSaleEndTime      *UpdateSaleEndTimeRequest `json:"update_sale_end_time"`
	PromotionUpdateRequest *PromotionUpdateRequest   `json:"promotion_update_request"`
	Invoice                *InvoiceRequest           `json:"invoice,omitempty"`
	Contract               *ContractRequest          `json:"contract,omitempty"`
}
