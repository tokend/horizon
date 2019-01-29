package regources

import "gitlab.com/tokend/go/xdr"

//CreateKYCRequestAttrs - details of corresponding op
type CreateKYCRequest struct {
	Key
	Attributes CreateKYCRequestAttrs `json:"attributes"`
}

//CreateKYCRequestAttrs - details of corresponding op
type CreateKYCRequestAttrs struct {
	AccountAddressToUpdateKYC string          `json:"account_address_to_update_kyc"`
	AccountTypeToSet          xdr.AccountType `json:"account_type_to_set"`
	KYCData                   Details         `json:"kyc_data"`
	AllTasks                  *uint32         `json:"all_tasks"`
	RequestDetails            Request         `json:"request_details"`
}
