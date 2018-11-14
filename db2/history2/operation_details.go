package history2

import "gitlab.com/tokend/go/xdr"

// OperationDetails - stores details of the operation performed in union switch form. Only one value must be selected at
// a type
type OperationDetails struct {
	Type          xdr.OperationType     `json:"type"`
	CreateAccount *CreateAccountDetails `json:"create_account"`
	Payment       *PaymentDetails       `json:"payment"`
}

// CreateAccountDetails - stores details of create account operation
type CreateAccountDetails struct {
}

// PaymentDetails - stores details of payment operation
type PaymentDetails struct {
}
