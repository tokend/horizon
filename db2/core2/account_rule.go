package core2

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/regources/v2"
)

//AccountRule - defines rule applicable for account roles
type AccountRule struct {
	ID       uint64                  `db:"id"`
	Resource xdr.AccountRuleResource `db:"resource"`
	Action   string                  `db:"action"`
	IsForbid bool                    `db:"is_forbid"`
	Details  regources.Details       `db:"details"`
}
