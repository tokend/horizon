package core2

import "gitlab.com/tokend/regources/v2"

//AccountRole - represents role applicable for account
type AccountRole struct {
	ID      uint64            `db:"id"`
	Details regources.Details `db:"details"`
}
