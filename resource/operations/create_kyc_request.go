package operations

type CreateKYCRequest struct {
	Base
	RequestID        uint64                 `json:"request_id"`
	UpdatedAccount   string                 `json:"updated_account"`
	AccountTypeToSet int32                  `json:"account_type_to_set"`
	KYCLevel         uint32                 `json:"kyc_level"`
	KYCData          map[string]interface{} `json:"kyc_data"`
}
