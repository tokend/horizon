package operations

type CreateUpdateKYCRequest struct {
	Base
	RequestID          uint64                 `json:"request_id"`
	AccountToUpdateKYC string                 `json:"account_to_update_kyc"`
	AccountTypeToSet   int32                  `json:"account_type_to_set"`
	KYCLevel           uint32                 `json:"kyc_level"`
	KYCData            map[string]interface{} `json:"kyc_data"`
	AllTasks           *uint32                `json:"all_tasks"`
}
