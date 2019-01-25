package regources

import "gitlab.com/tokend/go/xdr"

//CreateAccount - stores details of create account operation
type CreateAccount struct {
	Key
	Attributes CreateAccountAttrs `json:"attributes"`
}

// CreateAccountAttrs - stores details of create account operation
type CreateAccountAttrs struct {
	AccountAddress string          `json:"account_address"`
	AccountType    xdr.AccountType `json:"account_type"`
}
