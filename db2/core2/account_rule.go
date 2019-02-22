package core2

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/regources/v2"
)

//AccountRule - defines rule applicable for account roles
type AccountRule struct {
	ID       uint64                  `db:"id"`
	Resource xdr.AccountRuleResource `db:"resource"`
	Action   int32                   `db:"action"`
	Forbids  bool                    `db:"forbids"`
	Details  regources.Details       `db:"details"`
}
