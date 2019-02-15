package core2

type SignerRule struct {
	ID        uint64 `db:"id"`
	Resource  string `db:"resource"`
	Action    string `db:"action"`
	IsForbid  bool   `db:"is_forbid"`
	IsDefault bool   `db:"is_default"`
	OwnerID   string `db:"owner_id"`
	Details   string `db:"details"`
}
