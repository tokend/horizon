package core2

import (
	"time"
)

type ExternalSystemID struct {
	ID                 uint64    `db:"id"`
	AccountID          string    `db:"account_id"`
	ExternalSystemType int32     `db:"external_system_type"`
	Data               string    `db:"data"`
	IsDeleted          bool      `db:"is_deleted"`
	ExpiresAt          time.Time `db:"expires_at"`
	BindedAt           time.Time `db:"binded_at"`
}
