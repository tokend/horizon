package regources

import "gitlab.com/tokend/go/xdr"

//CreateAccountOp - stores details of create account operation
type CreateAccountOp struct {
	Key
	Attributes CreateAccountOpAttrs `json:"attributes"`
}

// CreateAccountOpAttrs - stores details of create account operation
type CreateAccountOpAttrs struct {
	AccountAddress string          `json:"account_address"`
	AccountType    xdr.AccountType `json:"account_type"`
}
