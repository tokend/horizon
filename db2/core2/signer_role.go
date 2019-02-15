package core2

type SignerRole struct {
	ID      uint64   `db:"id"`
	RuleIDs []uint64 `db:"rule_ids"`
	OwnerID string   `db:"owner_id"`
	Details string   `db:"details"`
}
