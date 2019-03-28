package core2

import "gitlab.com/tokend/regources/rgenerated"

//Signer - represents Signer Entry
type Signer struct {
	AccountID string             `db:"account_id"`
	PublicKey string             `db:"public_key"`
	Weight    uint32             `db:"weight"`
	RoleID    uint64             `db:"role_id"`
	Identity  uint32             `db:"identity"`
	Details   rgenerated.Details `db:"details"`
}
