package core2

import "gitlab.com/tokend/regources/rgenerated"

//AccountRole - represents role applicable for account
type AccountRole struct {
	ID      uint64             `db:"id"`
	Details rgenerated.Details `db:"details"`
}
