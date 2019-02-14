package core2

//Signer - represents Signer Entry
type Signer struct {
	AccountID string `db:"account_id"`
	ID        string `db:"public_key"`
	Weight    int    `db:"weight"`
	Type      int    `db:"role_id"`
	Identity  int    `db:"identity"`
	Details   string `db:"details"`
}
