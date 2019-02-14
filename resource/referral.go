package resource

import (
	"gitlab.com/tokend/horizon/db2/core"
)

type Referral struct {
	AccountID    string `json:"account_id"`
	AccountTypeI int32  `json:"account_type_i"`
	AccountType  string `json:"account_type"`
}

// Populate fills out the resource's fields
func (r *Referral) Populate(ca core.Account) {
	r.AccountID = ca.AccountID
	r.AccountTypeI = ca.AccountType
}
