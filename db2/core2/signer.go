package core2

//Signer - represents Signer Entry
type Signer struct {
	AccountID string `db:"accountid"`
	ID        string `db:"publickey"`
	Weight    int    `db:"weight"`
	Type      int    `db:"signer_type"`
	Identity  int    `db:"identity_id"`
	Name      string `db:"signer_name"`
}
