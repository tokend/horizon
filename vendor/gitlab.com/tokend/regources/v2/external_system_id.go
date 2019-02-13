package regources

import "time"

type ExternalSystemID struct {
	Key
	Attributes ExternalSystemIDAttr `json:"attributes"`
}

type ExternalSystemIDAttr struct {
	AccountID          string    `json:"account_id"`
	ExternalSystemType int32     `json:"external_system_type"`
	Data               Details   `json:"data"`
	IsDeleted          bool      `json:"is_deleted"`
	ExpiresAt          time.Time `json:"expires_at"`
	BindedAt           time.Time `json:"binded_at"`
}
