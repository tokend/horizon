package operations

type CreateKYCRequest struct {
	Base
	RequestID        uint64                 `json:"request_id"`
	AccountID        string                 `json:"account_id"`
	AccountTypeToSet int32                  `json:"account_type"`
	KYCLevel         uint32                 `json:"kyc_level"`
	KYCData          map[string]interface{} `json:"kyc_data"`
}
