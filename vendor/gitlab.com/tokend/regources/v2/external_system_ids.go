package regources

type ExternalSystemID struct {
	Key
	Attributes ExternalSystemIDAttr `json:"attributes"`
}

type ExternalSystemIDAttr struct {
	AccountID          string `json:"account_id"`
	ExternalSystemType int32  `json:"external_system_type"`
	Data               string `json:"data"`
	IsDeleted          bool   `json:"is_deleted"`
	ExpiresAt          int64  `json:"expires_at"`
	BindedAt           int64  `json:"binded_at"`
}
