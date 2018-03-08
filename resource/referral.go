package resource

import (
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/swarmfund/go/xdr"
)

type Referral struct {
	AccountID    string `json:"account_id"`
	AccountTypeI int32  `json:"account_type_i"`
	AccountType  string `json:"account_type"`
	Referrer     string `json:"referrer"`
}

// Populate fills out the resource's fields
func (r *Referral) Populate(ca core.Account) {
	r.AccountID = ca.AccountID
	r.AccountTypeI = ca.AccountType
	r.AccountType = xdr.AccountType(ca.AccountType).String()
	r.Referrer = ca.Referrer
}
