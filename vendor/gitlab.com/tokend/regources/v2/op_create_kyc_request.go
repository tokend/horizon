package regources

import "gitlab.com/tokend/go/xdr"

//CreateKYCRequestOpAttrs - details of corresponding op
type CreateKYCRequestOp struct {
	Key
	Attributes    CreateKYCRequestOpAttrs     `json:"attributes"`
	Relationships CreateKYCRequestOpRelations `json:"relationships"`
}

//CreateKYCRequestOpAttrs - details of corresponding op
type CreateKYCRequestOpAttrs struct {
	AccountTypeToSet xdr.AccountType `json:"account_type_to_set"`
	KYCData          Details         `json:"kyc_data"`
	AllTasks         *uint32         `json:"all_tasks"`
}

// CreateKYCRequestOpRelations - relationships of the operation
type CreateKYCRequestOpRelations struct {
	AccountToUpdateKYC *Relation `json:"account_to_update_kyc"`
	Request            *Relation `json:"request"`
}
