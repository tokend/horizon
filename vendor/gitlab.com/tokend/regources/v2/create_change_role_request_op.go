package regources

//CreateChangeRoleRequest - details of corresponding op
type CreateChangeRoleRequest struct {
	Key
	Attributes CreateChangeRoleRequestAttrs `json:"attributes"`
}

//CreateChangeRoleRequestAttrs - details of corresponding op
type CreateChangeRoleRequestAttrs struct {
	DestinationAccount string  `json:"destination_account"`
	AccountRoleToSet   uint64  `json:"account_role_to_set"`
	KYCData            Details `json:"kyc_data"`
	AllTasks           *uint32 `json:"all_tasks"`
	RequestDetails     Request `json:"request_details"`
}
