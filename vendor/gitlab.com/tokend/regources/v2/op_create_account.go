package regources

//CreateAccountOp - stores details of create account operation
type CreateAccountOp struct {
	Key
	Attributes CreateAccountOpAttrs `json:"attributes"`
}

// CreateAccountOpAttrs - stores details of create account operation
type CreateAccountOpAttrs struct {
	AccountAddress string `json:"account_address"`
	AccountRole    uint64 `json:"account_role"`
}
