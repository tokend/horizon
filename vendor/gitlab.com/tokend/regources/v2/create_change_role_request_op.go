package regources

//CreateKYCRequestAttrs - details of corresponding op
type CreateKYCRequest struct {
	Key
	Attributes CreateChangeRoleRequestAttrs `json:"attributes"`
}

//CreateChangeRoleRequestDetails - details of corresponding op
type CreateChangeRoleRequestAttrs struct {
	DestinationAccount string  `json:"destination_account"`
	AccountRoleToSet   uint64  `json:"account_role_to_set"`
	KYCData            Details `json:"kyc_data"`
	AllTasks           *uint32 `json:"all_tasks"`
	RequestDetails     Request `json:"request_details"`
}
