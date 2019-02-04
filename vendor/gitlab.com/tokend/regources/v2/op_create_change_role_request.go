package regources

//CreateChangeRoleRequest - details of corresponding op
type CreateChangeRoleRequest struct {
	Key
	Attributes    CreateChangeRoleRequestAttrs       `json:"attributes"`
	Relationships CreateChangeRoleRequestOpRelations `json:"relationships"`
}

//CreateChangeRoleRequestAttrs - details of corresponding op
type CreateChangeRoleRequestAttrs struct {
	AccountRoleToSet uint64  `json:"account_role_to_set"`
	KYCData          Details `json:"kyc_data"`
	AllTasks         *uint32 `json:"all_tasks"`
}

// CreateKYCRequestOpRelations - relationships of the operation
type CreateChangeRoleRequestOpRelations struct {
	AccountToUpdateRole *Relation `json:"account_to_update_role"`
	Request             *Relation `json:"request"`
}
