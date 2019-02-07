package core2

type ExternalSystemID struct {
	ID                 uint64 `db:"id"`
	AccountID          string `db:"account_id"`
	ExternalSystemType int32  `db:"external_system_type"`
	Data               string `db:"data"`
	IsDeleted          bool   `db:"is_deleted"`
	ExpiresAt          int64  `db:"expires_at"`
	BindedAt           int64  `db:"binded_at"`
}
