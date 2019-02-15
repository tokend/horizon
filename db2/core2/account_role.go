package core2

type AccountRole struct {
	ID      uint64   `db:"id"`
	RuleIDs []uint64 `db:"rule_ids"`
	Details string   `db:"details"`
}
