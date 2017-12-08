package core

import "gitlab.com/swarmfund/go/xdr"

// ExternalSystemAccountID - represents id generated for user in external system like Bitcoin/Ethereum
type ExternalSystemAccountID struct {
	// AccountID - account id of the owner of ID
	AccountID          string                 `db:"account_id"`
	// ExternalSystemType - type of the external system Bitcoin/Ethereum
	ExternalSystemType xdr.ExternalSystemType `db:"external_system_type"`
	// Data - ID from external system
	Data               string                 `db:"data"`
}
