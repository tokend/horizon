package core2

type AccountRule struct {
	ID       uint64 `db:"id"`
	Resource string `db:"resource"`
	Action   string `db:"action"`
	IsForbid bool   `db:"is_forbid"`
	Details  string `db:"details"`
}
