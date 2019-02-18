package core2

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/regources/v2"
)

type SignerRule struct {
	ID        uint64                 `db:"id"`
	Resource  xdr.SignerRuleResource `db:"resource"`
	Action    string                 `db:"action"`
	IsForbid  bool                   `db:"is_forbid"`
	IsDefault bool                   `db:"is_default"`
	OwnerID   string                 `db:"owner_id"`
	Details   regources.Details      `db:"details"`
}
