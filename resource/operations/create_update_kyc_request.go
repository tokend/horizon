package operations

type CreateUpdateKYCRequest struct {
	Base
	RequestID          uint64                 `json:"request_id"`
	AccountToUpdateKYC string                 `json:"account_to_update_kyc"`
	AccountTypeToSet   int32                  `json:"account_type_to_set"`
	KYCLevelToSet      uint32                 `json:"kyc_level_to_set"`
	KYCData            map[string]interface{} `json:"kyc_data"`
	AllTasks           *uint32                `json:"all_tasks"`
}
